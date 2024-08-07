//go:build embed

package lib

type StringSlice[M any] struct {
	*Slice[M, string, *String[M]]
}

func NewStringSlice[M any](key Key[M]) *StringSlice[M] {
	return &StringSlice[M]{
		Slice: NewSlice[M, string, *String[M]](key, NewString[M]),
	}
}

func (s *StringSlice[M]) Join(sep string) *String[M] {
	return NewString(s.key.fn("string::join", sep))
}
