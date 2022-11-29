package filter

import (
	"github.com/marcbinz/som/lib/builder"
)

type Of[T any] struct {
	builder.Where
}

func build[R any](key Key, op builder.Operator, val any, isCount bool) Of[R] {
	return Of[R]{&builder.Predicate{Op: op, Key: key.key, Val: val, Close: key.close, IsCount: isCount}}
}

func All[R any](filters []Of[R]) Of[R] {
	whereAll := &builder.WhereAll{}
	for _, f := range filters {
		whereAll.Where = append(whereAll.Where, builder.Where(f))
	}
	return Of[R]{whereAll}
}

func Any[R any](filters []Of[R]) Of[R] {
	whereAny := &builder.WhereAny{}
	for _, f := range filters {
		whereAny.Where = append(whereAny.Where, builder.Where(f))
	}
	return Of[R]{whereAny}
}
