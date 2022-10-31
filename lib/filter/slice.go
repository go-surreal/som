package filter

import (
	"github.com/marcbinz/sdb/lib/builder"
)

type Slice[T any, R any] struct {
	key string
}

func NewSlice[T, R any](key string) *Slice[T, R] {
	return &Slice[T, R]{key: key}
}

func (s *Slice[T, R]) Contains(val T) Of[R] {
	return newOf[R](builder.OpContains, s.key, val, false)
}

func (s *Slice[T, R]) ContainsNot(val T) Of[R] {
	return newOf[R](builder.OpContainsNot, s.key, val, false)
}

func (s *Slice[T, R]) ContainsAll(vals []T) Of[R] {
	return newOf[R](builder.OpContainsAll, s.key, vals, false)
}

func (s *Slice[T, R]) ContainsAny(vals []T) Of[R] {
	return newOf[R](builder.OpContainsAny, s.key, vals, false)
}

func (s *Slice[T, R]) ContainsNone(vals []T) Of[R] {
	return newOf[R](builder.OpContainsNone, s.key, vals, false)
}

func (s *Slice[T, R]) Count() *Numeric[int, R] {
	return newCountNumeric[int, R](s.key)
}
