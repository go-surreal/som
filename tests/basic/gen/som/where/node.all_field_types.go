// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import (
	lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/tests/basic/model"
	uuid "github.com/google/uuid"
)

var AllFieldTypes = newAllFieldTypes[model.AllFieldTypes](lib.NewKey[model.AllFieldTypes]())

func newAllFieldTypes[T any](key lib.Key[T]) allFieldTypes[T] {
	return allFieldTypes[T]{
		Bool:      lib.NewBool[T](lib.Field(key, "bool")),
		BoolPtr:   lib.NewBoolPtr[T](lib.Field(key, "bool_ptr")),
		CreatedAt: lib.NewTime[T](lib.Field(key, "created_at")),
		EnumPtr:   lib.NewBasePtr[model.Role, T](lib.Field(key, "enum_ptr")),
		Float32:   lib.NewNumeric[float32, T](lib.Field(key, "float_32")),
		Float64:   lib.NewNumeric[float64, T](lib.Field(key, "float_64")),
		ID:        lib.NewID[T](lib.Field(key, "id"), "all_field_types"),
		Int:       lib.NewNumeric[int, T](lib.Field(key, "int")),
		Int16:     lib.NewNumeric[int16, T](lib.Field(key, "int_16")),
		Int16Ptr:  lib.NewNumericPtr[*int16, T](lib.Field(key, "int_16_ptr")),
		Int32:     lib.NewNumeric[int32, T](lib.Field(key, "int_32")),
		Int32Ptr:  lib.NewNumericPtr[*int32, T](lib.Field(key, "int_32_ptr")),
		Int64:     lib.NewNumeric[int64, T](lib.Field(key, "int_64")),
		Int64Ptr:  lib.NewNumericPtr[*int64, T](lib.Field(key, "int_64_ptr")),
		Int8:      lib.NewNumeric[int8, T](lib.Field(key, "int_8")),
		Int8Ptr:   lib.NewNumericPtr[*int8, T](lib.Field(key, "int_8_ptr")),
		IntPtr:    lib.NewNumericPtr[*int, T](lib.Field(key, "int_ptr")),
		Role:      lib.NewBase[model.Role, T](lib.Field(key, "role")),
		Rune:      lib.NewNumeric[rune, T](lib.Field(key, "rune")),
		String:    lib.NewString[T](lib.Field(key, "string")),
		StringPtr: lib.NewStringPtr[T](lib.Field(key, "string_ptr")),
		Time:      lib.NewTime[T](lib.Field(key, "time")),
		TimePtr:   lib.NewTimePtr[T](lib.Field(key, "time_ptr")),
		UUID:      lib.NewBase[uuid.UUID, T](lib.Field(key, "uuid")),
		UUIDPtr:   lib.NewBasePtr[uuid.UUID, T](lib.Field(key, "uuid_ptr")),
		Uint16:    lib.NewNumeric[uint16, T](lib.Field(key, "uint_16")),
		Uint16Ptr: lib.NewNumericPtr[*uint16, T](lib.Field(key, "uint_16_ptr")),
		Uint32:    lib.NewNumeric[uint32, T](lib.Field(key, "uint_32")),
		Uint32Ptr: lib.NewNumericPtr[*uint32, T](lib.Field(key, "uint_32_ptr")),
		Uint8:     lib.NewNumeric[uint8, T](lib.Field(key, "uint_8")),
		Uint8Ptr:  lib.NewNumericPtr[*uint8, T](lib.Field(key, "uint_8_ptr")),
		UpdatedAt: lib.NewTime[T](lib.Field(key, "updated_at")),
		key:       key,
	}
}

type allFieldTypes[T any] struct {
	key       lib.Key[T]
	ID        *lib.ID[T]
	CreatedAt *lib.Time[T]
	UpdatedAt *lib.Time[T]
	String    *lib.String[T]
	StringPtr *lib.StringPtr[T]
	Int       *lib.Numeric[int, T]
	IntPtr    *lib.NumericPtr[*int, T]
	Int8      *lib.Numeric[int8, T]
	Int8Ptr   *lib.NumericPtr[*int8, T]
	Int16     *lib.Numeric[int16, T]
	Int16Ptr  *lib.NumericPtr[*int16, T]
	Int32     *lib.Numeric[int32, T]
	Int32Ptr  *lib.NumericPtr[*int32, T]
	Int64     *lib.Numeric[int64, T]
	Int64Ptr  *lib.NumericPtr[*int64, T]
	Uint8     *lib.Numeric[uint8, T]
	Uint8Ptr  *lib.NumericPtr[*uint8, T]
	Uint16    *lib.Numeric[uint16, T]
	Uint16Ptr *lib.NumericPtr[*uint16, T]
	Uint32    *lib.Numeric[uint32, T]
	Uint32Ptr *lib.NumericPtr[*uint32, T]
	Float32   *lib.Numeric[float32, T]
	Float64   *lib.Numeric[float64, T]
	Rune      *lib.Numeric[rune, T]
	Bool      *lib.Bool[T]
	BoolPtr   *lib.BoolPtr[T]
	Time      *lib.Time[T]
	TimePtr   *lib.TimePtr[T]
	UUID      *lib.Base[uuid.UUID, T]
	UUIDPtr   *lib.BasePtr[uuid.UUID, T]
	Role      *lib.Base[model.Role, T]
	EnumPtr   *lib.BasePtr[model.Role, T]
}

