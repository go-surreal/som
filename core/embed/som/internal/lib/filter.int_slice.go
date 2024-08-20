//go:build embed

package lib

import (
	"golang.org/x/exp/constraints"
)

type IntSlice[M any, T constraints.Integer] struct {
	*NumericSlice[M, T, *Int[M, T]]
}

func NewIntSlice[M any, T constraints.Integer](key Key[M]) *IntSlice[M, T] {
	return &IntSlice[M, T]{
		NumericSlice: NewNumericSlice[M, T](key, NewInt[M, T]),
	}
}

func (s *IntSlice[M, T]) Bottom(count int) *IntSlice[M, T] {
	return NewIntSlice[M, T](s.fn("math::bottom", count))
}

func (s *IntSlice[M, T]) Max() *Int[M, T] {
	return NewInt[M, T](s.fn("math::abs"))
}

func (s *IntSlice[M, T]) Min() *Int[M, T] {
	return NewInt[M, T](s.fn("math::min"))
}

func (s *IntSlice[M, T]) Mode() *Int[M, T] {
	return NewInt[M, T](s.fn("math::mode"))
}

func (s *IntSlice[M, T]) Product() *Int[M, T] {
	return NewInt[M, T](s.fn("math::product"))
}

func (s *IntSlice[M, T]) Spread() *Int[M, T] {
	return NewInt[M, T](s.fn("math::spread"))
}

func (s *IntSlice[M, T]) Sum() *Int[M, T] {
	return NewInt[M, T](s.fn("math::sum"))
}

func (s *IntSlice[M, T]) Top(count int) *IntSlice[M, T] {
	return NewIntSlice[M, T](s.fn("math::top", count))
}
