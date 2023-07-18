// Code generated by github.com/marcbinz/som, DO NOT EDIT.
package where

import (
	lib "github.com/marcbinz/som/examples/basic/gen/som/internal/lib"
	model "github.com/marcbinz/som/examples/basic/model"
)

func newGroupMemberMeta[T any](key lib.Key[T]) groupMemberMeta[T] {
	return groupMemberMeta[T]{
		IsActive: lib.NewBool[T](lib.Field(key, "is_active")),
		IsAdmin:  lib.NewBool[T](lib.Field(key, "is_admin")),
		key:      key,
	}
}

type groupMemberMeta[T any] struct {
	key      lib.Key[T]
	IsAdmin  *lib.Bool[T]
	IsActive *lib.Bool[T]
}

type groupMemberMetaEdges[T any] struct {
	lib.Filter[T]
	key lib.Key[T]
}

type groupMemberMetaSlice[T any] struct {
	lib.Filter[T]
	*lib.Slice[T, model.GroupMemberMeta]
}
