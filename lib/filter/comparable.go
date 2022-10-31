package filter

import (
	"github.com/marcbinz/sdb/lib/builder"
)

type Comparable[T any, R any] struct {
	key     string
	isCount bool
}

func (c *Comparable[T, R]) LessThan(val T) Of[R] {
	return newOf[R](builder.OpLessThan, c.key, val, c.isCount)
}

func (c *Comparable[T, R]) LessThanEqual(val T) Of[R] {
	return newOf[R](builder.OpLessThanEqual, c.key, val, c.isCount)
}

func (c *Comparable[T, R]) GreaterThan(val T) Of[R] {
	return newOf[R](builder.OpGreaterThan, c.key, val, c.isCount)
}

func (c *Comparable[T, R]) GreaterThanEqual(val T) Of[R] {
	return newOf[R](builder.OpGreaterThanEqual, c.key, val, c.isCount)
}
