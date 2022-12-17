package example

import (
	"context"
	"github.com/marcbinz/som/example/gen/som"
	"github.com/marcbinz/som/example/gen/som/where"
	"github.com/marcbinz/som/example/model"
	"gotest.tools/assert"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	containerName = "som_test_surrealdb"
)

func TestWithSurreal(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        "surrealdb/surrealdb:1.0.0-beta.8",
		Cmd:          []string{"start", "--log", "debug", "--user", "root", "--pass", "root", "memory"},
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForLog("Started web server on 0.0.0.0:8000"),
	}

	surreal, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Reuse:            true,
		},
	)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := surreal.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	endpoint, err := surreal.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
	}

	client, err := som.NewClient("ws://"+endpoint, "root", "root", "som_test", "som_test")
	if err != nil {
		t.Error(err)
	}

	defer client.Close()

	err = client.User().Create(ctx, &model.User{
		String: "some_user",
	})
	if err != nil {
		t.Error(err)
	}

	user, err := client.User().Query().
		Filter(
			where.User.String.Equal("some_user"),
		).
		First()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "some_user", user.String)
}
