package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/marcbinz/som/example/gen/som"
	"github.com/marcbinz/som/example/gen/som/where"
	"github.com/marcbinz/som/example/model"
	"log"
)

func main() {

	log.SetFlags(log.Lshortfile)

	ctx := context.Background()

	db, err := som.NewClient(som.Config{
		Address:   "ws://localhost:8010",
		Username:  "root",
		Password:  "root",
		Namespace: "som",
		Database:  "default",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.ApplySchema(); err != nil {
		log.Fatal(err)
	}

	group := &model.Group{
		Name: "some group",
	}

	err = db.Group().Create(ctx, group)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("group:", group.ID(), group.Name, group.CreatedAt())

	user := &model.User{
		String:    "Marc",
		MainGroup: *group,
	}

	err = db.User().Create(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user:", user.ID(), user.String, user.CreatedAt(), user.UpdatedAt().IsZero())

	edge := &model.MemberOf{
		User:  *user,
		Group: *group,
		Meta: model.MemberOfMeta{
			IsAdmin: true,
		},
	}

	err = db.User().Relate().MyGroups().Create(edge)
	if err != nil {
		log.Fatal(err)
	}

	user, ok, err := db.User().Read(ctx, user.ID())
	if err != nil {
		log.Fatal(err)
	}

	if !ok {
		log.Fatal("could not find user with id:", user.ID())
	}

	// fmt.Println("relation:", edge.ID)
	// fmt.Println("user:", edge.User.ID)
	// fmt.Println("group:", edge.Group.ID)

	fmt.Println("old user uuid:", user.UUID, user.ID(), user.UpdatedAt())

	user.UUID = uuid.New()

	value := "some value"
	user.StringPtr = &value

	fmt.Println("new user uuid:", user.UUID, user.ID())

	err = db.User().Update(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("updated user uuid:", user.UUID, user.ID(), user.UpdatedAt())

	query := db.User().Query().
		Filter(
			where.User.ID.NotEqual(""),
		)

	exists, err := query.
		Filter(where.User.StringPtr.NotNil()).
		Exists()

	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		log.Fatal("record with given filter does not exist")
	}

	user.StringPtr = nil

	err = db.User().Update(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	exists, err = query.
		Filter(where.User.StringPtr.Nil()).
		Exists()

	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		log.Fatal("record with given filter does not exist")
	}

	err = db.User().Delete(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("deleted user:", user.ID())

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
