//go:build embed

package lib

type makeFilter[M, F any] func(key Key[M]) F

// Slice is a filter that can be used for slice fields.
// M is the type of the outgoing model for the filter statement.
// E is the type of the slice elements.
// F is the type of the element filter.
type Slice[M, E any, F field[M]] struct {
	Key[M]
	makeElemFilter makeFilter[M, F]
}

// NewSlice creates a new slice filter.
func NewSlice[M, E any, F field[M]](key Key[M], makeElemFilter makeFilter[M, F]) *Slice[M, E, F] {
	return &Slice[M, E, F]{
		Key:            key,
		makeElemFilter: makeElemFilter,
	}
}

type SlicePtr[M, E any, F field[M]] struct {
	*Slice[M, E, F]
	*Nillable[M]
}

func NewSlicePtr[M, E any, F field[M]](key Key[M], makeElemFilter makeFilter[M, F]) *SlicePtr[M, E, F] {
	return &SlicePtr[M, E, F]{
		Slice:    NewSlice[M, E, F](key, makeElemFilter),
		Nillable: NewNillable[M](key),
	}
}

func NewSliceMaker[M, E any, F field[M]](makeElemFilter makeFilter[M, F]) makeFilter[M, *Slice[M, E, F]] {
	return func(key Key[M]) *Slice[M, E, F] {
		return NewSlice[M, E, F](key, makeElemFilter)
	}
}

func NewSliceMakerPtr[M, E any, F field[M]](makeElemFilter makeFilter[M, F]) makeFilter[M, *SlicePtr[M, E, F]] {
	return func(key Key[M]) *SlicePtr[M, E, F] {
		return NewSlicePtr[M, E, F](key, makeElemFilter)
	}
}

//
// -- COMPARISONS
//

func (s *Slice[M, E, F]) AnyEqual(val E) Filter[M] {
	return s.op(OpAnyEqual, val)
}

func (s *Slice[M, E, F]) AllEqual(val E) Filter[M] {
	return s.op(OpAllEqual, val)
}

func (s *Slice[M, E, F]) AnyFuzzyMatch(val E) Filter[M] {
	return s.op(OpAnyFuzzyMatch, val)
}

func (s *Slice[M, E, F]) AllFuzzyMatch(val E) Filter[M] {
	return s.op(OpAllFuzzyMatch, val)
}

func (s *Slice[M, E, F]) AllIn(val []E) Filter[M] {
	return s.op(OpAllIn, val)
}

func (s *Slice[M, E, F]) AnyIn(val []E) Filter[M] {
	return s.op(OpAnyIn, val)
}

func (s *Slice[M, E, F]) NoneIn(val []E) Filter[M] {
	return s.op(OpNoneIn, val)
}

func (s *Slice[M, E, F]) Contains(val E) Filter[M] {
	return s.op(OpContains, val)
}

func (s *Slice[M, E, F]) ContainsNot(val E) Filter[M] {
	return s.op(OpContainsNot, val)
}

func (s *Slice[M, E, F]) ContainsAll(vals []E) Filter[M] {
	return s.op(OpContainsAll, vals)
}

func (s *Slice[M, E, F]) ContainsAny(vals []E) Filter[M] {
	return s.op(OpContainsAny, vals)
}

func (s *Slice[M, E, F]) ContainsNone(vals []E) Filter[M] {
	return s.op(OpContainsNone, vals)
}

func (s *Slice[M, E, F]) Add(val E) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::add", val), s.makeElemFilter)
}

func (s *Slice[M, E, F]) All() *Bool[M] {
	return NewBool[M](s.fn("array::all"))
}

func (s *Slice[M, E, F]) Any() *Bool[M] {
	return NewBool[M](s.fn("array::any"))
}

func (s *Slice[M, E, F]) At(index int) F {
	return s.makeElemFilter(s.fn("array::at", index))
}

func (s *Slice[M, E, F]) Append(val E) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::append", val), s.makeElemFilter)
}

func (s *Slice[M, E, F]) BooleanAnd(val []E) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](s.fn("array::boolean_and", val), NewBool)
}

func (s *Slice[M, E, F]) BooleanOr(val []E) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](s.fn("array::boolean_or", val), NewBool)
}

func (s *Slice[M, E, F]) BooleanXor(val []E) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](s.fn("array::boolean_xor", val), NewBool)
}

func (s *Slice[M, E, F]) BooleanNot(val []E) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](s.fn("array::boolean_not", val), NewBool)
}

