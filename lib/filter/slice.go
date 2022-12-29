package filter

import (
	"github.com/marcbinz/som/lib/builder"
)

// Slice is a filter that can be used for slice fields.
// T is the type of the outgoing table for the filter statement.
// E is the type of the slice elements.
type Slice[T, E any] struct {
	key     Key
	filters []Of[E]
}

// NewSlice creates a new slice filter.
func NewSlice[T, E any](key Key, filters []Of[E]) *Slice[T, E] {
	return &Slice[T, E]{
		key:     key,
		filters: filters,
	}
}

func (s *Slice[T, E]) Contains(val T) Of[T] {
	return build[T](s.key, builder.OpContains, val, false)
}

func (s *Slice[T, E]) ContainsNot(val T) Of[T] {
	return build[T](s.key, builder.OpContainsNot, val, false)
}

func (s *Slice[T, E]) ContainsAll(vals []T) Of[T] {
	return build[T](s.key, builder.OpContainsAll, vals, false)
}

func (s *Slice[T, E]) ContainsAny(vals []T) Of[T] {
	return build[T](s.key, builder.OpContainsAny, vals, false)
}

func (s *Slice[T, E]) ContainsNone(vals []T) Of[T] {
	return build[T](s.key, builder.OpContainsNone, vals, false)
}

func (s *Slice[T, E]) Count() *Numeric[int, T] {
	return newNumeric[int, T](s.key, true)
}

type SlicePtr[T, E any] struct {
	*Slice[T, E]
	*Nillable[T]
}

func NewSlicePtr[T, E, F any](key Key) *SlicePtr[T, E] {
	return &SlicePtr[T, E]{
		Slice:    &Slice[T, E]{key: key},
		Nillable: &Nillable[T]{key: key},
	}
}
