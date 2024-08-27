//go:build embed

package lib

func (s *String[M]) FuzzyMatch_(field *String[M]) Filter[M] {
	return s.Base.op_(OpFuzzyMatch, field.Base.Key)
}

func (s *String[M]) NotFuzzyMatch_(field *String[M]) Filter[M] {
	return s.Base.op_(OpFuzzyNotMatch, field.Base.Key)
}

func (s *String[M]) Concat_(vals ...*String[M]) *Bool[M] {
	keys := make([]Key[M], len(vals))

	for i, val := range vals {
		keys[i] = val.key()
	}

	return NewBool(s.Base.fn_("string::concat", keys...))
}

func (s *String[M]) Contains_(field *String[M]) *Bool[M] {
	return NewBool(s.Base.fn_("string::contains", field.Base.Key))
}

func (s *String[M]) EndsWith_(field *String[M]) *Bool[M] {
	return NewBool(s.Base.fn_("string::endsWith", field.Base.Key))
}

// Join_ joins the given string fields with the base field as separator.
func (s *String[M]) Join_(vals ...*String[M]) *String[M] {
	keys := make([]Key[M], len(vals))

	for i, val := range vals {
		keys[i] = val.key()
	}

	return NewString(s.Base.fn_("string::join", keys...))
}

func (s *String[M]) Repeat_(times AnyInt[M]) *String[M] {
	return NewString(s.Base.fn_("string::repeat", times.key()))
}

func (s *String[M]) Replace_(old, new *String[M]) *String[M] {
	return NewString(s.Base.fn_("string::replace", old.key(), new.key()))
}

func (s *String[M]) Slice_(start, end *Numeric[M, int]) *String[M] {
	return NewString(s.Base.fn_("string::slice", start.Base.Key, end.Base.Key))
}

func (s *String[M]) Split_(field *String[M]) *Slice[M, string, *String[M]] {
	return NewSlice[M, string, *String[M]](s.Base.fn_("string::split", field.Base.Key), NewString[M])
}

func (s *String[M]) StartsWith_(field *String[M]) *Bool[M] {
	return NewBool(s.Base.fn_("string::startsWith", field.Base.Key))
}
