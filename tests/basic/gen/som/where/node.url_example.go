// Code generated by github.com/go-surreal/som, DO NOT EDIT.
package where

import (
	lib "github.com/go-surreal/som/tests/basic/gen/som/internal/lib"
	model "github.com/go-surreal/som/tests/basic/model"
)

var URLExample = newURLExample[model.URLExample](lib.NewKey[model.URLExample]())

func newURLExample[M any](key lib.Key[M]) urlexample[M] {
	return urlexample[M]{
		ID:           lib.NewID[M](lib.Field(key, "id"), "url_example"),
		Key:          key,
		SomeOtherURL: lib.NewURL[M](lib.Field(key, "some_other_url")),
		SomeURL:      lib.NewURLPtr[M](lib.Field(key, "some_url")),
	}
}

type urlexample[M any] struct {
	lib.Key[M]
	ID           *lib.ID[M]
	SomeURL      *lib.URLPtr[M]
	SomeOtherURL *lib.URL[M]
}

type urlexampleEdges[M any] struct {
	lib.Filter[M]
	lib.Key[M]
}
