//go:build embed

package lib

func (s *IntSlice[M, T]) Bottom_(count AnyInt[M]) *IntSlice[M, T] {
	return NewIntSlice[M, T](s.fn_("math::bottom", count.key()))
}

func (s *IntSlice[M, T]) Top_(count AnyInt[M]) *IntSlice[M, T] {
	return NewIntSlice[M, T](s.fn_("math::top", count.key()))
}
