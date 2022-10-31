package main

import (
	"context"
	"fmt"
	"github.com/marcbinz/sdb/example/gen/sdb"
	"github.com/marcbinz/sdb/example/model"
	"github.com/marcbinz/sdb/example/repo"
	"github.com/surrealdb/surrealdb.go"
	"log"
	"time"
)

func main() {
	ctx := context.Background()

	db, err := New("root", "root", "sdb", "default")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	userRepo := repo.User{DB: db}

	err = userRepo.Create(ctx, &model.User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	user, err := userRepo.FindById(ctx, "od5tfy9z8bzszl3p9l5l")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user)
}

func New(username, password, namespace, database string) (*sdb.Client, error) {
	db, err := surrealdb.New("ws://localhost:8010/rpc")
	if err != nil {
		return nil, fmt.Errorf("new failed: %v", err)
	}

	_, err = db.Signin(map[string]any{
		"user": username,
		"pass": password,
	})
	if err != nil {
		return nil, err
	}

	_, err = db.Use(namespace, database)
	if err != nil {
		return nil, err
	}

	return sdb.NewClient(db), nil
}
