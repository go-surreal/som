//go:build embed

package lib

type int_ interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~*int | ~*int8 | ~*int16 | ~*int32 | ~*int64 |
		~*uint | ~*uint8 | ~*uint16 | ~*uint32 | ~*uint64 | ~*uintptr
}

type AnyInt[M any] interface {
	field[M]
	anyInt()
}

type Int[M any, T int_] struct {
	*Numeric[M, T]
}

func NewInt[M any, T int_](key Key[M]) *Int[M, T] {
	return &Int[M, T]{
		Numeric: NewNumeric[M, T](key),
	}
}

type IntPtr[M any, T int_] struct {
	*Int[M, T]
	*Nillable[M]
}

func NewIntPtr[M any, T int_](key Key[M]) *IntPtr[M, T] {
	return &IntPtr[M, T]{
		Int:      NewInt[M, T](key),
		Nillable: NewNillable(key),
	}
}

func (n *Int[M, T]) anyInt() {}

func (n *Int[M, T]) Int() *Int[M, int] {
	return NewInt[M, int](n.key())
}

func (n *Int[M, T]) Int8() *Int[M, int8] {
	return NewInt[M, int8](n.key())
}

func (n *Int[M, T]) Int16() *Int[M, int16] {
	return NewInt[M, int16](n.key())
}

func (n *Int[M, T]) Int32() *Int[M, int32] {
	return NewInt[M, int32](n.key())
}

func (n *Int[M, T]) Int64() *Int[M, int64] {
	return NewInt[M, int64](n.key())
}

func (n *Int[M, T]) Uint() *Int[M, uint] {
	return NewInt[M, uint](n.key())
}

func (n *Int[M, T]) Uint8() *Int[M, uint8] {
	return NewInt[M, uint8](n.key())
}

func (n *Int[M, T]) Uint16() *Int[M, uint16] {
	return NewInt[M, uint16](n.key())
}

func (n *Int[M, T]) Uint32() *Int[M, uint32] {
	return NewInt[M, uint32](n.key())
}

func (n *Int[M, T]) Uint64() *Int[M, uint64] {
	return NewInt[M, uint64](n.key())
}

func (n *Int[M, T]) Float32() *Float[M, float32] {
	return NewFloat[M, float32](n.Base.prefix(CastFloat))
}

func (n *Int[M, T]) Float64() *Float[M, float64] {
	return NewFloat[M, float64](n.Base.prefix(CastFloat))
}
