package where

import (
	model "github.com/marcbinz/som/example/model"
	filter "github.com/marcbinz/som/lib/filter"
)

var Group = newGroup[model.Group](filter.NewKey())

func newGroup[T any](key filter.Key) group[T] {
	return group[T]{
		ID:   filter.NewID[T](key.Dot("id"), "group"),
		Name: filter.NewString[T](key.Dot("name")),
		key:  key,
	}
}

type group[T any] struct {
	key  filter.Key
	ID   *filter.ID[T]
	Name *filter.String[T]
}
type groupSlice[T any] struct {
	group[T]
	*filter.Slice[model.Group, T]
}

func (n group[T]) Members() memberOfOut[T] {
	return newMemberOfOut[T](n.key.Out("member_of"))
}