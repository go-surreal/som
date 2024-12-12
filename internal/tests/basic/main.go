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
	surrealDBVersion    = "2.1.3"
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
