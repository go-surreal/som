//go:build embed

package lib

type makeFilter[M, F any] func(key Key[M]) F

func NewSliceMaker[M, E, F any](makeElemFilter makeFilter[M, F]) makeFilter[M, *Slice[M, E, F]] {
	return func(key Key[M]) *Slice[M, E, F] {
		return NewSlice[M, E, F](key, makeElemFilter)
	}
}

// Slice is a filter that can be used for slice fields.
// M is the type of the outgoing model for the filter statement.
// E is the type of the slice elements.
// F is the type of the element filter.
type Slice[M, E, F any] struct {
	key            Key[M]
	makeElemFilter makeFilter[M, F]
}

// NewSlice creates a new slice filter.
func NewSlice[M, E, F any](key Key[M], makeElemFilter makeFilter[M, F]) *Slice[M, E, F] {
	return &Slice[M, E, F]{
		key:            key,
		makeElemFilter: makeElemFilter,
	}
}

func (s *Slice[M, E, F]) Contains(val E) Filter[M] {
	return s.key.op(OpContains, val)
}

func (s *Slice[M, E, F]) ContainsNot(val E) Filter[M] {
	return s.key.op(OpContainsNot, val)
}

func (s *Slice[M, E, F]) ContainsAll(vals []E) Filter[M] {
	return s.key.op(OpContainsAll, vals)
}

func (s *Slice[M, E, F]) ContainsAny(vals []E) Filter[M] {
	return s.key.op(OpContainsAny, vals)
}

func (s *Slice[M, E, F]) ContainsNone(vals []E) Filter[M] {
	return s.key.op(OpContainsNone, vals)
}

func (s *Slice[M, E, F]) All() *Bool[M] {
	return NewBool[M](s.key.fn("array::all"))
}

func (s *Slice[M, E, F]) Any() *Bool[M] {
	return NewBool[M](s.key.fn("array::any"))
}

func (s *Slice[M, E, F]) At(index int) F {
	return s.makeElemFilter(s.key.fn("array::at", index))
}

func (s *Slice[M, E, F]) Distinct() *Slice[M, E, F] {
	return NewSlice[M, E, F](s.key.fn("array::distinct"), s.makeElemFilter)
}

func (s *Slice[M, E, F]) FindIndex(val E) *Numeric[M, int] {
	return NewNumeric[M, int](s.key.fn("array::find_index", val))
}

func (s *Slice[M, E, F]) FilterIndex(val E) *Slice[M, int, *Numeric[M, int]] {
	return NewSlice[M, int, *Numeric[M, int]](
		s.key.fn("array::filter_index"),
		NewNumeric[M, int],
	)
}

func (s *Slice[M, E, F]) First() F {
	return s.makeElemFilter(s.key.fn("array::first"))
}

func (s *Slice[M, E, F]) Last() F {
	return s.makeElemFilter(s.key.fn("array::last"))
}

func (s *Slice[M, E, F]) Len() *Numeric[M, int] {
	return NewNumeric[M, int](s.key.fn("array::len"))
}

func (s *Slice[M, E, F]) Max() F {
	return s.makeElemFilter(s.key.fn("array::max"))
}

func (s *Slice[M, E, F]) Matches() *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](
		s.key.fn("array::matches"),
		NewBool[M],
	)
}

func (s *Slice[M, E, F]) Min() F {
	return s.makeElemFilter(s.key.fn("array::min"))
}

func (s *Slice[M, E, F]) Reverse() *Slice[M, E, F] {
	return NewSlice[M, E, F](s.key.fn("array::reverse"), s.makeElemFilter)
}

func (s *Slice[M, E, F]) SortAsc() *Slice[M, E, F] {
	return NewSlice[M, E, F](s.key.fn("array::sort::asc"), s.makeElemFilter)
}

func (s *Slice[M, E, F]) SortDesc() *Slice[M, E, F] {
	return NewSlice[M, E, F](s.key.fn("array::sort::desc"), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Slice(start, len int) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.key.fn("array::slice", start, len), s.makeElemFilter)
}

type SlicePtr[M, E any, F filter[M]] struct {
	*Slice[M, E, F]
	*Nillable[M]
}

func NewSlicePtr[M, E any, F filter[M]](key Key[M], makeElemFilter makeFilter[M, F]) *SlicePtr[M, E, F] {
	return &SlicePtr[M, E, F]{
		Slice:    NewSlice[M, E, F](key, makeElemFilter),
		Nillable: NewNillable[M](key),
	}
}
