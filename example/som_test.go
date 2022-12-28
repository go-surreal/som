package example

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/marcbinz/som/example/gen/som"
	"github.com/marcbinz/som/example/gen/som/where"
	"github.com/marcbinz/som/example/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gotest.tools/assert"
	"testing"
	"unicode/utf8"
)

const (
	containerName = "som_test_surrealdb"
)

func conf(endpoint string) som.Config {
	return som.Config{
		Address:   "ws://" + endpoint,
		Username:  "root",
		Password:  "root",
		Namespace: "som_test",
		Database:  "som_test",
	}
}

func TestQuery(t *testing.T) {
	client := &som.Client{}

	query := client.User().Query().
		Filter(
			where.User.String.In(nil),
			where.User.Bool.Is(true),
		)

	assert.Equal(t,
		"SELECT * FROM user WHERE (string INSIDE $0 AND bool == $1) ",
		query.Describe(),
	)

	query = query.Filter(
		where.Any(
			where.User.TimePtr.Nil(),
			where.User.UUID.Equal(uuid.New()),
		),
	)

	assert.Equal(t,
		"SELECT * FROM user WHERE (string INSIDE $0 AND bool == $1 AND (time_ptr == $2 OR uuid = $3)) ",
		query.Describe(),
	)
}

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

	client, err := som.NewClient(conf(endpoint))
	if err != nil {
		t.Error(err)
	}

	defer client.Close()

	if err := client.ApplySchema(); err != nil {
		t.Error(err)
	}

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

	client, err := som.NewClient(conf(endpoint))
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

func FuzzCustomModelIDs(f *testing.F) {
	fmt.Println("THIS IS A SINGLE RUN")

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

	client, err := som.NewClient(conf(endpoint))
	if err != nil {
		f.Error(err)
	}

	defer client.Close()

	f.Add("v9uitj942tv2403tnv")
	f.Add("vb92thj29v4tjn20d3")
	f.Add("ij024itvnjc20394it")
	f.Add(" 0")
	f.Add("\"0")
	f.Add("🙂")
	f.Add("✅")
	f.Add("👋😉")

	f.Fuzz(func(t *testing.T, id string) {
		if !utf8.ValidString(id) {
			t.Skip("id is not a valid utf8 string")
		}

		userIn := &model.User{
			String: "1",
		}

		userIn.ID = id

		err = client.User().Create(ctx, userIn)
		if err != nil {
			t.Error(err)
		}

		if userIn.ID == "" {
			t.Error("user ID must not be empty after create call")
		}

		userOut, ok, err := client.User().Read(ctx, userIn.ID)

		if err != nil {
			t.Error(err)
		}

		if !ok {
			t.Errorf("user with ID %s not found", userIn.ID)
		}

		assert.Equal(t, userIn.ID, userOut.ID)
		assert.Equal(t, "1", userOut.String)

		userOut.String = "2"

		err = client.User().Update(ctx, userOut)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "2", userOut.String)

		err = client.User().Delete(ctx, userOut)
		if err != nil {
			t.Error(err)
		}
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

	client, err := som.NewClient(conf(endpoint))
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
