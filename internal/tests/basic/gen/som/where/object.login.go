// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"

func newLogin[M any](key lib.Key[M]) login[M] {
	return login[M]{
		Key:      key,
		Password: lib.NewString[M](lib.Field(key, "password")),
		Username: lib.NewString[M](lib.Field(key, "username")),
	}
}

type login[M any] struct {
	lib.Key[M]
	Username *lib.String[M]
	Password *lib.String[M]
}

type loginEdges[M any] struct {
	lib.Filter[M]
	lib.Key[M]
}
