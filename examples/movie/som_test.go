package basic

import (
	"context"
	"github.com/google/go-cmp/cmp/cmpopts"
	sombase "github.com/marcbinz/som"
	"github.com/marcbinz/som/examples/movie/gen/som"
	"github.com/marcbinz/som/examples/movie/gen/som/where"
	"github.com/marcbinz/som/examples/movie/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/v3/assert"
	"testing"
)

const (
	surrealDBContainerVersion = "1.0.0-beta.8" // "nightly"
	containerName             = "som_test_surrealdb"
	containerStartedMsg       = "Started web server on 0.0.0.0:8000"
)

func conf(endpoint string) som.Config {
	return som.Config{
		Address:   "ws://" + endpoint,
		Username:  "root",
		Password:  "root",
		Namespace: "som_test",
		Database:  "example_movie",
	}
}

func TestWithDatabase(t *testing.T) {
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

	title := "Some Movie"

	movieNew := model.Movie{
		Title: title,
	}

	movieIn := movieNew

	err = client.MovieRepo().Create(ctx, &movieIn)
	if err != nil {
		t.Fatal(err)
	}

	movieOut, err := client.MovieRepo().Query().
		Filter(
			where.Movie.ID.Equal(movieIn.ID()),
			where.Movie.Title.Equal(title),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, title, movieOut.Title)

	assert.DeepEqual(t,
		movieNew, *movieOut,
		cmpopts.IgnoreUnexported(sombase.Node{}, sombase.Timestamps{}),
	)
}
