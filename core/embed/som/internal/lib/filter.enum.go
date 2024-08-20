//go:build embed

package lib

type Enum[M, E any] struct {
	*Base[M, E, *Enum[M, E], *Slice[M, E, *Enum[M, E]]]
}

func NewEnum[M, E any](key Key[M]) *Enum[M, E] {
	return &Enum[M, E]{
		Base: NewBase[M, E, *Enum[M, E], *Slice[M, E, *Enum[M, E]]](key),
	}
}

type EnumPtr[M, E any] struct {
	*Enum[M, E]
	*Nillable[M]
}

func NewEnumPtr[M, E any](key Key[M]) *EnumPtr[M, E] {
	return &EnumPtr[M, E]{
		Enum:     NewEnum[M, E](key),
		Nillable: NewNillable[M](key),
	}
}
