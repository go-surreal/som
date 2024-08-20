//go:build embed

package lib

import (
	"golang.org/x/exp/constraints"
)

type Float[M any, T constraints.Float] struct {
	*Numeric[M, T]
}

func NewFloat[M any, T constraints.Float](key Key[M]) *Float[M, T] {
	return &Float[M, T]{
		Numeric: NewNumeric[M, T](key),
	}
}

func (f *Float[M, T]) Ceil() *Int[M, int] {
	return NewInt[M, int](f.Base.fn("math::ceil"))
}

func (f *Float[M, T]) Fixed(places int) *Float[M, T] {
	return NewFloat[M, T](f.Base.fn("math::fixed", places))
}

func (f *Float[M, T]) Floor() *Int[M, int] {
	return NewInt[M, int](f.Base.fn("math::floor"))
}

func (f *Float[M, T]) Round() *Int[M, int] {
	return NewInt[M, int](f.Base.fn("math::round"))
}
