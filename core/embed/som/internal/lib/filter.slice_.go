//go:build embed

package lib

func (s *Slice[M, E, F]) Contains_(field F) Filter[M] {
	return s.op_(OpContains, field.key())
}

func (s *Slice[M, E, F]) ContainsNot_(field F) Filter[M] {
	return s.op_(OpContainsNot, field.key())
}

func (s *Slice[M, E, F]) ContainsAll_(field *Slice[M, E, F]) Filter[M] {
	return s.op_(OpContainsAll, field.key())
}

func (s *Slice[M, E, F]) ContainsAny_(field *Slice[M, E, F]) Filter[M] {
	return s.op_(OpContainsAny, field.key())
}

func (s *Slice[M, E, F]) ContainsNone_(field *Slice[M, E, F]) Filter[M] {
	return s.op_(OpContainsNone, field.key())
}

func (s *Slice[M, E, F]) At_(field *Numeric[M, int]) F {
	return s.makeElemFilter(s.fn_("array::at", field.key()))
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

func (s *Slice[M, E, F]) Matches_(field F) *Slice[M, bool, *Bool[M]] {
	return NewSlice[M, bool, *Bool[M]](
		s.fn_("array::matches", field.key()),
		NewBool[M],
	)
}

func (s *Slice[M, E, F]) Slice_(start, len *Numeric[M, int]) *Slice[M, E, F] {
	return NewSlice[M, E, F](
		s.fn_("array::slice", start.key(), len.key()),
		s.makeElemFilter,
	)
}
