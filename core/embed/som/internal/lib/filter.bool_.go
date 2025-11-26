//go:build embed

package lib

func (b *Bool[M]) Is_(field *Bool[M]) Filter[M] {
	return b.op_(OpExactlyEqual, field.Key)
}
