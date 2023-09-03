// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import (
	lib "github.com/go-surreal/som/examples/testing/gen/som/internal/lib"
	model "github.com/go-surreal/som/examples/testing/model"
)

var FieldsLikeDBResponse = newFieldsLikeDBResponse[model.FieldsLikeDBResponse](lib.NewKey[model.FieldsLikeDBResponse]())

func newFieldsLikeDBResponse[T any](key lib.Key[T]) fieldsLikeDBResponse[T] {
	return fieldsLikeDBResponse[T]{
		Detail: lib.NewString[T](lib.Field(key, "detail")),
		ID:     lib.NewID[T](lib.Field(key, "id"), "fields_like_db_response"),
		Status: lib.NewString[T](lib.Field(key, "status")),
		Time:   lib.NewString[T](lib.Field(key, "time")),
		key:    key,
	}
}

type fieldsLikeDBResponse[T any] struct {
	key    lib.Key[T]
	ID     *lib.ID[T]
	Time   *lib.String[T]
	Status *lib.String[T]
	Detail *lib.String[T]
}

func (n fieldsLikeDBResponse[T]) Result() *lib.Slice[T, string] {
	return lib.NewSlice[T, string](lib.Field(n.key, "result"))
}

type fieldsLikeDBResponseEdges[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

type fieldsLikeDBResponseSlice[T any] struct {
	lib.Filter[T]
	*lib.Slice[T, model.FieldsLikeDBResponse]
}
