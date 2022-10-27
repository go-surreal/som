package where

import (
	model "github.com/marcbinz/sdb/example/model"
	filter "github.com/marcbinz/sdb/lib/filter"
)

func newLogin[T any](key string) login[T] {
	return login[T]{
		Password: filter.NewString[T](key),
		Username: filter.NewString[T](key),
	}
}

type login[T any] struct {
	Username *filter.String[T]
	Password *filter.String[T]
}
type loginSlice[T any] struct {
	login[T]
	*filter.Slice[model.Login, T]
}
