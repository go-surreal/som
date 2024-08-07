// Code generated by github.com/go-surreal/som, DO NOT EDIT.

package lib

// Comparable is a filter with comparison operations.
// M is the type of the model this filter is for.
// M is the type of the field this filter is for.
type Comparable[M any, T any] struct {
	key  Key[M]
	conv func(T) any
}

func NewComparable[M, T any](key Key[M]) *Comparable[M, T] {
	return NewComparableConv[M, T](key, nil)
}

func NewComparableConv[M, T any](key Key[M], conv func(T) any) *Comparable[M, T] {
	return &Comparable[M, T]{key: key, conv: conv}
}

func (c *Comparable[M, T]) LessThan(val T) Filter[M] {
	if c.conv != nil {
		return c.key.op(OpLessThan, c.conv(val))
	}

	return c.key.op(OpLessThan, val)
}

func (c *Comparable[M, T]) LessThanEqual(val T) Filter[M] {
	if c.conv != nil {
		return c.key.op(OpLessThanEqual, c.conv(val))
	}

	return c.key.op(OpLessThanEqual, val)
}

func (c *Comparable[M, T]) GreaterThan(val T) Filter[M] {
	if c.conv != nil {
		return c.key.op(OpGreaterThan, c.conv(val))
	}

	return c.key.op(OpGreaterThan, val)
}

func (c *Comparable[M, T]) GreaterThanEqual(val T) Filter[M] {
	if c.conv != nil {
		return c.key.op(OpGreaterThanEqual, c.conv(val))
	}

	return c.key.op(OpGreaterThanEqual, val)
}
