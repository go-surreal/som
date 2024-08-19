//go:build embed

package lib

import (
	"slices"
	"strings"
)

type KeyPart interface {
	render(ctx *context) string
}

type BaseKeyPart[T any] struct {
	name      string
	separator string
	filters   []Filter[T]
}

func (p BaseKeyPart[T]) render(ctx *context) string {
	where := ""

	if len(p.filters) > 0 {
		var t T
		where = "[WHERE " + All[T](p.filters).build(ctx, t) + "]"
	}

	return p.separator + p.name + where
}

type FuncKeyPart[T any] struct {
	key    Key[T]
	fn     string
	params []any
}

func (p FuncKeyPart[T]) render(ctx *context) string {
	var mappedParams []string

	for _, param := range p.params {
		mappedParams = append(mappedParams, ctx.asVar(param))
	}

	paramString := strings.Join(mappedParams, ",")

	if paramString != "" {
		paramString = ", " + paramString
	}

	return p.fn + "(" + strings.TrimPrefix(p.key.render(ctx), ".") + paramString + ")"
}

type Key[M any] []KeyPart

func NewKey[T any]() Key[T] {
	return Key[T]{}
}

func (k Key[M]) fn(fn string, params ...any) Key[M] {
	return Key[M]{
		FuncKeyPart[M]{
			key:    slices.Clone(k),
			fn:     fn,
			params: params,
		},
	}
}

func (k Key[M]) op(op Operator, val any) Filter[M] {
	return filter[M](
		func(ctx *context, _ M) string {
			return strings.TrimPrefix(k.render(ctx), ".") +
				" " +
				string(op) +
				" " +
				ctx.asVar(val)
		},
	)
}

func (k Key[M]) op_(op Operator, k2 Key[M]) Filter[M] {
	return filter[M](
		func(ctx *context, _ M) string {
			return strings.TrimPrefix(k.render(ctx), ".") +
				" " +
				string(op) +
				" " +
				strings.TrimPrefix(k2.render(ctx), ".")
		},
	)
}

func (k Key[M]) render(ctx *context) string {
	var statement string

	for _, part := range k {
		statement += part.render(ctx)
	}

	return statement
}

func Field[T any](k Key[T], name string) Key[T] {
	return append(k, BaseKeyPart[T]{
		name:      name,
		separator: ".",
	})
}

func Node[T, S any](k Key[T], name string, filters []Filter[S]) Key[T] {
	return append(k, BaseKeyPart[S]{
		name:      name,
		separator: ".",
		filters:   filters,
	})
}

func EdgeIn[T, S any](k Key[T], name string, filters []Filter[S]) Key[T] {
	return append(k, BaseKeyPart[S]{
		name:      name,
		separator: "->",
		filters:   filters,
	})
}

func EdgeOut[T, S any](key Key[T], name string, filters []Filter[S]) Key[T] {
	return append(key, BaseKeyPart[S]{
		name:      name,
		separator: "<-",
		filters:   filters,
	})
}
