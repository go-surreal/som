package filter

import (
	"github.com/marcbinz/sdb/lib/builder"
)

type Base[T any, R any] struct {
	key     string
	isCount bool
}

func NewBase[T, R any](key string) *Base[T, R] {
	return &Base[T, R]{key: key}
}

func (b *Base[T, R]) Equal(val T) Of[R] {
	return newOf[R](builder.OpEqual, b.key, val, b.isCount)
}

func (b *Base[T, R]) NotEqual(val T) Of[R] {
	return newOf[R](builder.OpEqual, b.key, val, b.isCount)
}

func (b *Base[T, R]) In(vals []T) Of[R] {
	return newOf[R](builder.OpInside, b.key, vals, b.isCount)
}

func (b *Base[T, R]) NotIn(vals []T) Of[R] {
	return newOf[R](builder.OpNotInside, b.key, vals, b.isCount)
}
