// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package by

import (
	lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/tests/basic/model"
)

var URLExample = newURLExample[model.URLExample]("")

func newURLExample[T any](key string) urlexample[T] {
	return urlexample[T]{
		ID:  lib.NewBaseSort[M](keyed(key, "id")),
		key: key,
	}
}

type urlexample[T any] struct {
	key string
	ID  *lib.BaseSort[M]
}
