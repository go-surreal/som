package filter

import (
	"github.com/marcbinz/som/lib/builder"
)

type ID[R any] struct {
	key  Key
	node string
}

func NewID[R any](key Key, node string) *ID[R] {
	return &ID[R]{key: key, node: node}
}

func (b *ID[R]) Equal(val string) Of[R] {
	val = b.node + ":" + val
	return build[R](b.key, builder.OpEqual, val, false)
}

func (b *ID[R]) NotEqual(val string) Of[R] {
	val = b.node + ":" + val
	return build[R](b.key, builder.OpNotEqual, val, false)
}

func (b *ID[R]) In(vals []string) Of[R] {
	var mapped []string
	for _, val := range vals {
		mapped = append(mapped, b.node+":"+val)
	}
	return build[R](b.key, builder.OpInside, mapped, false)
}

func (b *ID[R]) NotIn(vals []string) Of[R] {
	var mapped []string
	for _, val := range vals {
		mapped = append(mapped, b.node+":"+val)
	}
	return build[R](b.key, builder.OpNotInside, mapped, false)
}
