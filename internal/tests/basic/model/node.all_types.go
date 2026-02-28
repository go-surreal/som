package model

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som"
	gofrsuuid "github.com/gofrs/uuid"
	"github.com/google/uuid"
)

type AllTypes struct {
	som.Node[som.ULID]
	som.Timestamps

	// basic types

	FieldString            string     `som:"fulltext(english_search)"`
	FieldStringPtr         *string    `som:"fulltext(english_search)"`
	FieldOther             []string   `som:"fulltext(english_search)"`
	FieldStringPtrSlice    []*string  `som:"fulltext(english_search)"`
	FieldStringSlicePtr    *[]string  `som:"fulltext(english_search)"` // TODO: cannot be filtered for nil!
	FieldStringPtrSlicePtr *[]*string `som:"fulltext(english_search)"`

	FieldInt            int
	FieldIntPtr         *int
	FieldIntSlice       []int
	FieldIntPtrSlice    []*int
	FieldIntSlicePtr    *[]int
	FieldIntPtrSlicePtr *[]*int

	FieldInt8     int8 // -128 to 127
	FieldInt8Ptr  *int8
	FieldInt16    int16 // -2^15 to 2^15-1 (-32768 to 32767)
	FieldInt16Ptr *int16
	FieldInt32    int32 // -2^31 to 2^31-1 (-2147483648 to 2147483647)
	FieldInt32Ptr *int32
	FieldInt64    int64 // -2^63 to 2^63-1 (-9223372036854775808 to 9223372036854775807)
	FieldInt64Ptr *int64

	//Uint      uint
	//UintPtr   *uint
	FieldUint8     uint8 // 0 to 255
	FieldUint8Ptr  *uint8
	FieldUint16    uint16 // 0 to 2^16-1 (0 to 65535)
	FieldUint16Ptr *uint16
	FieldUint32    uint32 // 0 to 2^32-1 (0 to 4294967295)
	FieldUint32Ptr *uint32
	//Uint64    uint64 // 0 to 2^64-1 (0 to 18446744073709551615)
	//Uint64Ptr *uint64

	//Uintptr    uintptr
	//UintptrPtr *uintptr

	FieldFloat32            float32 // -3.4e+38 to 3.4e+38.
	FieldFloat32Slice       []float32
	FieldFloat32SlicePtr    *[]float32
	FieldFloat32PtrSlice    []*float32
	FieldFloat32PtrSlicePtr *[]*float32

	FieldFloat64 float64 // -1.7e+308 to +1.7e+308.

	// Complex64  complex64
	// Complex128 complex128

	FieldRune      rune
	FieldRuneSlice []rune

	FieldBool      bool
	FieldBoolPtr   *bool
	FieldBoolSlice []bool

	// TODO: should we support the any type? (surrealdb seems to support it)

	// TODO: support math types?
	// BigInt   big.Int
	// BigRat   big.Rat
	// BigFloat big.Float

	// special types

	FieldTime           time.Time
	FieldTimePtr        *time.Time
	FieldTimeNil        *time.Time
	FieldTimeSlice      []time.Time
	FieldTimeSliceSlice [][]time.Time

	FieldDuration      time.Duration
	FieldDurationPtr   *time.Duration
	FieldDurationNil   *time.Duration
	FieldDurationSlice []time.Duration

	FieldMonth    time.Month
	FieldMonthPtr *time.Month

	FieldWeekday    time.Weekday
	FieldWeekdayPtr *time.Weekday

	FieldUUID      uuid.UUID
	FieldUUIDPtr   *uuid.UUID
	FieldUUIDNil   *uuid.UUID
	FieldUUIDSlice []uuid.UUID

	FieldUUIDGofrs      gofrsuuid.UUID
	FieldUUIDGofrsPtr   *gofrsuuid.UUID
	FieldUUIDGofrsNil   *gofrsuuid.UUID
	FieldUUIDGofrsSlice []gofrsuuid.UUID

	FieldURL      url.URL
	FieldURLPtr   *url.URL
	FieldURLNil   *url.URL
	FieldURLSlice []url.URL

	FieldRegex      regexp.Regexp
	FieldRegexPtr   *regexp.Regexp
	FieldRegexSlice []regexp.Regexp

	FieldEmail      som.Email
	FieldEmailPtr   *som.Email
	FieldEmailNil   *som.Email
	FieldEmailSlice []som.Email

	// enums

	FieldEnum            Role
	FieldEnumPtr         *Role
	FieldEnumSlice       []Role
	FieldEnumPtrSlice    []*Role
	FieldEnumPtrSlicePtr *[]*Role

	// structs

	FieldCredentials             Credentials
	FieldNestedDataPtr           *NestedData
	FieldNestedDataSlice         []NestedData
	FieldNestedDataPtrSlice      []*NestedData
	FieldNestedDataPtrSlicePtr   *[]*NestedData

	// nodes

	FieldNode       SpecialTypes   // node
	FieldNodePtr    *SpecialTypes  // node pointer
	FieldNodeSlice      []SpecialTypes // slice of Nodes
	FieldNodeSliceSlice [][]SpecialTypes
	FieldNodePtrSlice    []*SpecialTypes
	FieldNodePtrSlicePtr *[]*SpecialTypes

	// edges

	FieldEdgeRelations []EdgeRelation // slice of edges

	// other

	FieldSliceSlice       [][]string
	FieldSliceSliceSlice  [][][]string
	FieldSliceSliceSlice2 [][][]NestedData

	// maps (not (yet?) supported)

	FieldByte         byte
	FieldBytePtr      *byte
	FieldByteSlice    []byte
	FieldByteSlicePtr *[]byte // TODO: cannot be filtered for nil!

	//// MappedLogin  map[string]Credentials // map of string and struct
	//// MappedRoles  map[string]Role  // map of string and enum
	//// MappedGroups map[string]Group // map of string and node
	//// OtherMap     map[Role]string  // map of enum and string

	// hook fields
	FieldHookStatus string
	FieldHookDetail string
}

