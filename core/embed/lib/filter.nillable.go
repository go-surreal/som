package lib

// Nillable is a filter builder for values that can be nil.
// M is the type of the model the filter is for.
type Nillable[M any] struct {
	key Key[M]
}

// NewNillable creates a new nillable filter builder.
// M is the type of the model the filter is for.
func NewNillable[M any](key Key[M]) *Nillable[M] {
	return &Nillable[M]{key: key}
}

// Nil returns a filter that checks if the value is nil.
func (n *Nillable[M]) Nil(is bool) Filter[M] {
	if is {
		return n.key.op(OpExactlyEqual, nil)
	}

	return n.key.op(OpNotEqual, nil)
}
