package lib

type Sort[T any] SortBuilder

type BaseSort[T any] struct {
	key string
}

func NewBaseSort[T any](key string) *BaseSort[T] {
	return &BaseSort[T]{key: key}
}

func (s BaseSort[T]) Asc() *Sort[T] {
	return &Sort[T]{Field: s.key, Order: SortAsc}
}

func (s BaseSort[T]) Desc() *Sort[T] {
	return &Sort[T]{Field: s.key, Order: SortDesc}
}

type StringSort[T any] struct {
	key string
}

func NewStringSort[T any](key string) *StringSort[T] {
	return &StringSort[T]{key: key}
}

func (s StringSort[T]) Asc() *Sort[T] {
	return &Sort[T]{Field: s.key, Order: SortAsc}
}

func (s StringSort[T]) Desc() *Sort[T] {
	return &Sort[T]{Field: s.key, Order: SortDesc}
}

func (s StringSort[T]) Collate() StringCollate[T] {
	return StringCollate[T](s)
}

func (s StringSort[T]) Numeric() StringNumeric[T] {
	return StringNumeric[T](s)
}

type StringCollate[T any] struct {
	key string
}

func (w StringCollate[T]) Asc() *Sort[T] {
	return &Sort[T]{Field: w.key, Order: SortAsc, IsCollate: true}
}

func (w StringCollate[T]) Desc() *Sort[T] {
	return &Sort[T]{Field: w.key, Order: SortDesc, IsCollate: true}
}

type StringNumeric[T any] struct {
	key string
}

func (w StringNumeric[T]) Asc() *Sort[T] {
	return &Sort[T]{Field: w.key, Order: SortAsc, IsNumeric: true}
}

func (w StringNumeric[T]) Desc() *Sort[T] {
	return &Sort[T]{Field: w.key, Order: SortDesc, IsNumeric: true}
}
