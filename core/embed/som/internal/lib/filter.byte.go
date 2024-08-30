//go:build embed

package lib

type Byte[M any] struct {
	*Base[M, byte, *Byte[M], *Slice[M, byte, *Byte[M]]]
}

func NewByte[M any](key Key[M]) *Byte[M] {
	return &Byte[M]{
		Base: NewBase[M, byte, *Byte[M], *Slice[M, byte, *Byte[M]]](key),
	}
}

type BytePtr[M any] struct {
	*Byte[M]
	*Nillable[M]
}

func NewBytePtr[M any](key Key[M]) *BytePtr[M] {
	return &BytePtr[M]{
		Byte:     NewByte[M](key),
		Nillable: NewNillable[M](key),
	}
}
