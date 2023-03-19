// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package where

import (
	model "github.com/marcbinz/som/examples/movie/model"
	lib "github.com/marcbinz/som/lib"
)

var Movie = newMovie[model.Movie](lib.NewKey[model.Movie]())

func newMovie[T any](key lib.Key[T]) movie[T] {
	return movie[T]{
		ID:    lib.NewID[T](lib.Field(key, "id"), "movie"),
		Title: lib.NewString[T](lib.Field(key, "title")),
		key:   key,
	}
}

type movie[T any] struct {
	key   lib.Key[T]
	ID    *lib.ID[T]
	Title *lib.String[T]
}

type movieEdges[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

type movieSlice[T any] struct {
	lib.Filter[T]
	*lib.Slice[T, model.Movie]
}
