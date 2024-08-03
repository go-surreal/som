// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import (
	lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/tests/basic/model"
)

var Group = newGroup[model.Group](lib.NewKey[model.Group]())

func newGroup[T any](key lib.Key[T]) group[T] {
	return group[T]{
		CreatedAt: lib.NewTime[T](lib.Field(key, "created_at")),
		ID:        lib.NewID[T](lib.Field(key, "id"), "group"),
		Name:      lib.NewString[T](lib.Field(key, "name")),
		UpdatedAt: lib.NewTime[T](lib.Field(key, "updated_at")),
		key:       key,
	}
}

type group[T any] struct {
	key       lib.Key[T]
	ID        *lib.ID[T]
	CreatedAt *lib.Time[T]
	UpdatedAt *lib.Time[T]
	Name      *lib.String[T]
}

func (n group[T]) Members(filters ...lib.Filter[model.GroupMember]) groupMemberOut[T] {
	return newGroupMemberOut[T](lib.EdgeOut(n.key, "group_member", filters))
}

type groupEdges[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

func (n groupEdges[T]) Members(filters ...lib.Filter[model.GroupMember]) groupMemberOut[T] {
	return newGroupMemberOut[T](lib.EdgeOut(n.key, "group_member", filters))
}

type groupSlice[T any] struct {
	lib.Filter[T]
	*lib.Slice[T, model.Group, group[T]]
}
