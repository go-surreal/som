package filter

type Comparable[T any, R any] struct {
	key string
}

func (Comparable[T, R]) LessThan(val T) *Of[R] {
	return &Of[R]{}
}

func (Comparable[T, R]) LessThanEqual(val T) *Of[R] {
	return &Of[R]{}
}

func (Comparable[T, R]) GreaterThan(val T) *Of[R] {
	return &Of[R]{}
}

func (Comparable[T, R]) GreaterThanEqual(val T) *Of[R] {
	return &Of[R]{}
}
