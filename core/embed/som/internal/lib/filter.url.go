//go:build embed

package lib

import (
	"net/url"
)

type URL[M any] struct {
	*Base[M, url.URL, *URL[M], *Slice[M, url.URL, *URL[M]]]
}

func NewURL[M any](key Key[M]) *URL[M] {
	conv := func(val url.URL) any {
		return val.String()
	}

	return &URL[M]{
		Base: NewBaseConv[M, url.URL, *URL[M], *Slice[M, url.URL, *URL[M]]](key, conv),
	}
}

func (u *URL[M]) Domain() *String[M] {
	return NewString(u.fn("parse::url::domain"))
}

func (u *URL[M]) Fragment() *String[M] {
	return NewString(u.fn("parse::url::fragment"))
}

func (u *URL[M]) Host() *String[M] {
	return NewString(u.fn("parse::url::host"))
}

func (u *URL[M]) Path() *String[M] {
	return NewString(u.fn("parse::url::path"))
}

func (u *URL[M]) Port() *Numeric[M, int] {
	return NewNumeric[M, int](u.fn("parse::url::port"))
}

func (u *URL[M]) Query() *String[M] {
	return NewString(u.fn("parse::url::query"))
}

// TODO: https://surrealdb.com/docs/surrealdb/surrealql/functions/database/http

type URLPtr[M any] struct {
	*URL[M]
	*Nillable[M]
}

func NewURLPtr[M any](key Key[M]) *URLPtr[M] {
	return &URLPtr[M]{
		URL:      NewURL[M](key),
		Nillable: NewNillable[M](key),
	}
}
