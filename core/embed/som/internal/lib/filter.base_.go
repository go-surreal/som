//go:build embed

package lib

type field[M any] interface {
	key() Key[M]
}

func (b *Base[M, T, F, S]) Equal_(field F) Filter[M] {
	return b.op_(OpEqual, field.key())
}

func (b *Base[M, T, F, S]) NotEqual_(field F) Filter[M] {
	return b.op_(OpNotEqual, field.key())
}

func (b *Base[M, T, F, S]) In_(field S) Filter[M] {
	return b.op_(OpIn, field.key())
}

func (b *Base[M, T, F, S]) NotIn_(field S) Filter[M] {
	return b.op_(OpNotIn, field.key())
}
