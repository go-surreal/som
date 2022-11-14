package filter

import (
	"github.com/marcbinz/sdb/lib/builder"
)

type Bool[R any] struct {
	key Key
}

func NewBool[R any](key Key) *Bool[R] {
	return &Bool[R]{key: key}
}

func (b *Bool[R]) Is(val bool) Of[R] {
	return build[R](b.key, builder.OpExactlyEqual, val, false)
}
