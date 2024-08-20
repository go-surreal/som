//go:build embed

package lib

type FloatSlice[M any, T float_] struct {
	*NumericSlice[M, T, *Float[M, T]]
}

func NewFloatSlice[M any, T float_](key Key[M]) *FloatSlice[M, T] {
	return &FloatSlice[M, T]{
		NumericSlice: NewNumericSlice[M, T](key, NewFloat[M, T]),
	}
}

func (s *FloatSlice[M, T]) Bottom(count int) *FloatSlice[M, T] {
	return NewFloatSlice[M, T](s.fn("math::bottom", count))
}

func (s *FloatSlice[M, T]) Max() *Float[M, T] {
	return NewFloat[M, T](s.fn("math::abs"))
}

func (s *FloatSlice[M, T]) Min() *Float[M, T] {
	return NewFloat[M, T](s.fn("math::min"))
}

func (s *FloatSlice[M, T]) Mode() *Float[M, T] {
	return NewFloat[M, T](s.fn("math::mode"))
}

func (s *FloatSlice[M, T]) Product() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::product"))
}

func (s *FloatSlice[M, T]) Spread() *Float[M, T] {
	return NewFloat[M, T](s.fn("math::spread"))
}

func (s *FloatSlice[M, T]) Sum() *Float[M, float64] {
	return NewFloat[M, float64](s.fn("math::sum"))
}

func (s *FloatSlice[M, T]) Top(count int) *FloatSlice[M, T] {
	return NewFloatSlice[M, T](s.fn("math::top", count))
}

//
// -- POINTER
//

type FloatSlicePtr[M any, T float_] struct {
	*NumericSlicePtr[M, T, *Float[M, T]]
}

func NewFloatSlicePtr[M any, T float_](key Key[M]) *FloatSlicePtr[M, T] {
	return &FloatSlicePtr[M, T]{
		NumericSlicePtr: NewNumericSlicePtr[M, T](key, NewFloat[M, T]),
	}
}
