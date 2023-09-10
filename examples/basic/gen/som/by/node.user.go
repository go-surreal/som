// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package by

import (
	lib "github.com/go-surreal/som/examples/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/examples/basic/model"
)

var User = newUser[model.User]("")

func newUser[T any](key string) user[T] {
	return user[T]{
		CreatedAt:  lib.NewBaseSort[T](keyed(key, "created_at")),
		Float32:    lib.NewBaseSort[T](keyed(key, "float_32")),
		Float64:    lib.NewBaseSort[T](keyed(key, "float_64")),
		ID:         lib.NewBaseSort[T](keyed(key, "id")),
		Int:        lib.NewBaseSort[T](keyed(key, "int")),
		Int16:      lib.NewBaseSort[T](keyed(key, "int_16")),
		Int16Ptr:   lib.NewBaseSort[T](keyed(key, "int_16_ptr")),
		Int32:      lib.NewBaseSort[T](keyed(key, "int_32")),
		Int32Ptr:   lib.NewBaseSort[T](keyed(key, "int_32_ptr")),
		Int64:      lib.NewBaseSort[T](keyed(key, "int_64")),
		Int64Ptr:   lib.NewBaseSort[T](keyed(key, "int_64_ptr")),
		Int8:       lib.NewBaseSort[T](keyed(key, "int_8")),
		Int8Ptr:    lib.NewBaseSort[T](keyed(key, "int_8_ptr")),
		IntPtr:     lib.NewBaseSort[T](keyed(key, "int_ptr")),
		Rune:       lib.NewBaseSort[T](keyed(key, "rune")),
		String:     lib.NewStringSort[T](keyed(key, "string")),
		StringPtr:  lib.NewStringSort[T](keyed(key, "string_ptr")),
		Time:       lib.NewBaseSort[T](keyed(key, "time")),
		TimePtr:    lib.NewBaseSort[T](keyed(key, "time_ptr")),
		Uint:       lib.NewBaseSort[T](keyed(key, "uint")),
		Uint16:     lib.NewBaseSort[T](keyed(key, "uint_16")),
		Uint16Ptr:  lib.NewBaseSort[T](keyed(key, "uint_16_ptr")),
		Uint32:     lib.NewBaseSort[T](keyed(key, "uint_32")),
		Uint32Ptr:  lib.NewBaseSort[T](keyed(key, "uint_32_ptr")),
		Uint64:     lib.NewBaseSort[T](keyed(key, "uint_64")),
		Uint64Ptr:  lib.NewBaseSort[T](keyed(key, "uint_64_ptr")),
		Uint8:      lib.NewBaseSort[T](keyed(key, "uint_8")),
		Uint8Ptr:   lib.NewBaseSort[T](keyed(key, "uint_8_ptr")),
		UintPtr:    lib.NewBaseSort[T](keyed(key, "uint_ptr")),
		Uintptr:    lib.NewBaseSort[T](keyed(key, "uintptr")),
		UintptrPtr: lib.NewBaseSort[T](keyed(key, "uintptr_ptr")),
		UpdatedAt:  lib.NewBaseSort[T](keyed(key, "updated_at")),
		key:        key,
	}
}

type user[T any] struct {
	key        string
	ID         *lib.BaseSort[T]
	CreatedAt  *lib.BaseSort[T]
	UpdatedAt  *lib.BaseSort[T]
	String     *lib.StringSort[T]
	StringPtr  *lib.StringSort[T]
	Int        *lib.BaseSort[T]
	IntPtr     *lib.BaseSort[T]
	Int8       *lib.BaseSort[T]
	Int8Ptr    *lib.BaseSort[T]
	Int16      *lib.BaseSort[T]
	Int16Ptr   *lib.BaseSort[T]
	Int32      *lib.BaseSort[T]
	Int32Ptr   *lib.BaseSort[T]
	Int64      *lib.BaseSort[T]
	Int64Ptr   *lib.BaseSort[T]
	Uint       *lib.BaseSort[T]
	UintPtr    *lib.BaseSort[T]
	Uint8      *lib.BaseSort[T]
	Uint8Ptr   *lib.BaseSort[T]
	Uint16     *lib.BaseSort[T]
	Uint16Ptr  *lib.BaseSort[T]
	Uint32     *lib.BaseSort[T]
	Uint32Ptr  *lib.BaseSort[T]
	Uint64     *lib.BaseSort[T]
	Uint64Ptr  *lib.BaseSort[T]
	Uintptr    *lib.BaseSort[T]
	UintptrPtr *lib.BaseSort[T]
	Float32    *lib.BaseSort[T]
	Float64    *lib.BaseSort[T]
	Rune       *lib.BaseSort[T]
	Time       *lib.BaseSort[T]
	TimePtr    *lib.BaseSort[T]
}

func (n user[T]) MainGroup() group[T] {
	return newGroup[T](keyed(n.key, "main_group"))
}

func (n user[T]) MainGroupPtr() group[T] {
	return newGroup[T](keyed(n.key, "main_group_ptr"))
}
