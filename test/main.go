package main

import (
	"context"
	"fmt"
	"github.com/marcbinz/sdb/example/gen/sdb"
	"github.com/marcbinz/sdb/example/model"
	"github.com/marcbinz/sdb/example/repo"
	"log"
	"time"
)

func main() {

	log.SetFlags(log.Lshortfile)

	ctx := context.Background()

	db, err := sdb.NewClient("ws://localhost:8010", "root", "root", "sdb", "default")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	groupRepo := repo.Group(db)
	userRepo := repo.User(db)

	group := &model.Group{
		Name: "some group",
	}

	err = groupRepo.Create(ctx, group)
	if err != nil {
		log.Fatal(err)
	}

	user := &model.User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		String:    "group:test",
		MainGroup: *group,
	}

	err = userRepo.Create(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	user2, err := userRepo.FindById(ctx, user.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user2", user2)

	users, err := userRepo.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		fmt.Println(user.ID, user.String, user.MainGroup.Name)
	}
}
