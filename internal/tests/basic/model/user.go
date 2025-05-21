package model

import (
	"github.com/go-surreal/som"
	sdir "github.com/go-surreal/som/exp/model/subdir"
)

type User struct {
	som.Node
	//som.Timestamps

	x, y func()

	dafuq struct{ bla string }

	//Name string
	//Age  int
	//
	//UUID uuid.UUID
	//
	StringPtr *string
	//
	//Login        Login    // struct
	//Role         Role     // enum
	//Groups       []Group  // slice of Nodes
	//MainGroup    Group    // node
	//MainGroupPtr *Group   // node pointer
	Other []string // slice of strings
	//More         []float32
	//Roles        []Role // slice of enum

	X sdir.X
}

type Some int
type _ string
type D[T any] bool

type a int
type B a
type V sdir.X
type user2 User

func X() {
	// this is a function declaration
	// it should not be parsed
}

// this is a bad declaration
// it should not be parsed
//type P int

//type Login struct {
//	Username string
//	Password string
//}
//
//type Role som.Enum
//
//const (
//	RoleUser  Role = "user"
//	RoleAdmin Role = "admin"
//)
//
//type Group struct {
//	som.Node
//	som.Timestamps
//
//	Name string
//}