func (n allFieldTypes[T]) Other() *lib.Slice[T, string] {
	return lib.NewSlice[T, string](lib.Field(n.key, "other"))
}

func (n allFieldTypes[T]) StringPtrSlice() *lib.Slice[T, *string] {
	return lib.NewSlice[T, *string](lib.Field(n.key, "string_ptr_slice"))
}

func (n allFieldTypes[T]) StringSlicePtr() *lib.Slice[T, string] {
	return lib.NewSlice[T, string](lib.Field(n.key, "string_slice_ptr"))
}

func (n allFieldTypes[T]) More() *lib.Slice[T, float32] {
	return lib.NewSlice[T, float32](lib.Field(n.key, "more"))
}

func (n allFieldTypes[T]) Roles() *lib.Slice[T, model.Role] {
	return lib.NewSlice[T, model.Role](lib.Field(n.key, "roles"))
}

func (n allFieldTypes[T]) EnumPtrSlice() *lib.Slice[T, model.Role] {
	return lib.NewSlice[T, model.Role](lib.Field(n.key, "enum_ptr_slice"))
}

func (n allFieldTypes[T]) EnumPtrSlicePtr() *lib.Slice[T, model.Role] {
	return lib.NewSlice[T, model.Role](lib.Field(n.key, "enum_ptr_slice_ptr"))
}

func (n allFieldTypes[T]) Login() login[T] {
	return newLogin[T](lib.Field(n.key, "login"))
}

func (n allFieldTypes[T]) StructPtr() someStruct[T] {
	return newSomeStruct[T](lib.Field(n.key, "struct_ptr"))
}

func (n allFieldTypes[T]) StructSlice() *lib.Slice[T, model.SomeStruct] {
	return lib.NewSlice[T, model.SomeStruct](lib.Field(n.key, "struct_slice"))
}

func (n allFieldTypes[T]) StructPtrSlice() *lib.Slice[T, *model.SomeStruct] {
	return lib.NewSlice[T, *model.SomeStruct](lib.Field(n.key, "struct_ptr_slice"))
}

func (n allFieldTypes[T]) StructPtrSlicePtr() *lib.Slice[T, *model.SomeStruct] {
	return lib.NewSlice[T, *model.SomeStruct](lib.Field(n.key, "struct_ptr_slice_ptr"))
}

func (n allFieldTypes[T]) MainGroup() group[T] {
	return newGroup[T](lib.Field(n.key, "main_group"))
}

func (n allFieldTypes[T]) MainGroupPtr() group[T] {
	return newGroup[T](lib.Field(n.key, "main_group_ptr"))
}

func (n allFieldTypes[T]) Groups(filters ...lib.Filter[model.Group]) groupSlice[T] {
	key := lib.Node(n.key, "groups", filters)
	return groupSlice[T]{lib.KeyFilter[T](key), lib.NewSlice[T, model.Group](key)}
}

func (n allFieldTypes[T]) NodePtrSlice(filters ...lib.Filter[model.Group]) groupSlice[T] {
	key := lib.Node(n.key, "node_ptr_slice", filters)
	return groupSlice[T]{lib.KeyFilter[T](key), lib.NewSlice[T, model.Group](key)}
}

func (n allFieldTypes[T]) NodePtrSlicePtr(filters ...lib.Filter[model.Group]) groupSlice[T] {
	key := lib.Node(n.key, "node_ptr_slice_ptr", filters)
	return groupSlice[T]{lib.KeyFilter[T](key), lib.NewSlice[T, model.Group](key)}
}

func (n allFieldTypes[T]) MemberOf(filters ...lib.Filter[model.GroupMember]) groupMemberIn[T] {
	return newGroupMemberIn[T](lib.EdgeIn(n.key, "group_member", filters))
}

func (n allFieldTypes[T]) SliceSlice() *lib.Slice[T, []string] {
	return lib.NewSlice[T, []string](lib.Field(n.key, "slice_slice"))
}

type allFieldTypesEdges[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

func (n allFieldTypesEdges[T]) MemberOf(filters ...lib.Filter[model.GroupMember]) groupMemberIn[T] {
	return newGroupMemberIn[T](lib.EdgeIn(n.key, "group_member", filters))
}

type allFieldTypesSlice[T any] struct {
	lib.Filter[T]
	*lib.Slice[T, model.AllFieldTypes]
}
