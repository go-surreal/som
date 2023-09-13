package basic

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func prepareDatabase(ctx context.Context, tb testing.TB) (som.Client, func()) {
	tb.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        "surrealdb/surrealdb:" + surrealDBContainerVersion,
		Cmd:          []string{"start", "--strict", "--allow-funcs", "--user", "root", "--pass", "root", "--log", "debug", "memory"},
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

	client, err := som.NewClient(ctx, conf(endpoint))
	if err != nil {
		tb.Fatal(err)
	}

	if err := client.ApplySchema(ctx); err != nil {
		tb.Fatal(err)
	}

	cleanup := func() {
		client.Close()

		if err := surreal.Terminate(ctx); err != nil {
			tb.Fatalf("failed to terminate container: %s", err.Error())
		}
	}

	return client, cleanup
}
