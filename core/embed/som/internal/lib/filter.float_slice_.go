//go:build embed

package lib

func (s *FloatSlice[M, T]) Bottom_(count AnyInt[M]) *FloatSlice[M, T] {
	return NewFloatSlice[M, T](s.fn_("math::bottom", count.key()))
}

func (s *FloatSlice[M, T]) Top_(count AnyInt[M]) *FloatSlice[M, T] {
	return NewFloatSlice[M, T](s.fn_("math::top", count.key()))
}
