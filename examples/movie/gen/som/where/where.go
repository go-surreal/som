// Code generated by github.com/marcbinz/som, DO NOT EDIT.

package where

import "github.com/marcbinz/som/examples/movie/gen/som/internal/lib"

func All[T any](filters ...lib.Filter[T]) lib.Filter[T] {
	return lib.All[T](filters)
}

func Any[T any](filters ...lib.Filter[T]) lib.Filter[T] {
	return lib.Any[T](filters)
}
