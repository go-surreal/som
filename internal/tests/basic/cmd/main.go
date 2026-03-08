package main

import (
	"context"

	"github.com/go-surreal/som/tests/basic/gen/som/repo"
)

func main() {
	client, err := repo.NewClient(context.Background(), repo.Config{})
	if err != nil {
		panic(err)
	}

	_ = client
}