type contextKey string

const AbortDeleteKey contextKey = "abortDelete"
const AfterDeleteCalledKey contextKey = "afterDeleteCalled"

func (f *AllTypes) OnBeforeCreate(_ context.Context) error {
	f.FieldHookStatus = "[created]" + f.FieldHookStatus
	return nil
}

func (f *AllTypes) OnAfterCreate(_ context.Context) error {
	f.FieldHookDetail = f.FieldHookDetail + "[after-create]"
	return nil
}

func (f *AllTypes) OnBeforeUpdate(_ context.Context) error {
	f.FieldHookStatus = "[updated]" + f.FieldHookStatus
	return nil
}

func (f *AllTypes) OnAfterUpdate(_ context.Context) error {
	f.FieldHookDetail = f.FieldHookDetail + "[after-update]"
	return nil
}

func (f *AllTypes) OnBeforeDelete(ctx context.Context) error {
	if ptr, ok := ctx.Value(AbortDeleteKey).(*bool); ok && *ptr {
		return errors.New("delete aborted by model hook")
	}
	return nil
}

func (f *AllTypes) OnAfterDelete(ctx context.Context) error {
	if ptr, ok := ctx.Value(AfterDeleteCalledKey).(*bool); ok {
		*ptr = true
	}
	return nil
}

type Credentials struct {
	Username string `som:"index"`

	Password    som.Password[som.Bcrypt]
	PasswordPtr *som.Password[som.Argon2]
}

type NestedData struct {
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

type EdgeRelation struct {
	som.Edge
	som.Timestamps

	AllTypes     AllTypes      `som:"in"`
	SpecialTypes SpecialTypes  `som:"out"`

	Meta EdgeMeta
}

type EdgeMeta struct {
	IsAdmin  bool
	IsActive bool
}
