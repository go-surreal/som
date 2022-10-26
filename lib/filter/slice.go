package filter

type Slice[T any, R any] struct {
	*Base[T, R]
}

func NewSlice[T, R any](key string) *Slice[T, R] {
	return &Slice[T, R]{
		Base: &Base[T, R]{key: key},
	}
}

func (String[R]) Contains(val string) *Of[R] {
	return &Of[R]{}
}

func (String[R]) ContainsNot(val string) *Of[R] {
	return &Of[R]{}
}

func (String[R]) ContainsAll(vals []string) *Of[R] {
	return &Of[R]{}
}

func (String[R]) ContainsAny(val string) *Of[R] {
	return &Of[R]{}
}

func (String[R]) ContainsNone(val string) *Of[R] {
	return &Of[R]{}
}
