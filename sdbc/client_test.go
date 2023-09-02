package sdbc

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
	"log/slog"
	"os"
	"testing"
)

const (
	surrealDBContainerVersion = "1.0.0-beta.10"
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
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Name:  containerName,
		Image: "surrealdb/surrealdb:" + surrealDBContainerVersion,
		Cmd: []string{
			"start", "--auth", "--strict", "--allow-funcs",
			"--user", surrealUser,
			"--pass", surrealPass,
			"--log", "trace", "memory",
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

	client, err := NewClient(ctx, conf(endpoint))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	_, err = client.Query(ctx, 0, "define table test schemaless", nil)
	if err != nil {
		t.Fatal(err)
	}

	create, err := client.Create(ctx, 0, "test", nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(create), string(create))
}
