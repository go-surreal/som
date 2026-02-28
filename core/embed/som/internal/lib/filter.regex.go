//go:build embed

package lib

import (
	"regexp"
)

type Regex[M any] struct {
	*Base[M, regexp.Regexp, *Regex[M], *Slice[M, regexp.Regexp, *Regex[M]]]
}

func NewRegex[M any](key Key[M]) *Regex[M] {
	conv := func(val regexp.Regexp) any {
		return val.String()
	}

	return &Regex[M]{
		Base: NewBaseConv[M, regexp.Regexp, *Regex[M], *Slice[M, regexp.Regexp, *Regex[M]]](key, conv),
	}
}

type RegexPtr[M any] struct {
	*Regex[M]
	*Nillable[M]
}

func NewRegexPtr[M any](key Key[M]) *RegexPtr[M] {
	return &RegexPtr[M]{
		Regex:    NewRegex[M](key),
		Nillable: NewNillable[M](key),
	}
}
