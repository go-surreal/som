//go:build embed

package lib

func (f *Float[M, T]) Fixed_(places *Int[M, int]) *Float[M, T] {
	return NewFloat[M, T](f.Base.fn_("math::fixed", places.key()))
}