//func (s *Slice[M, E, F]) Combine(val []E) *Slice[M, []E, *Slice[M, E, F]] {
//	return NewSlice[M, []E, *Slice[M, E, F]](s.fn("array::combine", val), NewSliceMaker[M, E, F](s.makeElemFilter))
//} --> error: instantiation cycle

func (s *Slice[M, E, F]) Complement(val []E) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::complement", val), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Concat(val []E) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::concat", val), s.makeElemFilter)
}

//func (s *Slice[M, E, F]) Clump(size int) *Slice[M, []E, *Slice[M, E, F]] {
//	return NewSlice[M, []E, *Slice[M, E, F]](s.fn("array::clump", size), NewSliceMaker[M, E, F](s.makeElemFilter))
//} --> error: instantiation cycle

func (s *Slice[M, E, F]) Diff(val []E) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::difference", val), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Distinct() *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::distinct"), s.makeElemFilter)
}

// func (s *Slice[M, E, F]) Flatten() *Slice[M, E, F] {} -> only works for arrays of arrays

func (s *Slice[M, E, F]) FindIndex(val E) *Numeric[M, int] {
	return NewNumeric[M, int](s.fn("array::find_index", val))
}

func (s *Slice[M, E, F]) FilterIndex(val E) *Slice[M, int, *Numeric[M, int]] {
	return NewSlice[M, int, *Numeric[M, int]](
		s.fn("array::filter_index", val),
		NewNumeric[M, int],
	)
}

func (s *Slice[M, E, F]) First() F {
	return s.makeElemFilter(s.fn("array::first"))
}

// func (s *Slice[M, E, F]) Group() *Slice[M, E, F] {} -> only works for arrays of arrays

func (s *Slice[M, E, F]) Insert(val E, pos int) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::insert", val, pos), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Intersect(val []E) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::intersect", val), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Join(sep string) *String[M] {
	return NewString[M](s.fn("array::join", sep)) // TODO: works for any type?
}

func (s *Slice[M, E, F]) Last() F {
	return s.makeElemFilter(s.fn("array::last"))
}

func (s *Slice[M, E, F]) Len() *Numeric[M, int] {
	return NewNumeric[M, int](s.fn("array::len"))
}

func (s *Slice[M, E, F]) LogicalAnd(val []E) *Slice[M, bool, *Bool[M]] {
	// TODO: return type might vary if the arrays are different in length
	return NewSlice[M, bool, *Bool[M]](s.fn("array::logical_and", val), NewBool)
}

func (s *Slice[M, E, F]) LogicalOr(val []E) *Slice[M, bool, *Bool[M]] {
	// TODO: return type might vary if the arrays are different in length
	return NewSlice[M, bool, *Bool[M]](s.fn("array::logical_or", val), NewBool)
}

func (s *Slice[M, E, F]) LogicalXor(val []E) *Slice[M, bool, *Bool[M]] {
	// TODO: return type might vary if the arrays are different in length
	return NewSlice[M, bool, *Bool[M]](s.fn("array::logical_xor", val), NewBool)
}

func (s *Slice[M, E, F]) Max() F {
	return s.makeElemFilter(s.fn("array::max"))
}

func (s *Slice[M, E, F]) Matches(val E) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](s.fn("array::matches", val), NewBool[M])
}

func (s *Slice[M, E, F]) Min() F {
	return s.makeElemFilter(s.fn("array::min"))
}

// func (s *Slice[M, E, F]) Pop() F {} -> needed? - might modify the model?!

func (s *Slice[M, E, F]) Prepend(val E) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::prepend", val), s.makeElemFilter)
}

// func (s *Slice[M, E, F]) Push(val E) F {} -> needed? - might modify the model?!

// func (s *Slice[M, E, F]) Remove(pos int) F {} -> needed? - might modify the model?!

func (s *Slice[M, E, F]) Reverse() *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::reverse"), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Slice(start, len int) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::slice", start, len), s.makeElemFilter)
}

func (s *Slice[M, E, F]) SortAsc() *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::sort::asc"), s.makeElemFilter)
}

func (s *Slice[M, E, F]) SortDesc() *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::sort::desc"), s.makeElemFilter)
}

// TODO: Transpose

func (s *Slice[M, E, F]) Union(val []E) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn("array::union", val), s.makeElemFilter)
}

// TODO: Windows (v2.0.0)

// TODO: https://surrealdb.com/docs/surrealdb/surrealql/functions/database/vector
