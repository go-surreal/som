package lib

import "net/url"

type URL[M any] struct {
	key Key[M]
}

func NewURL[M any](key Key[M]) *URL[M] {
	return &URL[M]{key: key}
}

func (u *URL[M]) Equal(val url.URL) Filter[M] {
	return u.key.op(OpEqual, val.String())
}

func (u *URL[M]) NotEqual(val url.URL) Filter[M] {
	return u.key.op(OpNotEqual, val.String())
}

func (u *URL[M]) In(vals []url.URL) Filter[M] {
	var mapped []string

	for _, val := range vals {
		mapped = append(mapped, val.String())
	}

	return u.key.op(OpInside, mapped)
}

func (u *URL[M]) NotIn(vals []url.URL) Filter[M] {
	var mapped []string

	for _, val := range vals {
		mapped = append(mapped, val.String())
	}

	return u.key.op(OpNotInside, mapped)
}

func (u *URL[M]) Domain() *String[M] {
	return NewString(u.key.fn("parse::url::domain"))
}

func (u *URL[M]) Fragment() *String[M] {
	return NewString(u.key.fn("parse::url::fragment"))
}

func (u *URL[M]) Host() *String[M] {
	return NewString(u.key.fn("parse::url::host"))
}

func (u *URL[M]) Path() *String[M] {
	return NewString(u.key.fn("parse::url::path"))
}

func (u *URL[M]) Port() *Numeric[M, int] {
	return NewNumeric[M, int](u.key.fn("parse::url::port"))
}

func (u *URL[M]) Query() *String[M] {
	return NewString(u.key.fn("parse::url::query"))
}

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
