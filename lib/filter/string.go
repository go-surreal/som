package filter

import (
	"github.com/marcbinz/sdb/lib/builder"
)

type String[R any] struct {
	*Base[string, R]
	*Comparable[string, R]
}

func NewString[R any](key string) *String[R] {
	return &String[R]{
		Base:       &Base[string, R]{key: key},
		Comparable: &Comparable[string, R]{key: key},
	}
}

func (s *String[R]) FuzzyMatch(val string) Of[R] {
	return newOf[R](builder.OpFuzzyMatch, s.Base.key, val, false)
}

func (s *String[R]) NotFuzzyMatch(val string) Of[R] {
	return newOf[R](builder.OpFuzzyNotMatch, s.Base.key, val, false)
}
