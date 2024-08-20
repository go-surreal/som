//go:build embed

package lib

import (
	"github.com/go-surreal/som"
)

// SemVer is a filter builder for string values.
// M is the model this filter is for.
type SemVer[M any] struct {
	*Base[M, som.SemVer, *SemVer[M], *Slice[M, som.SemVer, *SemVer[M]]]
}

func NewSemVer[M any](key Key[M]) *SemVer[M] {
	return &SemVer[M]{
		Base: NewBase[M, som.SemVer, *SemVer[M], *Slice[M, som.SemVer, *SemVer[M]]](key),
	}
}

func (b *SemVer[M]) Compare(val som.SemVer) *Numeric[M, int] {
	return NewNumeric[M, int](b.fn("string::semver::compare", val))
}

func (b *SemVer[M]) Major() *Numeric[M, int] {
	return NewNumeric[M, int](b.Base.fn("string::semver::major"))
}

func (b *SemVer[M]) Minor() *Numeric[M, int] {
	return NewNumeric[M, int](b.Base.fn("string::semver::minor"))
}

func (b *SemVer[M]) Patch() *Numeric[M, int] {
	return NewNumeric[M, int](b.Base.fn("string::semver::patch"))
}
