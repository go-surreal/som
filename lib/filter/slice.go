package filter

import (
	"github.com/marcbinz/som/lib/builder"
)

type Slice[T any, R any] struct {
	key Key
}

func NewSlice[T, R any](key Key) *Slice[T, R] {
	return &Slice[T, R]{key: key}
}

func (s *Slice[T, R]) Contains(val T) Of[R] {
	return build[R](s.key, builder.OpContains, val, false)
}

func (s *Slice[T, R]) ContainsNot(val T) Of[R] {
	return build[R](s.key, builder.OpContainsNot, val, false)
}

func (s *Slice[T, R]) ContainsAll(vals []T) Of[R] {
	return build[R](s.key, builder.OpContainsAll, vals, false)
}

func (s *Slice[T, R]) ContainsAny(vals []T) Of[R] {
	return build[R](s.key, builder.OpContainsAny, vals, false)
}

func (s *Slice[T, R]) ContainsNone(vals []T) Of[R] {
	return build[R](s.key, builder.OpContainsNone, vals, false)
}

func (s *Slice[T, R]) Count() *Numeric[int, R] {
	return newNumeric[int, R](s.key, true)
}

type SlicePtr[T, R any] struct {
	*Slice[T, R]
	*Nillable[R]
}

func NewSlicePtr[T, R any](key Key) *SlicePtr[T, R] {
	return &SlicePtr[T, R]{
		Slice:    &Slice[T, R]{key: key},
		Nillable: &Nillable[R]{key: key},
	}
}
