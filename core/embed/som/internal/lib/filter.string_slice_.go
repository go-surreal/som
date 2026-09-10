//go:build embed

package lib

func (s *StringSlice[M]) Join_(field *String[M]) *String[M] {
	return NewString(s.fn_("string::join", field.key()))
}
