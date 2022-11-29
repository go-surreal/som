package by

import (
	model "github.com/marcbinz/som/example/model"
	sort "github.com/marcbinz/som/lib/sort"
)

var Group = newGroup[model.Group]("")

func newGroup[T any](key string) group[T] {
	return group[T]{
		ID:   sort.NewSort[T](keyed(key, "id")),
		Name: sort.NewString[T](keyed(key, "name")),
		key:  key,
	}
}

type group[T any] struct {
	key  string
	ID   *sort.Sort[T]
	Name *sort.String[T]
}
