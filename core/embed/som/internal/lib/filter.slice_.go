//go:build embed

package lib

func (s *Slice[M, E, F]) AnyEqual_(val F) Filter[M] {
	return s.op_(OpAnyEqual, val.key())
}

func (s *Slice[M, E, F]) AllEqual_(val F) Filter[M] {
	return s.op_(OpAllEqual, val.key())
}

func (s *Slice[M, E, F]) AnyFuzzyMatch_(val F) Filter[M] {
	return s.op_(OpAnyFuzzyMatch, val.key())
}

func (s *Slice[M, E, F]) AllFuzzyMatch_(val F) Filter[M] {
	return s.op_(OpAllFuzzyMatch, val.key())
}

func (s *Slice[M, E, F]) AllIn_(val *Slice[M, E, F]) Filter[M] {
	return s.op_(OpAllIn, val.key())
}

func (s *Slice[M, E, F]) AnyIn_(val *Slice[M, E, F]) Filter[M] {
	return s.op_(OpAnyIn, val.key())
}

func (s *Slice[M, E, F]) NoneIn_(val *Slice[M, E, F]) Filter[M] {
	return s.op_(OpNoneIn, val.key())
}

func (s *Slice[M, E, F]) Contains_(val F) Filter[M] {
	return s.op_(OpContains, val.key())
}

func (s *Slice[M, E, F]) ContainsNot_(val F) Filter[M] {
	return s.op_(OpContainsNot, val.key())
}

func (s *Slice[M, E, F]) ContainsAll_(val *Slice[M, E, F]) Filter[M] {
	return s.op_(OpContainsAll, val.key())
}

func (s *Slice[M, E, F]) ContainsAny_(val *Slice[M, E, F]) Filter[M] {
	return s.op_(OpContainsAny, val.key())
}

func (s *Slice[M, E, F]) ContainsNone_(val *Slice[M, E, F]) Filter[M] {
	return s.op_(OpContainsNone, val.key())
}

func (s *Slice[M, E, F]) Add_(val F) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn_("array::add", val.key()), s.makeElemFilter)
}

func (s *Slice[M, E, F]) At_(pos *Numeric[M, int]) F { // TODO: type is NONE if out of bounds
	return s.makeElemFilter(s.fn_("array::at", pos.key()))
}

func (s *Slice[M, E, F]) Append_(val F) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn_("array::append", val.key()), s.makeElemFilter)
}

func (s *Slice[M, E, F]) BooleanAnd_(val *Slice[M, E, F]) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](s.fn_("array::boolean_and", val.key()), NewBool)
}

func (s *Slice[M, E, F]) BooleanOr_(val *Slice[M, E, F]) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](s.fn_("array::boolean_or", val.key()), NewBool)
}

func (s *Slice[M, E, F]) BooleanXor_(val *Slice[M, E, F]) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](s.fn_("array::boolean_xor", val.key()), NewBool)
}

func (s *Slice[M, E, F]) BooleanNot_(val *Slice[M, E, F]) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](s.fn_("array::boolean_not", val.key()), NewBool)
}

func (s *Slice[M, E, F]) Complement_(val *Slice[M, E, F]) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn_("array::complement", val.key()), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Concat_(val *Slice[M, E, F]) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn_("array::concat", val.key()), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Diff_(val *Slice[M, E, F]) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn_("array::difference", val.key()), s.makeElemFilter)
}

func (s *Slice[M, E, F]) FindIndex_(field F) *Numeric[M, int] {
	return NewNumeric[M, int](s.fn_("array::find_index", field.key()))
}

func (s *Slice[M, E, F]) FilterIndex_(field F) *Slice[M, int, *Numeric[M, int]] {
	return NewSlice[M, int, *Numeric[M, int]](
		s.fn_("array::filter_index", field.key()),
		NewNumeric[M, int],
	)
}

func (s *Slice[M, E, F]) Insert_(val F, pos AnyInt[M]) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn_("array::insert", val.key(), pos.key()), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Intersect_(val *Slice[M, E, F]) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn_("array::intersect", val.key()), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Join_(sep *String[M]) *String[M] {
	return NewString[M](s.fn_("array::join", sep.key())) // TODO: works for any type?
}

func (s *Slice[M, E, F]) LogicalAnd_(val *Slice[M, E, F]) *Slice[M, bool, *Bool[M]] {
	// TODO: return type might vary if the arrays are different in length
	return NewSlice[M, bool, *Bool[M]](s.fn_("array::logical_and", val.key()), NewBool)
}

func (s *Slice[M, E, F]) LogicalOr_(val *Slice[M, E, F]) *Slice[M, bool, *Bool[M]] {
	// TODO: return type might vary if the arrays are different in length
	return NewSlice[M, bool, *Bool[M]](s.fn_("array::logical_or", val.key()), NewBool)
}

func (s *Slice[M, E, F]) LogicalXor_(val *Slice[M, E, F]) *Slice[M, bool, *Bool[M]] {
	// TODO: return type might vary if the arrays are different in length
	return NewSlice[M, bool, *Bool[M]](s.fn_("array::logical_xor", val.key()), NewBool)
}

func (s *Slice[M, E, F]) Matches_(val F) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](
		s.fn_("array::matches", val.key()),
		NewBool[M],
	)
}

func (s *Slice[M, E, F]) Prepend_(val F) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn_("array::prepend", val.key()), s.makeElemFilter)
}

func (s *Slice[M, E, F]) Slice_(start, len AnyInt[M]) *Slice[M, E, F] {
	return NewSlice[M, E, F](
		s.fn_("array::slice", start.key(), len.key()),
		s.makeElemFilter,
	)
}

func (s *Slice[M, E, F]) Union_(val *Slice[M, E, F]) *Slice[M, E, F] {
	return NewSlice[M, E, F](s.fn_("array::union", val.key()), s.makeElemFilter)
}
