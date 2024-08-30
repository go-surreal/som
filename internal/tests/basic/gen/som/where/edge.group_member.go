// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import (
	lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/tests/basic/model"
)

var GroupMember = newGroupMember[model.GroupMember](lib.NewKey[model.GroupMember]())

func newGroupMember[M any](key lib.Key[M]) groupMember[M] {
	return groupMember[M]{
		CreatedAt: lib.NewTime[M](lib.Field(key, "created_at")),
		ID:        lib.NewID[M](lib.Field(key, "id"), "group_member"),
		Key:       key,
		UpdatedAt: lib.NewTime[M](lib.Field(key, "updated_at")),
	}
}

type groupMember[M any] struct {
	lib.Key[M]
	ID        *lib.ID[M]
	CreatedAt *lib.Time[M]
	UpdatedAt *lib.Time[M]
}

func (n groupMember[M]) Meta() groupMemberMeta[M] {
	return newGroupMemberMeta[M](lib.Field(n.Key, "meta"))
}

type groupMemberIn[M any] struct {
	lib.Filter[M]
	key lib.Key[M]
}

func newGroupMemberIn[M any](key lib.Key[M]) groupMemberIn[M] {
	return groupMemberIn[M]{lib.KeyFilter(key), key}
}

func (i groupMemberIn[M]) Group(filters ...lib.Filter[model.Group]) groupEdges[M] {
	key := lib.EdgeIn(i.key, "group", filters)
	return groupEdges[M]{lib.KeyFilter(key), key}
}

type groupMemberOut[M any] struct {
	lib.Filter[M]
	key lib.Key[M]
}

func newGroupMemberOut[M any](key lib.Key[M]) groupMemberOut[M] {
	return groupMemberOut[M]{lib.KeyFilter(key), key}
}

func (o groupMemberOut[M]) User(filters ...lib.Filter[model.AllFieldTypes]) allFieldTypesEdges[M] {
	key := lib.EdgeOut(o.key, "user", filters)
	return allFieldTypesEdges[M]{lib.KeyFilter(key), key}
}
