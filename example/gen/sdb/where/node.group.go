package where

import (
	model "github.com/marcbinz/sdb/example/model"
	filter "github.com/marcbinz/sdb/lib/filter"
)

var Group = newGroup[model.Group]("")

func newGroup[T any](key string) group[T] {
	return group[T]{
		ID:   filter.NewID[T](keyed(key, "id"), "group"),
		Name: filter.NewString[T](keyed(key, "name")),
		key:  key,
	}
}

type group[T any] struct {
	key  string
	ID   *filter.ID[T]
	Name *filter.String[T]
}
type groupSlice[T any] struct {
	group[T]
	*filter.Slice[model.Group, T]
}
