package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/marcbinz/som/example/gen/som"
	"github.com/marcbinz/som/example/model"
	"github.com/marcbinz/som/example/repo"
	"log"
	"time"
)

func main() {

	log.SetFlags(log.Lshortfile)

	ctx := context.Background()

	db, err := som.NewClient("ws://localhost:8010", "root", "root", "som", "default")
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

	fmt.Println("group:", group.ID, group.Name)

	user := &model.User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		String:    "Marc",
		MainGroup: *group,
	}

	err = userRepo.Create(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user:", user.ID, user.String)

	edge := &model.MemberOf{
		CreatedAt: time.Now(),
		User:      *user,
		Group:     *group,
		Meta: model.MemberOfMeta{
			IsAdmin: true,
		},
	}

	err = userRepo.Relate(ctx, edge)
	if err != nil {
		log.Fatal(err)
	}

	user, ok, err := userRepo.Read(ctx, user.ID)
	if err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Fatal("could not find user with id:", user.ID)
	}

	// fmt.Println("relation:", edge.ID)
	// fmt.Println("user:", edge.User.ID)
	// fmt.Println("group:", edge.Group.ID)

	fmt.Println("old user uuid:", user.UUID, user.ID)

	user.UUID, _ = uuid.NewUUID()

	fmt.Println("new user uuid:", user.UUID, user.ID)

	err = userRepo.Update(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("updated user uuid:", user.UUID, user.ID)

	err = userRepo.Delete(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("deleted user:", user.ID)

	//
	// user2, err := userRepo.FindById(ctx, user.ID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// fmt.Println("user2", user2)

	// users, err := userRepo.List(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// for _, user := range users {
	// 	fmt.Println("result:", user.ID, user.MainGroup.Name)
	// }
}
