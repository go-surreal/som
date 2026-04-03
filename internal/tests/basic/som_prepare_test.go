package basic

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"som.test/gen/som/repo"
)

const (
	surrealDBVersion    = "3.0.5"
	containerStartedMsg = "Started web server on "
)

// errAlreadyInProgress is a regular expression that matches the
// error for a container removal that is already in progress.
var errAlreadyInProgress = regexp.MustCompile(`removal of container .* is already in progress`)

// errNoSuchContainer is a regular expression that matches the
// error for a container that does not exist.
var errNoSuchContainer = regexp.MustCompile(`No such container`)

var sharedEndpoint string
var sharedUsername string
var sharedPassword string

func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	ctx := context.Background()

	sharedUsername = gofakeit.Username()
	sharedPassword = gofakeit.Password(true, true, true, true, true, 32)

	req := testcontainers.ContainerRequest{
		Name:  "som_test_shared",
		Image: "surrealdb/surrealdb:v" + surrealDBVersion,
		Env: map[string]string{
			"SURREAL_PATH":   "memory",
			"SURREAL_STRICT": "true",
			"SURREAL_USER":   sharedUsername,
			"SURREAL_PASS":   sharedPassword,
		},
		Cmd: []string{
			"start", "--allow-funcs", "--log", "trace",
		},
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForLog(containerStartedMsg),
		HostConfigModifier: func(conf *container.HostConfig) {
			conf.AutoRemove = true
		},
	}

	surreal, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Reuse:            false,
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start shared container: %v\n", err)
		os.Exit(1)
	}

	endpoint, err := surreal.Endpoint(ctx, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get container endpoint: %v\n", err)
		os.Exit(1)
	}

	sharedEndpoint = endpoint

	code := m.Run()

	if err := surreal.Terminate(ctx); err != nil {
		if !errAlreadyInProgress.MatchString(err.Error()) && !errNoSuchContainer.MatchString(err.Error()) {
			fmt.Fprintf(os.Stderr, "failed to terminate container: %v\n", err)
		}
	}

	os.Exit(code)
}

func prepareDatabase(ctx context.Context, tb testing.TB) (repo.Client, func()) {
	tb.Helper()

	namespace := gofakeit.FirstName()
	database := gofakeit.LastName()

	config := repo.Config{
		Address:   "ws://" + sharedEndpoint,
		Username:  sharedUsername,
		Password:  sharedPassword,
		Namespace: namespace,
		Database:  database,
	}

	var client *repo.ClientImpl
	var err error
	for range 5 {
		client, err = repo.NewClient(ctx, config)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		tb.Fatal(err)
	}

	if err := client.ApplySchema(ctx); err != nil {
		tb.Fatal(err)
	}

	cleanup := func() {
		client.Close()
	}

	return client, cleanup
}

