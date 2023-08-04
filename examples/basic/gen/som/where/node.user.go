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
		Bool:       lib.NewBool[T](lib.Field(key, "bool")),
		BoolPtr:    lib.NewBoolPtr[T](lib.Field(key, "bool_ptr")),
		CreatedAt:  lib.NewTime[T](lib.Field(key, "created_at")),
		Float32:    lib.NewNumeric[float32, T](lib.Field(key, "float_32")),
		Float64:    lib.NewNumeric[float64, T](lib.Field(key, "float_64")),
		ID:         lib.NewID[T](lib.Field(key, "id"), "user"),
		Int:        lib.NewNumeric[int, T](lib.Field(key, "int")),
		Int16:      lib.NewNumeric[int16, T](lib.Field(key, "int_16")),
		Int16Ptr:   lib.NewNumericPtr[*int16, T](lib.Field(key, "int_16_ptr")),
		Int32:      lib.NewNumeric[int32, T](lib.Field(key, "int_32")),
		Int32Ptr:   lib.NewNumericPtr[*int32, T](lib.Field(key, "int_32_ptr")),
		Int64:      lib.NewNumeric[int64, T](lib.Field(key, "int_64")),
		Int64Ptr:   lib.NewNumericPtr[*int64, T](lib.Field(key, "int_64_ptr")),
		Int8:       lib.NewNumeric[int8, T](lib.Field(key, "int_8")),
		Int8Ptr:    lib.NewNumericPtr[*int8, T](lib.Field(key, "int_8_ptr")),
		IntPtr:     lib.NewNumericPtr[*int, T](lib.Field(key, "int_ptr")),
		Role:       lib.NewBase[model.Role, T](lib.Field(key, "role")),
		Rune:       lib.NewNumeric[rune, T](lib.Field(key, "rune")),
		String:     lib.NewString[T](lib.Field(key, "string")),
		StringPtr:  lib.NewStringPtr[T](lib.Field(key, "string_ptr")),
		TimePtr:    lib.NewTimePtr[T](lib.Field(key, "time_ptr")),
		UUID:       lib.NewBase[uuid.UUID, T](lib.Field(key, "uuid")),
		Uint:       lib.NewNumeric[uint, T](lib.Field(key, "uint")),
		Uint16:     lib.NewNumeric[uint16, T](lib.Field(key, "uint_16")),
		Uint16Ptr:  lib.NewNumericPtr[*uint16, T](lib.Field(key, "uint_16_ptr")),
		Uint32:     lib.NewNumeric[uint32, T](lib.Field(key, "uint_32")),
		Uint32Ptr:  lib.NewNumericPtr[*uint32, T](lib.Field(key, "uint_32_ptr")),
		Uint64:     lib.NewNumeric[uint64, T](lib.Field(key, "uint_64")),
		Uint64Ptr:  lib.NewNumericPtr[*uint64, T](lib.Field(key, "uint_64_ptr")),
		Uint8:      lib.NewNumeric[uint8, T](lib.Field(key, "uint_8")),
		Uint8Ptr:   lib.NewNumericPtr[*uint8, T](lib.Field(key, "uint_8_ptr")),
		UintPtr:    lib.NewNumericPtr[*uint, T](lib.Field(key, "uint_ptr")),
		Uintptr:    lib.NewNumeric[uintptr, T](lib.Field(key, "uintptr")),
		UintptrPtr: lib.NewNumericPtr[*uintptr, T](lib.Field(key, "uintptr_ptr")),
		UpdatedAt:  lib.NewTime[T](lib.Field(key, "updated_at")),
		UuidPtr:    lib.NewBasePtr[uuid.UUID, T](lib.Field(key, "uuid_ptr")),
		key:        key,
	}
}

type user[T any] struct {
	key        lib.Key[T]
	ID         *lib.ID[T]
	CreatedAt  *lib.Time[T]
	UpdatedAt  *lib.Time[T]
	String     *lib.String[T]
	StringPtr  *lib.StringPtr[T]
	Int        *lib.Numeric[int, T]
	IntPtr     *lib.NumericPtr[*int, T]
	Int8       *lib.Numeric[int8, T]
	Int8Ptr    *lib.NumericPtr[*int8, T]
	Int16      *lib.Numeric[int16, T]
	Int16Ptr   *lib.NumericPtr[*int16, T]
	Int32      *lib.Numeric[int32, T]
	Int32Ptr   *lib.NumericPtr[*int32, T]
	Int64      *lib.Numeric[int64, T]
	Int64Ptr   *lib.NumericPtr[*int64, T]
	Uint       *lib.Numeric[uint, T]
	UintPtr    *lib.NumericPtr[*uint, T]
	Uint8      *lib.Numeric[uint8, T]
	Uint8Ptr   *lib.NumericPtr[*uint8, T]
	Uint16     *lib.Numeric[uint16, T]
	Uint16Ptr  *lib.NumericPtr[*uint16, T]
	Uint32     *lib.Numeric[uint32, T]
	Uint32Ptr  *lib.NumericPtr[*uint32, T]
	Uint64     *lib.Numeric[uint64, T]
	Uint64Ptr  *lib.NumericPtr[*uint64, T]
	Uintptr    *lib.Numeric[uintptr, T]
	UintptrPtr *lib.NumericPtr[*uintptr, T]
	Float32    *lib.Numeric[float32, T]
	Float64    *lib.Numeric[float64, T]
	Rune       *lib.Numeric[rune, T]
	Bool       *lib.Bool[T]
	BoolPtr    *lib.BoolPtr[T]
	UUID       *lib.Base[uuid.UUID, T]
	Role       *lib.Base[model.Role, T]
	TimePtr    *lib.TimePtr[T]
	UuidPtr    *lib.BasePtr[uuid.UUID, T]
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
