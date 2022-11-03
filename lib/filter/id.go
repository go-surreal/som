package filter

import (
	"github.com/marcbinz/sdb/lib/builder"
)

type ID[R any] struct {
	key  string
	node string
}

func NewID[R any](key string, node string) *ID[R] {
	return &ID[R]{key: key, node: node}
}

func (b *ID[R]) Equal(val string) Of[R] {
	val = b.node + ":" + val
	return newOf[R](builder.OpEqual, b.key, val, false)
}

func (b *ID[R]) NotEqual(val string) Of[R] {
	val = b.node + ":" + val
	return newOf[R](builder.OpEqual, b.key, val, false)
}

func (b *ID[R]) In(vals []string) Of[R] {
	var mapped []string
	for _, val := range vals {
		mapped = append(mapped, b.node+":"+val)
	}
	return newOf[R](builder.OpInside, b.key, mapped, false)
}

func (b *ID[R]) NotIn(vals []string) Of[R] {
	var mapped []string
	for _, val := range vals {
		mapped = append(mapped, b.node+":"+val)
	}
	return newOf[R](builder.OpNotInside, b.key, mapped, false)
}
