package model

import (
	"github.com/marcbinz/sdb"
	"time"
)

type User struct {
	sdb.Node     `surrealdb:"user"`
	ID           string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Login        Login
	Text         string
	Role         Role
	Groups       []Group
	Other        []string
	Roles        []Role
	MappedLogin  map[string]Login
	MappedRoles  map[string]Role
	MappedGroups map[string]Group
	OtherMap     map[Role]string
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
	Name     string
}
