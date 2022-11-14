package model

import (
	"github.com/google/uuid"
	"github.com/marcbinz/sdb"
	"time"
)

type User struct {
	sdb.Node `surrealdb:"user"`
	ID       string

	CreatedAt time.Time
	UpdatedAt time.Time

	String  string
	Int     int
	Int32   int32
	Int64   int64
	Float32 float32
	Float64 float64
	Bool    bool
	Bool2   bool

	UUID uuid.UUID

	Login     Login    // struct
	Role      Role     // enum
	Groups    []Group  // slice of Nodes
	MainGroup Group    // node
	Other     []string // slice of strings
	More      []float32
	Roles     []Role // slice of enum

	MyGroups []MemberOf

	// MappedLogin  map[string]Login // map of string and struct
	// MappedRoles  map[string]Role  // map of string and enum
	// MappedGroups map[string]Group // map of string and node
	// OtherMap     map[Role]string  // map of enum and string
}

func (u *User) GetGroups() []Group {
	var nodes []Group
	for _, edge := range u.MyGroups {
		nodes = append(nodes, edge.Group)
	}
	return nodes
}

type Login struct {
	Username string
	Password string
}

type Role sdb.Enum

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Group struct {
	sdb.Node `surrealdb:"group"`

	ID   string
	Name string

	Members []MemberOf
}

func (g *Group) GetMembers() []User {
	var nodes []User
	for _, edge := range g.Members {
		nodes = append(nodes, edge.User)
	}
	return nodes
}

type MemberOf struct {
	sdb.Edge

	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time

	User  User  `som:"in"`
	Group Group `som:"out"`

	Meta MemberOfMeta
}

type MemberOfMeta struct {
	IsAdmin  bool
	IsActive bool
}
