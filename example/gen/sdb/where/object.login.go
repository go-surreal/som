package where

import (
	model "github.com/marcbinz/sdb/example/model"
	filter "github.com/marcbinz/sdb/lib/filter"
)

func newLogin[T any](key filter.Key) login[T] {
	return login[T]{
		Password: filter.NewString[T](key.Dot("password")),
		Username: filter.NewString[T](key.Dot("username")),
		key:      key,
	}
}

type login[T any] struct {
	key      filter.Key
	Username *filter.String[T]
	Password *filter.String[T]
}
type loginSlice[T any] struct {
	login[T]
	*filter.Slice[model.Login, T]
}
