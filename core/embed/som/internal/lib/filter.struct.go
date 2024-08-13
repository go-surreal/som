//go:build embed

package lib

// Struct is a filter that can be used for struct fields.
// M is the type of the model for the filter statement.
// T is the type of the struct.
type Struct[M, T any] struct {
	key Key[M]
}

// NewStruct creates a new struct filter.
func NewStruct[M, T any](key Key[M]) *Struct[M, T] {
	return &Struct[M, T]{
		key: key,
	}
}

//func (s *Struct[M, T]) Entries() [][]any {}

func (s *Struct[M, T]) Keys() *Slice[M, string, *String[M]] {
	return NewSlice[M, string, *String[M]](s.key.fn("object::keys"), NewString[M])
}

//func (s *Struct[M, T]) Values() any {}

func (s *Struct[M, T]) Len() *Numeric[M, int] {
	return NewNumeric[M, int](s.key.fn("object::len"))
}
