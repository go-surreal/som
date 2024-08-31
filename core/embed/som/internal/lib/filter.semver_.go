//go:build embed

package lib

func (b *SemVer[M]) Compare_(other *SemVer[M]) *Numeric[M, int] {
	return NewNumeric[M, int](b.fn_("string::semver::compare", other.key()))
}

func (b *SemVer[M]) SetMajor_(val AnyInt[M]) *SemVer[M] {
	return NewSemVer[M](b.Base.fn_("string::semver::set::major", val.key()))
}

func (b *SemVer[M]) SetMinor_(val AnyInt[M]) *SemVer[M] {
	return NewSemVer[M](b.Base.fn_("string::semver::set::minor", val.key()))
}

func (b *SemVer[M]) SetPatch_(val AnyInt[M]) *SemVer[M] {
	return NewSemVer[M](b.Base.fn_("string::semver::set::patch", val.key()))
}
