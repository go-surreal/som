// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package where

import (
	lib "github.com/marcbinz/som/examples/basic/gen/som/internal/lib"
	model "github.com/marcbinz/som/examples/basic/model"
)

var GroupMember = newGroupMember[model.GroupMember](lib.NewKey[model.GroupMember]())

func newGroupMember[T any](key lib.Key[T]) groupMember[T] {
	return groupMember[T]{
		CreatedAt: lib.NewTime[T](lib.Field(key, "created_at")),
		ID:        lib.NewID[T](lib.Field(key, "id"), "group_member"),
		UpdatedAt: lib.NewTime[T](lib.Field(key, "updated_at")),
		key:       key,
	}
}

type groupMember[T any] struct {
	key       lib.Key[T]
	ID        *lib.ID[T]
	CreatedAt *lib.Time[T]
	UpdatedAt *lib.Time[T]
}

func (n groupMember[T]) Meta() groupMemberMeta[T] {
	return newGroupMemberMeta[T](lib.Field(n.key, "meta"))
}

type groupMemberIn[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

func newGroupMemberIn[T any](key lib.Key[T]) groupMemberIn[T] {
	return groupMemberIn[T]{lib.KeyFilter(key), key}
}

func (i groupMemberIn[T]) Group(filters ...lib.Filter[model.Group]) groupEdges[T] {
	key := lib.EdgeIn(i.key, "group", filters)
	return groupEdges[T]{lib.KeyFilter(key), key}
}

type groupMemberOut[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

func newGroupMemberOut[T any](key lib.Key[T]) groupMemberOut[T] {
	return groupMemberOut[T]{lib.KeyFilter(key), key}
}

func (o groupMemberOut[T]) User(filters ...lib.Filter[model.User]) userEdges[T] {
	key := lib.EdgeOut(o.key, "user", filters)
	return userEdges[T]{lib.KeyFilter(key), key}
}