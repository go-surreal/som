// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package where

import (
	uuid "github.com/google/uuid"
	lib "github.com/marcbinz/som/examples/basic/gen/som/internal/lib"
	model "github.com/marcbinz/som/examples/basic/model"
)

var User = newUser[model.User](lib.NewKey[model.User]())

func newUser[T any](key lib.Key[T]) user[T] {
	return user[T]{
		Bool:        lib.NewBool[T](lib.Field(key, "bool")),
		Bool2:       lib.NewBool[T](lib.Field(key, "bool_2")),
		CreatedAt:   lib.NewTime[T](lib.Field(key, "created_at")),
		Duration:    lib.NewDuration[T](lib.Field(key, "duration")),
		DurationPtr: lib.NewDurationPtr[T](lib.Field(key, "duration_ptr")),
		Float32:     lib.NewNumeric[float32, T](lib.Field(key, "float_32")),
		Float64:     lib.NewNumeric[float64, T](lib.Field(key, "float_64")),
		ID:          lib.NewID[T](lib.Field(key, "id"), "user"),
		Int:         lib.NewNumeric[int, T](lib.Field(key, "int")),
		Int32:       lib.NewNumeric[int32, T](lib.Field(key, "int_32")),
		Int64:       lib.NewNumeric[int64, T](lib.Field(key, "int_64")),
		IntPtr:      lib.NewNumericPtr[*int, T](lib.Field(key, "int_ptr")),
		Role:        lib.NewBase[model.Role, T](lib.Field(key, "role")),
		String:      lib.NewString[T](lib.Field(key, "string")),
		StringPtr:   lib.NewStringPtr[T](lib.Field(key, "string_ptr")),
		Time:        lib.NewTime[T](lib.Field(key, "time")),
		TimePtr:     lib.NewTimePtr[T](lib.Field(key, "time_ptr")),
		UUID:        lib.NewBase[uuid.UUID, T](lib.Field(key, "uuid")),
		UUIDPtr:     lib.NewBasePtr[uuid.UUID, T](lib.Field(key, "uuid_ptr")),
		UpdatedAt:   lib.NewTime[T](lib.Field(key, "updated_at")),
		key:         key,
	}
}

type user[T any] struct {
	key         lib.Key[T]
	ID          *lib.ID[T]
	CreatedAt   *lib.Time[T]
	UpdatedAt   *lib.Time[T]
	String      *lib.String[T]
	Int         *lib.Numeric[int, T]
	Int32       *lib.Numeric[int32, T]
	Int64       *lib.Numeric[int64, T]
	Float32     *lib.Numeric[float32, T]
	Float64     *lib.Numeric[float64, T]
	Bool        *lib.Bool[T]
	Bool2       *lib.Bool[T]
	Role        *lib.Base[model.Role, T]
	StringPtr   *lib.StringPtr[T]
	IntPtr      *lib.NumericPtr[*int, T]
	Time        *lib.Time[T]
	TimePtr     *lib.TimePtr[T]
	Duration    *lib.Duration[T]
	DurationPtr *lib.DurationPtr[T]
	UUID        *lib.Base[uuid.UUID, T]
	UUIDPtr     *lib.BasePtr[uuid.UUID, T]
}

func (n user[T]) Login() login[T] {
	return newLogin[T](lib.Field(n.key, "login"))
}

func (n user[T]) Groups(filters ...lib.Filter[model.Group]) groupSlice[T] {
	key := lib.Node(n.key, "groups", filters)
	return groupSlice[T]{lib.KeyFilter[T](key), lib.NewSlice[T, model.Group](key)}
}

func (n user[T]) MainGroup() group[T] {
	return newGroup[T](lib.Field(n.key, "main_group"))
}

func (n user[T]) MainGroupPtr() group[T] {
	return newGroup[T](lib.Field(n.key, "main_group_ptr"))
}

func (n user[T]) Other() *lib.Slice[T, string] {
	return lib.NewSlice[T, string](lib.Field(n.key, "other"))
}

func (n user[T]) More() *lib.Slice[T, float32] {
	return lib.NewSlice[T, float32](lib.Field(n.key, "more"))
}

func (n user[T]) Roles() *lib.Slice[T, model.Role] {
	return lib.NewSlice[T, model.Role](lib.Field(n.key, "roles"))
}

func (n user[T]) MemberOf(filters ...lib.Filter[model.GroupMember]) groupMemberIn[T] {
	return newGroupMemberIn[T](lib.EdgeIn(n.key, "group_member", filters))
}

func (n user[T]) StructPtr() someStruct[T] {
	return newSomeStruct[T](lib.Field(n.key, "struct_ptr"))
}

func (n user[T]) StringPtrSlice() *lib.Slice[T, *string] {
	return lib.NewSlice[T, *string](lib.Field(n.key, "string_ptr_slice"))
}

func (n user[T]) StringSlicePtr() *lib.Slice[T, string] {
	return lib.NewSlice[T, string](lib.Field(n.key, "string_slice_ptr"))
}

func (n user[T]) StructPtrSlice() *lib.Slice[T, *model.SomeStruct] {
	return lib.NewSlice[T, *model.SomeStruct](lib.Field(n.key, "struct_ptr_slice"))
}

func (n user[T]) StructPtrSlicePtr() *lib.Slice[T, *model.SomeStruct] {
	return lib.NewSlice[T, *model.SomeStruct](lib.Field(n.key, "struct_ptr_slice_ptr"))
}

func (n user[T]) EnumPtrSlice() *lib.Slice[T, model.Role] {
	return lib.NewSlice[T, model.Role](lib.Field(n.key, "enum_ptr_slice"))
}

func (n user[T]) NodePtrSlice(filters ...lib.Filter[model.Group]) groupSlice[T] {
	key := lib.Node(n.key, "node_ptr_slice", filters)
	return groupSlice[T]{lib.KeyFilter[T](key), lib.NewSlice[T, model.Group](key)}
}

func (n user[T]) NodePtrSlicePtr(filters ...lib.Filter[model.Group]) groupSlice[T] {
	key := lib.Node(n.key, "node_ptr_slice_ptr", filters)
	return groupSlice[T]{lib.KeyFilter[T](key), lib.NewSlice[T, model.Group](key)}
}

func (n user[T]) SliceSlice() *lib.Slice[T, []string] {
	return lib.NewSlice[T, []string](lib.Field(n.key, "slice_slice"))
}

type userEdges[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

func (n userEdges[T]) MemberOf(filters ...lib.Filter[model.GroupMember]) groupMemberIn[T] {
	return newGroupMemberIn[T](lib.EdgeIn(n.key, "group_member", filters))
}

type userSlice[T any] struct {
	lib.Filter[T]
	*lib.Slice[T, model.User]
}
