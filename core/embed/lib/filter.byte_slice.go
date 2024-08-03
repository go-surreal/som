package lib

// ByteSlice is a filter that can be used for byte slice fields.
// M is the type of the outgoing model for the filter statement.
type ByteSlice[M any] struct {
	*Base[M, []byte]
}

// NewByteSlice creates a new slice filter.
func NewByteSlice[M any](key Key[M]) *ByteSlice[M] {
	return &ByteSlice[M]{
		Base: NewBase[M, []byte](key),
	}
}

func (s *ByteSlice[M]) Base64Encode() *String[M] {
	return &String[M]{
		Base: NewBase[M, string](s.Base.key.fn("encoding::base64::encode")),
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
