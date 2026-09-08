package main

import (
	"context"

	"som.test/gen/som/repo"
)

func main() {
	client, err := repo.NewClient(context.Background(), repo.Config{})
	if err != nil {
		panic(err)
	}

	_ = client
}
