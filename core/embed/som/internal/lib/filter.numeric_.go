//go:build embed

package lib

func (n *Numeric[M, T]) Add_(val AnyNumeric[M]) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc_(OpAdd, val.key()))
}

func (n *Numeric[M, T]) Sub_(val AnyNumeric[M]) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc_(OpSub, val.key()))
}

func (n *Numeric[M, T]) Mul_(val AnyNumeric[M]) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc_(OpMul, val.key()))
}

func (n *Numeric[M, T]) Div_(val AnyNumeric[M]) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc_(OpDiv, val.key()))
}

func (n *Numeric[M, T]) Raise_(val AnyNumeric[M]) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc_(OpRaise, val.key()))
}
