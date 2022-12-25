package filter

import (
	"github.com/marcbinz/som/lib/builder"
)

type Base[T any, R any] struct {
	key     Key
	isCount bool
}

func NewBase[T, R any](key Key) *Base[T, R] {
	return &Base[T, R]{key: key}
}

func (b *Base[T, R]) Equal(val T) Of[R] {
	return build[R](b.key, builder.OpEqual, val, b.isCount)
}

func (b *Base[T, R]) NotEqual(val T) Of[R] {
	return build[R](b.key, builder.OpNotEqual, val, b.isCount)
}

func (b *Base[T, R]) In(vals []T) Of[R] {
	return build[R](b.key, builder.OpInside, vals, b.isCount)
}

func (b *Base[T, R]) NotIn(vals []T) Of[R] {
	return build[R](b.key, builder.OpNotInside, vals, b.isCount)
}

type BasePtr[T, R any] struct {
	*Base[T, R]
	*Nillable[R]
}

func NewBasePtr[T, R any](key Key) *BasePtr[T, R] {
	return &BasePtr[T, R]{
		Base:     &Base[T, R]{key: key},
		Nillable: &Nillable[R]{key: key},
	}
}
