package model

import (
	"github.com/google/uuid"
	"github.com/marcbinz/som"
	"time"
)

type User struct {
	som.Node
	som.Timestamps

	String  string
	Int     int
	Int32   int32
	Int64   int64
	Float32 float32
	Float64 float64
	Bool    bool
	Bool2   bool

	Login        Login    // struct
	Role         Role     // enum
	Groups       []Group  // slice of Nodes
	MainGroup    Group    // node
	MainGroupPtr *Group   // node pointer
	Other        []string // slice of strings
	More         []float32
	Roles        []Role // slice of enum

	MemberOf []GroupMember

	StringPtr *string
	IntPtr    *int

	StructPtr         *SomeStruct
	StringPtrSlice    []*string
	StringSlicePtr    *[]string
	StructPtrSlice    []*SomeStruct
	StructPtrSlicePtr *[]*SomeStruct
	EnumPtrSlice      []*Role
	NodePtrSlice      []*Group
	NodePtrSlicePtr   *[]*Group
	SliceSlice        [][]string

	Time    time.Time
	TimePtr *time.Time

	Duration    time.Duration
	DurationPtr *time.Duration

	UUID    uuid.UUID
	UUIDPtr *uuid.UUID

	// MappedLogin  map[string]Login // map of string and struct
	// MappedRoles  map[string]Role  // map of string and enum
	// MappedGroups map[string]Group // map of string and node
	// OtherMap     map[Role]string  // map of enum and string
}

func (u *User) GetGroups() []Group {
	var nodes []Group
	for _, edge := range u.MemberOf {
		nodes = append(nodes, edge.Group)
	}
	return nodes
}

type Login struct {
	Username string
	Password string
}

type SomeStruct struct {
	StringPtr *string
	IntPtr    *int
	TimePtr   *time.Time
	UuidPtr   *uuid.UUID
}

type Role som.Enum

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Group struct {
	som.Node
	som.Timestamps

	Name string

	Members []GroupMember
}

func (g *Group) GetMembers() []User {
	var nodes []User
	for _, edge := range g.Members {
		nodes = append(nodes, edge.User)
	}
	return nodes
}

type GroupMember struct {
	som.Edge
	som.Timestamps

	User  User  `som:"in"`
	Group Group `som:"out"`

	Meta GroupMemberMeta
}

type GroupMemberMeta struct {
	IsAdmin  bool
	IsActive bool
}
