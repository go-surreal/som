package main

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/docker/docker/api/types/container"
	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
	"os"
)

const (
	surrealDBVersion    = "2.0.0-beta.1"
	containerStartedMsg = "Started web server on "
)

func main() {
	ctx := context.Background()

	client, cleanup, err := prepareDatabase(ctx, "main")
	if err != nil {
		panic(err)
	}

	defer cleanup()

	_, err = client.AllFieldTypesRepo().Query().All(ctx)
	if err != nil {
		panic(err)
	}
}

func prepareDatabase(ctx context.Context, name string) (som.Client, func(), error) {
	username := gofakeit.Username()
	password := gofakeit.Password(true, true, true, true, true, 32)
	namespace := gofakeit.FirstName()
	database := gofakeit.LastName()

	req := testcontainers.ContainerRequest{
		Name:  "sdbc_" + name,
		Image: "surrealdb/surrealdb:v" + surrealDBVersion,
		Env: map[string]string{
			"SURREAL_PATH":   "memory",
			"SURREAL_STRICT": "true",
			"SURREAL_USER":   username,
			"SURREAL_PASS":   password,
		},
		Cmd: []string{
			"start", "--allow-funcs", "--log", "trace",
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
		return nil, nil, err
	}

	endpoint, err := surreal.Endpoint(ctx, "")
	if err != nil {
		return nil, nil, err
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	config := som.Config{
		Host:      endpoint,
		Username:  username,
		Password:  password,
		Namespace: namespace,
		Database:  database,
	}

	opts := []som.Option{
		som.WithLogger(slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)),
	}

	client, err := som.NewClient(ctx, config, opts...)
	if err != nil {
		return nil, nil, err
	}

	if err := client.ApplySchema(ctx); err != nil {
		return nil, nil, err
	}

	//err = client.Execute(ctx,
	//	"DEFINE INDEX OVERWRITE test ON "+client.TableAllFieldTypes()+" FIELDS string, string_ptr CONCURRENTLY;", // or: IF NOT EXISTS
	//	map[string]any{}, // INFO FOR TABLE all_field_types; shows the status of the index
	//) // REBUILD
	//if err != nil {
	//	tb.Fatal(err)
	//}

	cleanup := func() {
		client.Close()

		if err := surreal.Terminate(ctx); err != nil {
			slog.ErrorContext(ctx, "Failed to terminate container.",
				"error", err.Error(),
			)
		}
	}

	return client, cleanup, nil
}
