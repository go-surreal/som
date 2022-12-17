package example

import (
	"context"
	"github.com/marcbinz/som/example/gen/som"
	"github.com/marcbinz/som/example/gen/som/where"
	"github.com/marcbinz/som/example/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/assert"
	"testing"
)

const (
	containerName = "som_test_surrealdb"
)

func TestWithDatabase(t *testing.T) {
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

	str := "Some User"

	err = client.User().Create(ctx, &model.User{
		String: str,
	})
	if err != nil {
		t.Error(err)
	}

	user, err := client.User().Query().
		Filter(
			where.User.String.Equal(str),
		).
		First()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, str, user.String)
}

func FuzzWithDatabase(f *testing.F) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "surrealdb/surrealdb:1.0.0-beta.8",
		Cmd:          []string{"start", "--log", "debug", "--user", "root", "--pass", "root", "memory"},
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForLog("Started web server on 0.0.0.0:8000"),
	}

	surreal, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		f.Error(err)
	}

	defer func() {
		if err := surreal.Terminate(ctx); err != nil {
			f.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	endpoint, err := surreal.Endpoint(ctx, "")
	if err != nil {
		f.Error(err)
	}

	client, err := som.NewClient("ws://"+endpoint, "root", "root", "som_test", "som_test")
	if err != nil {
		f.Error(err)
	}

	defer client.Close()

	f.Add("Some User")

	f.Fuzz(func(t *testing.T, str string) {
		userIn := &model.User{
			String: str,
		}

		err = client.User().Create(ctx, userIn)
		if err != nil {
			t.Error(err)
		}

		if userIn.ID == "" {
			t.Error("user ID must not be empty after create call")
		}

		userOut, err := client.User().Query().
			Filter(
				where.User.ID.Equal(userIn.ID),
			).
			First()

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, userIn.String, userOut.String)
	})
}

func BenchmarkWithDatabase(b *testing.B) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "surrealdb/surrealdb:1.0.0-beta.8",
		Cmd:          []string{"start", "--log", "debug", "--user", "root", "--pass", "root", "memory"},
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForLog("Started web server on 0.0.0.0:8000"),
	}

	surreal, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		b.Error(err)
	}

	defer func() {
		if err := surreal.Terminate(ctx); err != nil {
			b.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	endpoint, err := surreal.Endpoint(ctx, "")
	if err != nil {
		b.Error(err)
	}

	client, err := som.NewClient("ws://"+endpoint, "root", "root", "som_test", "som_test")
	if err != nil {
		b.Error(err)
	}

	defer client.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		userIn := &model.User{
			String: "Some User",
		}

		err = client.User().Create(ctx, userIn)
		if err != nil {
			b.Error(err)
		}

		if userIn.ID == "" {
			b.Error("user ID must not be empty after create call")
		}

		userOut, err := client.User().Query().
			Filter(
				where.User.ID.Equal(userIn.ID),
			).
			First()

		if err != nil {
			b.Error(err)
		}

		assert.Equal(b, userIn.String, userOut.String)
	}
}
