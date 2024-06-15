package model

import "github.com/marcbinz/som"

type User struct {
	som.Node
	//som.Timestamps

	x, y func()

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
}

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
