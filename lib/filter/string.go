package filter

import (
	"github.com/marcbinz/som/lib/builder"
)

type String[R any] struct {
	*Base[string, R]
	*Comparable[string, R]
}

func NewString[R any](key Key) *String[R] {
	return &String[R]{
		Base:       &Base[string, R]{key: key},
		Comparable: &Comparable[string, R]{key: key},
	}
}

func (s *String[R]) FuzzyMatch(val string) Of[R] {
	return build[R](s.Base.key, builder.OpFuzzyMatch, val, false)
}

func (s *String[R]) NotFuzzyMatch(val string) Of[R] {
	return build[R](s.Base.key, builder.OpFuzzyNotMatch, val, false)
}

type StringPtr[R any] struct {
	*String[R]
	*Nillable[R]
}

func NewStringPtr[R any](key Key) *StringPtr[R] {
	return &StringPtr[R]{
		String:   NewString[R](key),
		Nillable: &Nillable[R]{key: key},
	}
}
