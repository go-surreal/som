package where

import filter "github.com/marcbinz/som/lib/filter"

func All[T any](filters ...filter.Of[T]) filter.Of[T] {
	return filter.All[T](filters)
}

func Any[T any](filters ...filter.Of[T]) filter.Of[T] {
	return filter.Any[T](filters)
}