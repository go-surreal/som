// Code generated by github.com/go-surreal/som, DO NOT EDIT.

package lib

func (b *SemVer[M]) Compare_(field *SemVer[M]) *Numeric[M, int] {
	return NewNumeric[M, int](b.fn_("string::semver::compare", field.key()))
}
