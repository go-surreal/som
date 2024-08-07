//go:build embed

package lib

import (
	"github.com/go-surreal/som"
)

// Email is a filter builder for email values.
// M is the model this filter is for.
type Email[M any] struct {
	*Base[M, som.Email]
}

func NewEmail[M any](key Key[M]) *Email[M] {
	return &Email[M]{
		Base: NewBase[M, som.Email](key),
	}
}

func (e *Email[M]) User() *String[M] {
	return NewString[M](e.key.fn("parse::email::user"))
}

func (e *Email[M]) Host() *String[M] {
	return NewString[M](e.key.fn("parse::email::host"))
}
