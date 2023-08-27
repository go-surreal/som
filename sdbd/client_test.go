package sdbd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
	"log/slog"
	"os"
	"testing"
)

const (
	surrealDBContainerVersion = "nightly"
	containerName             = "sdbd_test_surrealdb"
	containerStartedMsg       = "Started web server on 0.0.0.0:8000"
	surrealUser               = "root"
	surrealPass               = "root"
)

func conf(endpoint string) Config {
	return Config{
		Address:   "ws://" + endpoint + "/rpc",
		Username:  surrealUser,
		Password:  surrealPass,
		Namespace: "test",
		Database:  "test",
	}
}

func TestClient(t *testing.T) {
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Name:  containerName,
		Image: "surrealdb/surrealdb:" + surrealDBContainerVersion,
		Cmd: []string{
			"start", "--strict",
			"--user", surrealUser,
			"--pass", surrealPass,
			"--log", "debug", "memory",
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
		t.Fatal(err)
	}

	defer func() {
		if err := surreal.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	endpoint, err := surreal.Endpoint(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	))

	slog.Info("pn3fo3enf")

	client, err := NewClient(ctx, conf(endpoint))
	if err != nil {
		t.Fatal(err)
	}

	slog.Info("304ifn349i")

	defer func() {
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	slog.Info("pn3fo3enf")

	_, err = client.Query(ctx, 0, "define table test schemaless", nil)
	if err != nil {
		t.Fatal(err)
	}

	create, err := client.Create(ctx, 0, "test", nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("create:", string(create))
	slog.Info("epo3fm3ßfon34")

	assert.Equal(t, string(create), "")
}
