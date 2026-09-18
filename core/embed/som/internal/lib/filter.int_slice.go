//go:build embed

package lib

// The integer slice filters drop the element width for the same reason Int
// does: it only ever typed the comparison arguments, and instantiated the
// whole Slice machinery once per width per model.

type IntSlice[M any] struct {
	*NumericSlice[M, int64, *Int[M]]
}

func NewIntSlice[M any](key Key[M]) *IntSlice[M] {
	return &IntSlice[M]{
		NumericSlice: NewNumericSlice[M, int64](key, NewInt[M]),
	}
}

type IntSlicePtr[M any] struct {
	*NumericSlicePtr[M, int64, *Int[M]]
}

func NewIntSlicePtr[M any](key Key[M]) *IntSlicePtr[M] {
	return &IntSlicePtr[M]{
		NumericSlicePtr: NewNumericSlicePtr[M, int64](key, NewInt[M]),
	}
}

type IntPtrSlice[M any] struct {
	*NumericSlice[M, int64, *IntPtr[M]]
}

func NewIntPtrSlice[M any](key Key[M]) *IntPtrSlice[M] {
	return &IntPtrSlice[M]{
		NumericSlice: NewNumericSlice[M, int64](key, NewIntPtr[M]),
	}
}

type IntPtrSlicePtr[M any] struct {
	*NumericSlicePtr[M, int64, *IntPtr[M]]
}

func NewIntPtrSlicePtr[M any](key Key[M]) *IntPtrSlicePtr[M] {
	return &IntPtrSlicePtr[M]{
		NumericSlicePtr: NewNumericSlicePtr[M, int64](key, NewIntPtr[M]),
	}
}

//
//
//

func (s *IntSlice[M]) Bottom(count int) *IntSlice[M] {
	return NewIntSlice[M](s.fn("math::bottom", count))
}

func (s *IntSlice[M]) Max() *Int[M] {
	return NewInt[M](s.fn("math::abs"))
}

func (s *IntSlice[M]) Min() *Int[M] {
	return NewInt[M](s.fn("math::min"))
}

func (s *IntSlice[M]) Mode() *Int[M] {
	return NewInt[M](s.fn("math::mode"))
}

func (s *IntSlice[M]) Product() *Int[M] {
	return NewInt[M](s.fn("math::product"))
}

func (s *IntSlice[M]) Spread() *Int[M] {
	return NewInt[M](s.fn("math::spread"))
}

func (s *IntSlice[M]) Sum() *Int[M] {
	return NewInt[M](s.fn("math::sum"))
}

func (s *IntSlice[M]) Top(count int) *IntSlice[M] {
	return NewIntSlice[M](s.fn("math::top", count))
}
