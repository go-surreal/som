package filter

import (
	"github.com/marcbinz/som/lib/builder"
)

type Comparable[T any, R any] struct {
	key     Key
	isCount bool
}

func (c *Comparable[T, R]) LessThan(val T) Of[R] {
	return build[R](c.key, builder.OpLessThan, val, c.isCount)
}

func (c *Comparable[T, R]) LessThanEqual(val T) Of[R] {
	return build[R](c.key, builder.OpLessThanEqual, val, c.isCount)
}

func (c *Comparable[T, R]) GreaterThan(val T) Of[R] {
	return build[R](c.key, builder.OpGreaterThan, val, c.isCount)
}

func (c *Comparable[T, R]) GreaterThanEqual(val T) Of[R] {
	return build[R](c.key, builder.OpGreaterThanEqual, val, c.isCount)
}
