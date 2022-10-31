package where

import filter "github.com/marcbinz/sdb/lib/filter"

func All[T any](filters ...filter.Of[T]) filter.Of[T] {
	return filter.All[T](filters)
}
func Any[T any](filters ...filter.Of[T]) filter.Of[T] {
	return filter.Any[T](filters)
}
func keyed(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}
