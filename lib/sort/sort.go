package sort

import (
	"github.com/marcbinz/som/lib/builder"
)

type Of[T any] builder.Sort

type Sort[T any] struct {
	key string
}

func NewSort[T any](key string) *Sort[T] {
	return &Sort[T]{key: key}
}

func (w Sort[T]) Asc() *Of[T] {
	return &Of[T]{Field: w.key, Order: builder.SortAsc}
}

func (w Sort[T]) Desc() *Of[T] {
	return &Of[T]{Field: w.key, Order: builder.SortDesc}
}

type String[T any] struct {
	key string
}

func NewString[T any](key string) *String[T] {
	return &String[T]{key: key}
}

func (w String[T]) Asc() *Of[T] {
	return &Of[T]{Field: w.key, Order: builder.SortAsc}
}

func (w String[T]) Desc() *Of[T] {
	return &Of[T]{Field: w.key, Order: builder.SortDesc}
}

func (w String[T]) Collate() StringCollate[T] {
	return StringCollate[T](w)
}

func (w String[T]) Numeric() StringNumeric[T] {
	return StringNumeric[T](w)
}

type StringCollate[T any] struct {
	key string
}

func (w StringCollate[T]) Asc() *Of[T] {
	return &Of[T]{Field: w.key, Order: builder.SortAsc, IsCollate: true}
}

func (w StringCollate[T]) Desc() *Of[T] {
	return &Of[T]{Field: w.key, Order: builder.SortDesc, IsCollate: true}
}

type StringNumeric[T any] struct {
	key string
}

func (w StringNumeric[T]) Asc() *Of[T] {
	return &Of[T]{Field: w.key, Order: builder.SortAsc, IsNumeric: true}
}

func (w StringNumeric[T]) Desc() *Of[T] {
	return &Of[T]{Field: w.key, Order: builder.SortDesc, IsNumeric: true}
}
