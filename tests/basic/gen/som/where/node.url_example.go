// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import (
	lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/tests/basic/model"
)

var URLExample = newURLExample[model.URLExample](lib.NewKey[model.URLExample]())

func newURLExample[T any](key lib.Key[T]) urlexample[T] {
	return urlexample[T]{
		ID:           lib.NewID[T](lib.Field(key, "id"), "url_example"),
		SomeOtherURL: lib.NewURL[T](lib.Field(key, "some_other_url")),
		SomeURL:      lib.NewURLPtr[T](lib.Field(key, "some_url")),
		key:          key,
	}
}

type urlexample[T any] struct {
	key          lib.Key[T]
	ID           *lib.ID[T]
	SomeURL      *lib.URLPtr[T]
	SomeOtherURL *lib.URL[T]
}

type urlexampleEdges[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

type urlexampleSlice[T any] struct {
	lib.Filter[T]
	*lib.Slice[T, model.URLExample, urlexample[T]]
}
