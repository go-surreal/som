package filter

type Bool[R any] struct {
	key string
}

func NewBool[R any](key string) *Bool[R] {
	return &Bool[R]{key: key}
}

func (Bool[R]) Is(val bool) *Of[R] {
	return &Of[R]{}
}
