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

	user := &model.User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = userRepo.Create(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	user2, err := userRepo.FindById(ctx, user.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user2)
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
