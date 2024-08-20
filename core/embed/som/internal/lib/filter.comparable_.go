//go:build embed

package lib

func (c *Comparable[M, T, F]) LessThan_(field F) Filter[M] {
	return c.op_(OpLessThan, field.key())
}

func (c *Comparable[M, T, F]) LessThanEqual_(field F) Filter[M] {
	return c.op_(OpLessThanEqual, field.key())
}

func (c *Comparable[M, T, F]) GreaterThan_(field F) Filter[M] {
	return c.op_(OpGreaterThan, field.key())
}

func (c *Comparable[M, T, F]) GreaterThanEqual_(field F) Filter[M] {
	return c.op_(OpGreaterThanEqual, field.key())
}
