package model

import (
	"github.com/go-surreal/som"
	"github.com/google/uuid"
	"time"
)

type AllFieldTypes struct {
	som.Node
	som.Timestamps

	// basic types

	String         string
	StringPtr      *string
	Other          []string
	StringPtrSlice []*string
	StringSlicePtr *[]string // TODO: cannot be filtered for nil!

	Int    int
	IntPtr *int

	Int32 int32
	Int64 int64

	Float32 float32
	More    []float32

	Float64 float64
	Bool    bool
	Bool2   bool

	// special types

	Time    time.Time
	TimePtr *time.Time

	UUID    uuid.UUID
	UUIDPtr *uuid.UUID

	// enums

	Role            Role
	EnumPtr         *Role
	Roles           []Role
	EnumPtrSlice    []*Role
	EnumPtrSlicePtr *[]*Role

	// structs

	Login             Login
	StructPtr         *SomeStruct
	StructSlice       []SomeStruct
	StructPtrSlice    []*SomeStruct
	StructPtrSlicePtr *[]*SomeStruct

	// nodes

	MainGroup       Group   // node
	MainGroupPtr    *Group  // node pointer
	Groups          []Group // slice of Nodes
	NodePtrSlice    []*Group
	NodePtrSlicePtr *[]*Group

	// edges

	MemberOf []GroupMember

	// other

	SliceSlice [][]string

	// maps (not (yet?) supported)

	Byte         byte
	BytePtr      *byte
	ByteSlice    []byte
	ByteSlicePtr *[]byte // TODO: cannot be filtered for nil!

	// MappedLogin  map[string]Login // map of string and struct
	// MappedRoles  map[string]Role  // map of string and enum
	// MappedGroups map[string]Group // map of string and node
	// OtherMap     map[Role]string  // map of enum and string
}

func (u *AllFieldTypes) GetGroups() []Group {
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

func (g *Group) GetMembers() []AllFieldTypes {
	var nodes []AllFieldTypes
	for _, edge := range g.Members {
		nodes = append(nodes, edge.User)
	}
	return nodes
}

type GroupMember struct {
	som.Edge
	som.Timestamps

	User  AllFieldTypes `som:"in"`
	Group Group         `som:"out"`

	Meta GroupMemberMeta
}

type GroupMemberMeta struct {
	IsAdmin  bool
	IsActive bool
}
