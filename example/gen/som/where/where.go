// Code generated by github.com/marcbinz/som, DO NOT EDIT.

package where

import "github.com/marcbinz/som/lib"

func All[T any](filters ...lib.Filter[T]) lib.Filter[T] {
	return lib.Filter[T](lib.All(lib.ToWhere(filters)))
}

func Any[T any](filters ...lib.Filter[T]) lib.Filter[T] {
	return lib.Filter[T](lib.Any(lib.ToWhere(filters)))
}
