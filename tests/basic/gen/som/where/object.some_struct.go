// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"

func newSomeStruct[M any](key lib.Key[M]) someStruct[M] {
	return someStruct[M]{
		IntPtr:    lib.NewNumericPtr[M, *int](lib.Field(key, "int_ptr")),
		Key:       key,
		StringPtr: lib.NewStringPtr[M](lib.Field(key, "string_ptr")),
		TimePtr:   lib.NewTimePtr[M](lib.Field(key, "time_ptr")),
		UuidPtr:   lib.NewUUIDPtr[M](lib.Field(key, "uuid_ptr")),
	}
}

type someStruct[M any] struct {
	lib.Key[M]
	StringPtr *lib.StringPtr[M]
	IntPtr    *lib.NumericPtr[M, *int]
	TimePtr   *lib.TimePtr[M]
	UuidPtr   *lib.UUIDPtr[M]
}

type someStructEdges[M any] struct {
	lib.Filter[M]
	lib.Key[M]
}
