package filter

type Slice[T any, R any] struct {
	key string
}

func NewSlice[T, R any](key string) *Slice[T, R] {
	return &Slice[T, R]{key: key}
}

func (Slice[T, R]) Contains(val T) *Of[R] {
	return &Of[R]{}
}

func (Slice[T, R]) ContainsNot(val T) *Of[R] {
	return &Of[R]{}
}

func (Slice[T, R]) ContainsAll(vals []T) *Of[R] {
	return &Of[R]{}
}

func (Slice[T, R]) ContainsAny(val []T) *Of[R] {
	return &Of[R]{}
}

func (Slice[T, R]) ContainsNone(val []T) *Of[R] {
	return &Of[R]{}
}

func (s Slice[T, R]) Count() *Numeric[int, R] {
	return NewNumeric[int, R](s.key) // TODO: pass info that this is a "COUNT"
}
