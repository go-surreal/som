package filter

type Numeric[T, R any] struct {
	*Base[T, R]
	*Comparable[T, R]
}

func NewNumeric[T, R any](key Key) *Numeric[T, R] {
	return newNumeric[T, R](key, false)
}

func newNumeric[T, R any](key Key, count bool) *Numeric[T, R] {
	return &Numeric[T, R]{
		Base:       &Base[T, R]{key: key, isCount: count},
		Comparable: &Comparable[T, R]{key: key, isCount: count},
	}
}

type NumericPtr[T, R any] struct {
	*Numeric[T, R]
	*Nillable[R]
}

func NewNumericPtr[T, R any](key Key) *NumericPtr[T, R] {
	return &NumericPtr[T, R]{
		Numeric:  NewNumeric[T, R](key),
		Nillable: &Nillable[R]{key: key},
	}
}
