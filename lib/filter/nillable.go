package filter

import (
	"github.com/marcbinz/som/lib/builder"
)

type Nillable[R any] struct {
	key Key
}

func (n *Nillable[R]) Nil() Of[R] {
	return build[R](n.key, builder.OpExactlyEqual, nil, false)
}

func (n *Nillable[R]) NotNil() Of[R] {
	return build[R](n.key, builder.OpNotEqual, nil, false)
}
