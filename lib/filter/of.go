package filter

import (
	"github.com/marcbinz/sdb/lib/builder"
)

type Of[T any] builder.Where

func newOf[T any](op builder.Operator, key string, val any, isCount bool) Of[T] {
	return &builder.Predicate{Op: op, Key: key, Val: val, IsCount: isCount}
}

func All[R any](filters []Of[R]) Of[R] {
	whereAll := &builder.WhereAll{}
	for _, f := range filters {
		whereAll.Where = append(whereAll.Where, builder.Where(f))
	}
	return whereAll
}

func Any[R any](filters []Of[R]) Of[R] {
	whereAny := &builder.WhereAny{}
	for _, f := range filters {
		whereAny.Where = append(whereAny.Where, builder.Where(f))
	}
	return whereAny
}
