// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package where

import (
	uuid "github.com/google/uuid"
	model "github.com/marcbinz/som/example/model"
	lib "github.com/marcbinz/som/lib"
)

var User = newUser[model.User](lib.NewKey())

func newUser[T any](key lib.Key) user[T] {
	return user[T]{
		Bool:      lib.NewBool[T](key.Field("bool")),
		Bool2:     lib.NewBool[T](key.Field("bool_2")),
		CreatedAt: lib.NewTime[T](key.Field("created_at")),
		Float32:   lib.NewNumeric[float32, T](key.Field("float_32")),
		Float64:   lib.NewNumeric[float64, T](key.Field("float_64")),
		ID:        lib.NewID[T](key.Field("id"), "user"),
		Int:       lib.NewNumeric[int, T](key.Field("int")),
		Int32:     lib.NewNumeric[int32, T](key.Field("int_32")),
		Int64:     lib.NewNumeric[int64, T](key.Field("int_64")),
		IntPtr:    lib.NewNumericPtr[*int, T](key.Field("int_ptr")),
		Role:      lib.NewBase[model.Role, T](key.Field("role")),
		String:    lib.NewString[T](key.Field("string")),
		StringPtr: lib.NewStringPtr[T](key.Field("string_ptr")),
		TimePtr:   lib.NewTimePtr[T](key.Field("time_ptr")),
		UUID:      lib.NewBase[uuid.UUID, T](key.Field("uuid")),
		UpdatedAt: lib.NewTime[T](key.Field("updated_at")),
		UuidPtr:   lib.NewBasePtr[uuid.UUID, T](key.Field("uuid_ptr")),
		key:       key,
	}
}

type user[T any] struct {
	key       lib.Key
	ID        *lib.ID[T]
	CreatedAt *lib.Time[T]
	UpdatedAt *lib.Time[T]
	String    *lib.String[T]
	Int       *lib.Numeric[int, T]
	Int32     *lib.Numeric[int32, T]
	Int64     *lib.Numeric[int64, T]
	Float32   *lib.Numeric[float32, T]
	Float64   *lib.Numeric[float64, T]
	Bool      *lib.Bool[T]
	Bool2     *lib.Bool[T]
	UUID      *lib.Base[uuid.UUID, T]
	Role      *lib.Base[model.Role, T]
	StringPtr *lib.StringPtr[T]
	IntPtr    *lib.NumericPtr[*int, T]
	TimePtr   *lib.TimePtr[T]
	UuidPtr   *lib.BasePtr[uuid.UUID, T]
}
type userSlice[T any] struct {
	user[T]
	*lib.Slice[T, model.User]
}

func (n user[T]) Login() login[T] {
	return newLogin[T](n.key.Field("login"))
}
func (n user[T]) Groups(filters ...lib.Filter[model.Group]) groupSlice[T] {
	key := n.key.Node("groups", lib.Filters(filters))
	return groupSlice[T]{newGroup[T](key), lib.NewSlice[T, model.Group](key)}
}
func (n user[T]) MainGroup() group[T] {
	return newGroup[T](n.key.Field("main_group"))
}
func (n user[T]) Other() *lib.Slice[T, string] {
	return lib.NewSlice[T, string](n.key.Field("other"))
}
func (n user[T]) More() *lib.Slice[T, float32] {
	return lib.NewSlice[T, float32](n.key.Field("more"))
}
func (n user[T]) Roles() *lib.Slice[T, model.Role] {
	return lib.NewSlice[T, model.Role](n.key.Field("roles"))
}
func (n user[T]) MyGroups(filters ...lib.Filter[model.MemberOf]) memberOfIn[T] {
	return newMemberOfIn[T](n.key.EdgeIn("member_of", lib.Filters(filters)))
}
func (n user[T]) StructPtr() someStruct[T] {
	return newSomeStruct[T](n.key.Field("struct_ptr"))
}
func (n user[T]) StringPtrSlice() *lib.Slice[T, *string] {
	return lib.NewSlice[T, *string](n.key.Field("string_ptr_slice"))
}
func (n user[T]) StringSlicePtr() *lib.Slice[T, string] {
	return lib.NewSlice[T, string](n.key.Field("string_slice_ptr"))
}
func (n user[T]) StructPtrSlice() *lib.Slice[T, *model.SomeStruct] {
	return lib.NewSlice[T, *model.SomeStruct](n.key.Field("struct_ptr_slice"))
}
func (n user[T]) StructPtrSlicePtr() *lib.Slice[T, *model.SomeStruct] {
	return lib.NewSlice[T, *model.SomeStruct](n.key.Field("struct_ptr_slice_ptr"))
}
func (n user[T]) EnumPtrSlice() *lib.Slice[T, model.Role] {
	return lib.NewSlice[T, model.Role](n.key.Field("enum_ptr_slice"))
}
func (n user[T]) NodePtrSlice(filters ...lib.Filter[model.Group]) groupSlice[T] {
	key := n.key.Node("node_ptr_slice", lib.Filters(filters))
	return groupSlice[T]{newGroup[T](key), lib.NewSlice[T, model.Group](key)}
}
func (n user[T]) NodePtrSlicePtr(filters ...lib.Filter[model.Group]) groupSlice[T] {
	key := n.key.Node("node_ptr_slice_ptr", lib.Filters(filters))
	return groupSlice[T]{newGroup[T](key), lib.NewSlice[T, model.Group](key)}
}
func (n user[T]) SliceSlice() *lib.Slice[T, []string] {
	return lib.NewSlice[T, []string](n.key.Field("slice_slice"))
}
