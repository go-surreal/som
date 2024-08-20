//go:build embed

package lib

// Comparable is a filter with comparison operations.
// M is the type of the model this filter is for.
// T is the type of the field this filter is for.
type Comparable[M any, T any, F field[M]] struct {
	Key[M]
	conv func(T) any
}

func NewComparable[M, T any, F field[M]](key Key[M]) *Comparable[M, T, F] {
	return NewComparableConv[M, T, F](key, nil)
}

func NewComparableConv[M, T any, F field[M]](key Key[M], conv func(T) any) *Comparable[M, T, F] {
	return &Comparable[M, T, F]{Key: key, conv: conv}
}

func (c *Comparable[M, T, F]) LessThan(val T) Filter[M] {
	if c.conv != nil {
		return c.op(OpLessThan, c.conv(val))
	}

	return c.op(OpLessThan, val)
}

func (c *Comparable[M, T, F]) LessThanEqual(val T) Filter[M] {
	if c.conv != nil {
		return c.op(OpLessThanEqual, c.conv(val))
	}

	return c.op(OpLessThanEqual, val)
}

func (c *Comparable[M, T, F]) GreaterThan(val T) Filter[M] {
	if c.conv != nil {
		return c.op(OpGreaterThan, c.conv(val))
	}

	return c.op(OpGreaterThan, val)
}

func (c *Comparable[M, T, F]) GreaterThanEqual(val T) Filter[M] {
	if c.conv != nil {
		return c.op(OpGreaterThanEqual, c.conv(val))
	}

	return c.op(OpGreaterThanEqual, val)
}
