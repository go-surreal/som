// Code generated by github.com/go-surreal/som, DO NOT EDIT.

package lib

// ByteSlice is a filter that can be used for byte slice fields.
// M is the type of the outgoing model for the filter statement.
type ByteSlice[M any] struct {
	*Base[M, []byte, *ByteSlice[M], *Slice[M, []byte, *ByteSlice[M]]]
}

// NewByteSlice creates a new slice filter.
func NewByteSlice[M any](key Key[M]) *ByteSlice[M] {
	return &ByteSlice[M]{
		Base: NewBase[M, []byte, *ByteSlice[M], *Slice[M, []byte, *ByteSlice[M]]](key),
	}
}

type ByteSlicePtr[M any] struct {
	*ByteSlice[M]
	*Nillable[M]
}

func NewByteSlicePtr[M any](key Key[M]) *ByteSlicePtr[M] {
	return &ByteSlicePtr[M]{
		ByteSlice: NewByteSlice[M](key),
		Nillable:  NewNillable[M](key),
	}
}

type BytePtrSlice[M any] struct {
	*Base[M, []*byte, *ByteSlice[M], *Slice[M, []*byte, *ByteSlice[M]]]
}

func NewBytePtrSlice[M any](key Key[M]) *BytePtrSlice[M] {
	return &BytePtrSlice[M]{
		Base: NewBase[M, []*byte, *ByteSlice[M], *Slice[M, []*byte, *ByteSlice[M]]](key),
	}
}

type BytePtrSlicePtr[M any] struct {
	*BytePtrSlice[M]
	*Nillable[M]
}

func NewBytePtrSlicePtr[M any](key Key[M]) *BytePtrSlicePtr[M] {
	return &BytePtrSlicePtr[M]{
		BytePtrSlice: NewBytePtrSlice[M](key),
		Nillable:     NewNillable[M](key),
	}
}

func (s *ByteSlice[M]) Base64Encode() *String[M] {
	return &String[M]{
		Base: NewBase[M, string, *String[M], *Slice[M, string, *String[M]]](s.fn("encoding::base64::encode")),
	}
}
