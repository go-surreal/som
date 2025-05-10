package model

import (
	"github.com/go-surreal/som/tests/basic/gen/som/sombase"
	"github.com/google/uuid"
	"net/url"
	"time"
)

type AllFieldTypes struct {
	sombase.Node
	sombase.Timestamps

	// basic types

	String         string
	StringPtr      *string
	Other          []string
	StringPtrSlice []*string
	StringSlicePtr *[]string // TODO: cannot be filtered for nil!

	Int            int
	IntPtr         *int
	IntSlice       []int
	IntPtrSlice    []*int
	IntSlicePtr    *[]int
	IntPtrSlicePtr *[]*int

	Int8     int8 // -128 to 127
	Int8Ptr  *int8
	Int16    int16 // -2^15 to 2^15-1 (-32768 to 32767)
	Int16Ptr *int16
	Int32    int32 // -2^31 to 2^31-1 (-2147483648 to 2147483647)
	Int32Ptr *int32
	Int64    int64 // -2^63 to 2^63-1 (-9223372036854775808 to 9223372036854775807)
	Int64Ptr *int64

	//Uint      uint
	//UintPtr   *uint
	Uint8     uint8 // 0 to 255
	Uint8Ptr  *uint8
	Uint16    uint16 // 0 to 2^16-1 (0 to 65535)
	Uint16Ptr *uint16
	Uint32    uint32 // 0 to 2^32-1 (0 to 4294967295)
	Uint32Ptr *uint32
	//Uint64    uint64 // 0 to 2^64-1 (0 to 18446744073709551615)
	//Uint64Ptr *uint64

	//Uintptr    uintptr
	//UintptrPtr *uintptr

	Float32            float32 // -3.4e+38 to 3.4e+38.
	Float32Slice       []float32
	Float32SlicePtr    *[]float32
	Float32PtrSlice    []*float32
	Float32PtrSlicePtr *[]*float32

	Float64 float64 // -1.7e+308 to +1.7e+308.

	// Complex64  complex64
	// Complex128 complex128

	Rune      rune
	RuneSlice []rune

	Bool      bool
	BoolPtr   *bool
	BoolSlice []bool

	// TODO: should we support the any type? (surrealdb seems to support it)

	// TODO: support math types?
	// BigInt   big.Int
	// BigRat   big.Rat
	// BigFloat big.Float

	// special types

	Time           time.Time
	TimePtr        *time.Time
	TimeNil        *time.Time
	TimeSlice      []time.Time
	TimeSliceSlice [][]time.Time

	Duration      time.Duration
	DurationPtr   *time.Duration
	DurationNil   *time.Duration
	DurationSlice []time.Duration

	UUID      uuid.UUID
	UUIDPtr   *uuid.UUID
	UUIDNil   *uuid.UUID
	UUIDSlice []uuid.UUID

	URL      url.URL
	URLPtr   *url.URL
	URLNil   *url.URL
	URLSlice []url.URL

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
	GroupsSlice     [][]Group
	NodePtrSlice    []*Group
	NodePtrSlicePtr *[]*Group

	// edges

	MemberOf []GroupMember // slice of edges

	// other

	SliceSlice       [][]string
	SliceSliceSlice  [][][]string
	SliceSliceSlice2 [][][]SomeStruct

	// maps (not (yet?) supported)

	Byte         byte
	BytePtr      *byte
	ByteSlice    []byte
	ByteSlicePtr *[]byte // TODO: cannot be filtered for nil!

	//// MappedLogin  map[string]Login // map of string and struct
	//// MappedRoles  map[string]Role  // map of string and enum
	//// MappedGroups map[string]Group // map of string and node
	//// OtherMap     map[Role]string  // map of enum and string
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

type Role sombase.Enum

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Group struct {
	sombase.Node
	sombase.Timestamps

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
	sombase.Edge
	sombase.Timestamps

	User  AllFieldTypes `som:"in"`
	Group Group         `som:"out"`

	Meta GroupMemberMeta
}

type GroupMemberMeta struct {
	IsAdmin  bool
	IsActive bool
}
