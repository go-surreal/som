// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package by

import (
	model "github.com/marcbinz/som/examples/movie/model"
	lib "github.com/marcbinz/som/lib"
)

var Movie = newMovie[model.Movie]("")

func newMovie[T any](key string) movie[T] {
	return movie[T]{
		ID:    lib.NewBaseSort[T](keyed(key, "id")),
		Title: lib.NewStringSort[T](keyed(key, "title")),
		key:   key,
	}
}

type movie[T any] struct {
	key   string
	ID    *lib.BaseSort[T]
	Title *lib.StringSort[T]
}
