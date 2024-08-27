//go:build embed

package lib

func (t *Time[M]) Before_(field *Time[M]) Filter[M] {
	return t.comp.LessThan_(field)
}

func (t *Time[M]) BeforeOrEqual_(field *Time[M]) Filter[M] {
	return t.comp.LessThanEqual_(field)
}

func (t *Time[M]) After_(field *Time[M]) Filter[M] {
	return t.comp.GreaterThan_(field)
}

func (t *Time[M]) AfterOrEqual_(field *Time[M]) Filter[M] {
	return t.comp.GreaterThanEqual_(field)
}

func (t *Time[M]) Add_(dur *Duration[M]) *Time[M] {
	return NewTime[M](t.calc_(OpAdd, dur.key()))
}

func (t *Time[M]) Sub_(dur *Duration[M]) *Time[M] {
	return NewTime[M](t.calc_(OpSub, dur.key()))
}

func (t *Time[M]) Floor_(field *Duration[M]) *Time[M] {
	return NewTime[M](t.fn_("time::floor", field.key()))
}

func (t *Time[M]) Format_(field *String[M]) *String[M] {
	return NewString[M](t.fn_("time::format", field.key()))
}

// TODO!
//func (t *Time[M]) Group_(group Group) *Time[M] {
//	return NewTime[M](t.fn("time::group", group))
//}

func (t *Time[M]) Round_(field *Duration[M]) *Time[M] {
	return NewTime[M](t.fn_("time::round", field.key()))
}
