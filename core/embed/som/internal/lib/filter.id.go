//go:build embed

package lib

import "github.com/go-surreal/som"

type ID[M any] struct {
	key  Key[M]
	node string // TODO!
}

func NewID[M any](key Key[M], node string) *ID[M] {
	return &ID[M]{key: key, node: node}
}

func (b *ID[M]) Equal(val *som.ID) Filter[M] {
	// val = b.node + ":" + val
	return b.key.op(OpEqual, val)
}

func (b *ID[M]) NotEqual(val *som.ID) Filter[M] {
	// val = b.node + ":" + val
	return b.key.op(OpNotEqual, val)
}

func (b *ID[M]) In(vals []*som.ID) Filter[M] {
	return b.key.op(OpIn, vals)
}

func (b *ID[M]) NotIn(vals []*som.ID) Filter[M] {
	return b.key.op(OpNotIn, vals)
}
