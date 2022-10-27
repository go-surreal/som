package sort

type Of[T any] struct{}

type Sort[T any] struct {
	key string
}

func NewSort[T any](key string) *Sort[T] {
	return &Sort[T]{key: key}
}

func (w Sort[T]) Asc() *Of[T] {
	return &Of[T]{}
}

func (w Sort[T]) Desc() *Of[T] {
	return &Of[T]{}
}

type String[T any] struct {
	key string
}

func NewString[T any](key string) *String[T] {
	return &String[T]{key: key}
}

func (w String[T]) Asc() *Of[T] {
	return &Of[T]{}
}

func (w String[T]) Desc() *Of[T] {
	return &Of[T]{}
}

func (w String[T]) Collate() StringCollate[T] {
	return StringCollate[T]{w.key}
}

func (w String[T]) Numeric() StringNumeric[T] {
	return StringNumeric[T]{w.key}
}

type StringCollate[T any] struct {
	key string
}

func (w StringCollate[T]) Asc() *Of[T] {
	return &Of[T]{}
}

func (w StringCollate[T]) Desc() *Of[T] {
	return &Of[T]{}
}

type StringNumeric[T any] struct {
	Origin string
}

func (w StringNumeric[T]) Asc() *Of[T] {
	return &Of[T]{}
}

func (w StringNumeric[T]) Desc() *Of[T] {
	return &Of[T]{}
}
