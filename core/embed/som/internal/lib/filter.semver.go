//go:build embed

package lib

import (
	"github.com/go-surreal/som/core/embed/som/som"
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

func (b *SemVer[M]) Compare(other som.SemVer) *Numeric[M, int] {
	return NewNumeric[M, int](b.fn("string::semver::compare", other))
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

func (b *SemVer[M]) IncMajor() *SemVer[M] {
	return NewSemVer[M](b.Base.fn("string::semver::inc::major"))
}

func (b *SemVer[M]) IncMinor() *SemVer[M] {
	return NewSemVer[M](b.Base.fn("string::semver::inc::minor"))
}

func (b *SemVer[M]) IncPatch() *SemVer[M] {
	return NewSemVer[M](b.Base.fn("string::semver::inc::patch"))
}

func (b *SemVer[M]) SetMajor(val int) *SemVer[M] {
	return NewSemVer[M](b.Base.fn("string::semver::set::major", val))
}

func (b *SemVer[M]) SetMinor(val int) *SemVer[M] {
	return NewSemVer[M](b.Base.fn("string::semver::set::minor", val))
}

func (b *SemVer[M]) SetPatch(val int) *SemVer[M] {
	return NewSemVer[M](b.Base.fn("string::semver::set::patch", val))
}
