package testing

import (
	"context"
	"github.com/marcbinz/som/examples/testing/gen/som"
	"github.com/marcbinz/som/examples/testing/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

const (
	randMin = 5
	randMax = 20
)

const (
	surrealDBContainerVersion = "1.0.0-beta.9"
	containerName             = "som_test_surrealdb"
	containerStartedMsg       = "Started web server on 0.0.0.0:8000"
)

func conf(endpoint string) som.Config {
	return som.Config{
		Address:   "ws://" + endpoint,
		Username:  "root",
		Password:  "root",
		Namespace: "som_test",
		Database:  "example_basic",
	}
}

func TestCreateWithFieldsLikeDBResponse(t *testing.T) {
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

	newModel := &model.FieldsLikeDBResponse{
		Status: "some value",
	}

	err = client.FieldsLikeDBResponseRepo().Create(ctx, newModel)
	if err != nil {
		t.Fatal(err)
	}

	readModel, exists, err := client.FieldsLikeDBResponseRepo().Read(ctx, newModel.ID())
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, exists)
	assert.Equal(t, "some value", readModel.Status)

	readModel.Status = "some other value"

	err = client.FieldsLikeDBResponseRepo().Update(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "some other value", readModel.Status)

	err = client.FieldsLikeDBResponseRepo().Delete(ctx, readModel)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: add database cleanup?
}