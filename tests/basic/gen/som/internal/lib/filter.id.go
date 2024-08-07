// Code generated by github.com/go-surreal/som, DO NOT EDIT.

package lib

import "github.com/go-surreal/sdbc"

type ID[M any] struct {
	key  Key[M]
	node string // TODO!
}

func NewID[M any](key Key[M], node string) *ID[M] {
	return &ID[M]{key: key, node: node}
}

func (b *ID[M]) Equal(val *sdbc.ID) Filter[M] {
	// val = b.node + ":" + val
	return b.key.op(OpEqual, val)
}

func (b *ID[M]) NotEqual(val *sdbc.ID) Filter[M] {
	// val = b.node + ":" + val
	return b.key.op(OpNotEqual, val)
}

func (b *ID[M]) In(vals []*sdbc.ID) Filter[M] {
	return b.key.op(OpInside, vals)
}

func (b *ID[M]) NotIn(vals []*sdbc.ID) Filter[M] {
	return b.key.op(OpNotInside, vals)
}
