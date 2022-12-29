// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package by

import (
	model "github.com/marcbinz/som/example/model"
	sort "github.com/marcbinz/som/lib/sort"
)

var Group = newGroup[model.Group]("")

func newGroup[T any](key string) group[T] {
	return group[T]{
		CreatedAt: sort.NewSort[T](keyed(key, "created_at")),
		ID:        sort.NewSort[T](keyed(key, "id")),
		Name:      sort.NewString[T](keyed(key, "name")),
		UpdatedAt: sort.NewSort[T](keyed(key, "updated_at")),
		key:       key,
	}
}

type group[T any] struct {
	key       string
	ID        *sort.Sort[T]
	CreatedAt *sort.Sort[T]
	UpdatedAt *sort.Sort[T]
	Name      *sort.String[T]
}
