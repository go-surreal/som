package filter

import (
	"golang.org/x/exp/constraints"
)

type Numeric[T Number, R any] struct {
	*Base[T, R]
	*Comparable[T, R]
}

func NewNumeric[T Number, R any](key string) *Numeric[T, R] {
	return &Numeric[T, R]{
		Base:       &Base[T, R]{key: key},
		Comparable: &Comparable[T, R]{key: key},
	}
}

type Number interface {
	constraints.Integer | constraints.Float
}
