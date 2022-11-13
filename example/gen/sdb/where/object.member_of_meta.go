package where

import (
	model "github.com/marcbinz/sdb/example/model"
	filter "github.com/marcbinz/sdb/lib/filter"
)

func newMemberOfMeta[T any](key filter.Key) memberOfMeta[T] {
	return memberOfMeta[T]{
		IsActive: filter.NewBool[T](key.Dot("is_active")),
		IsAdmin:  filter.NewBool[T](key.Dot("is_admin")),
		key:      key,
	}
}

type memberOfMeta[T any] struct {
	key      filter.Key
	IsAdmin  *filter.Bool[T]
	IsActive *filter.Bool[T]
}
type memberOfMetaSlice[T any] struct {
	memberOfMeta[T]
	*filter.Slice[model.MemberOfMeta, T]
}
