package filter

import (
	"github.com/marcbinz/som/lib/builder"
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

type BoolPtr[R any] struct {
	*Bool[R]
	*Nillable[R]
}

func NewBoolPtr[R any](key Key) *BoolPtr[R] {
	return &BoolPtr[R]{
		Bool:     &Bool[R]{key: key},
		Nillable: &Nillable[R]{key: key},
	}
}
