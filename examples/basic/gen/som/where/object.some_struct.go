// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import (
	lib "github.com/go-surreal/som/examples/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/examples/basic/model"
	uuid "github.com/google/uuid"
)

func newSomeStruct[T any](key lib.Key[T]) someStruct[T] {
	return someStruct[T]{
		IntPtr:    lib.NewNumericPtr[*int, T](lib.Field(key, "int_ptr")),
		StringPtr: lib.NewStringPtr[T](lib.Field(key, "string_ptr")),
		TimePtr:   lib.NewTimePtr[T](lib.Field(key, "time_ptr")),
		UuidPtr:   lib.NewBasePtr[uuid.UUID, T](lib.Field(key, "uuid_ptr")),
		key:       key,
	}
}

type someStruct[T any] struct {
	key       lib.Key[T]
	StringPtr *lib.StringPtr[T]
	IntPtr    *lib.NumericPtr[*int, T]
	TimePtr   *lib.TimePtr[T]
	UuidPtr   *lib.BasePtr[uuid.UUID, T]
}

type someStructEdges[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

type someStructSlice[T any] struct {
	lib.Filter[T]
	*lib.Slice[T, model.SomeStruct]
}
