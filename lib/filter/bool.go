package filter

import (
	"github.com/marcbinz/sdb/lib/builder"
)

type Bool[R any] struct {
	key string
}

func NewBool[R any](key string) *Bool[R] {
	return &Bool[R]{key: key}
}

func (b *Bool[R]) Is(val bool) Of[R] {
	return newOf[R](builder.OpExactlyEqual, b.key, val, false)
}
