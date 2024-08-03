package lib

type Numeric[M, T any] struct {
	*Base[M, T]
	*Comparable[M, T]
}

func NewNumeric[M, T any](key Key[M]) *Numeric[M, T] {
	return &Numeric[M, T]{
		Base:       NewBase[M, T](key),
		Comparable: NewComparable[M, T](key),
	}
}

//func (n *Numeric[M, M]) METHODS(min, max M) Filter[M] {
//	// https://surrealdb.com/docs/surrealdb/surrealql/functions/database/math
//}

type NumericPtr[M, T any] struct {
	*Numeric[M, T]
	*Nillable[M]
}

func NewNumericPtr[M, T any](key Key[M]) *NumericPtr[M, T] {
	return &NumericPtr[M, T]{
		Numeric:  NewNumeric[M, T](key),
		Nillable: NewNillable[M](key),
	}
}
