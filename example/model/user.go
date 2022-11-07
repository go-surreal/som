package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/marcbinz/sdb"
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
	Roles     []Role   // slice of enum

	// MappedLogin  map[string]Login // map of string and struct
	// MappedRoles  map[string]Role  // map of string and enum
	// MappedGroups map[string]Group // map of string and node
	// OtherMap     map[Role]string  // map of enum and string
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
}
