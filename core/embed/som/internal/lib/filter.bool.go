//go:build embed

package lib

// Bool is a filter builder for boolean fields.
// M is the type of the field the filter is for.
type Bool[M any] struct {
	Key[M]
}

func NewBool[M any](key Key[M]) *Bool[M] {
	return &Bool[M]{Key: key}
}

func (b *Bool[M]) Is(val bool) Filter[M] {
	return b.op(OpExactlyEqual, val)
}

func (b *Bool[M]) Not() *Bool[M] {
	return NewBool(b.prefix(OpNot))
}

type BoolPtr[M any] struct {
	*Bool[M]
	*Nillable[M]
}

func NewBoolPtr[M any](key Key[M]) *BoolPtr[M] {
	return &BoolPtr[M]{
		Bool:     NewBool[M](key),
		Nillable: NewNillable[M](key),
	}
}
