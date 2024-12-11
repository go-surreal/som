package basic

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/docker/docker/api/types/container"
	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"testing"
)

const (
	surrealDBVersion    = "2.1.3"
	containerStartedMsg = "Started web server on "
)

// errAlreadyInProgress is a regular expression that matches the
// error for a container removal that is already in progress.
var errAlreadyInProgress = regexp.MustCompile(`removal of container .* is already in progress`)

// errNoSuchContainer is a regular expression that matches the
// error for a container that does not exist.
var errNoSuchContainer = regexp.MustCompile(`No such container`)

func prepareDatabase(ctx context.Context, tb testing.TB) (som.Client, func()) {
	tb.Helper()

	username := gofakeit.Username()
	password := gofakeit.Password(true, true, true, true, true, 32)
	namespace := gofakeit.FirstName()
	database := gofakeit.LastName()

	req := testcontainers.ContainerRequest{
		Name:  "sdbc_" + toSlug(tb.Name()),
		Image: "surrealdb/surrealdb:v" + surrealDBVersion,
		Env: map[string]string{
			"SURREAL_PATH":   "memory",
			"SURREAL_STRICT": "true",
			"SURREAL_USER":   username,
			"SURREAL_PASS":   password,
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
			Reuse:            true,
		},
	)
	if err != nil {
		tb.Fatal(err)
	}

	endpoint, err := surreal.Endpoint(ctx, "")
	if err != nil {
		tb.Fatal(err)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	config := som.Config{
		Host:      endpoint,
		Username:  username,
		Password:  password,
		Namespace: namespace,
		Database:  database,
	}

	opts := []som.Option{
		som.WithLogger(slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)),
	}

	client, err := som.NewClient(ctx, config, opts...)
	if err != nil {
		tb.Fatal(err)
	}

	if err := client.ApplySchema(ctx); err != nil {
		tb.Fatal(err)
	}

	cleanup := func() {
		client.Close()

		if err := surreal.Terminate(ctx); err != nil {
			if errAlreadyInProgress.MatchString(err.Error()) || errNoSuchContainer.MatchString(err.Error()) {
				return // this "error" is not caught by the Terminate method even though it is safe to ignore
			}

			tb.Fatalf("failed to terminate container: %s", err.Error())
		}
	}

	return client, cleanup
}

func toSlug(input string) string {
	// Remove special characters
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	processedString := reg.ReplaceAllString(input, " ")

	// Remove leading and trailing spaces
	processedString = strings.TrimSpace(processedString)

	// Replace spaces with dashes
	slug := strings.ReplaceAll(processedString, " ", "-")

	// Convert to lowercase
	slug = strings.ToLower(slug)

	return slug
}
