package filter

type Base[T any, R any] struct {
	key string
}

func NewBase[T, R any](key string) *Base[T, R] {
	return &Base[T, R]{key: key}
}

func (Base[T, R]) Equal(val T) *Of[R] {
	return &Of[R]{}
}

func (Base[T, R]) NotEqual(val T) *Of[R] {
	return &Of[R]{}
}

func (Base[T, R]) In(vals []T) *Of[R] {
	return &Of[R]{}
}

func (Base[T, R]) NotIn(vals []T) *Of[R] {
	return &Of[R]{}
}
