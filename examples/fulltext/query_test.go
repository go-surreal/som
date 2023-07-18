package fulltext

import (
	"context"
	"github.com/marcbinz/som/examples/basic/gen/som"
	"github.com/marcbinz/som/examples/basic/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
	"math/rand"
	"testing"
)

const (
	randMin = 5
	randMax = 20
)

func TestQueryCount(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        "surrealdb/surrealdb:" + surrealDBContainerVersion,
		Cmd:          []string{"start", "--log", "debug", "--user", "root", "--pass", "root", "memory"},
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForLog(containerStartedMsg),
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

	client, err := som.NewClient(conf(endpoint))
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	if err := client.ApplySchema(); err != nil {
		t.Fatal(err)
	}

	count := rand.Intn(randMax-randMin) + randMin

	for i := 0; i < count; i++ {
		err = client.UserRepo().Create(ctx, &model.User{})
		if err != nil {
			t.Fatal(err)
		}
	}

	dbCount, err := client.UserRepo().Query().Count(ctx)

	if err != nil {
		t.Fatal(err)
	}

	// TODO: add database cleanup?

	assert.Equal(t, count, dbCount)
}
