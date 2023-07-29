package sdbd

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
	"testing"
)

const (
	surrealDBContainerVersion = "nightly"
	containerName             = "sdbd_test_surrealdb"
	containerStartedMsg       = "Started web server on 0.0.0.0:8000"
)

func conf(endpoint string) Config {
	return Config{
		Address:   "ws://" + endpoint + "/rpc",
		Username:  "root",
		Password:  "root",
		Namespace: "test",
		Database:  "test",
	}
}

func TestClient(t *testing.T) {
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        "surrealdb/surrealdb:" + surrealDBContainerVersion,
		Cmd:          []string{"start", "--strict", "--user", "root", "--pass", "root", "--log", "debug", "memory"},
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

	client, err := NewClient(ctx, conf(endpoint))
	if err != nil {
		t.Fatal(err)
	}

	create, err := client.Create(ctx, "test", nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(create), "")
}
