//go:build embed

package lib

import (
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

type CountKeyPart[T any] struct {
	key Key[T]
}

func (p CountKeyPart[T]) render(ctx *context) string {
	return "count(" + strings.TrimPrefix(p.key.render(ctx), ".") + ")"
}

type Key[T any] []KeyPart

func NewKey[T any]() Key[T] {
	return Key[T]{}
}

func (k Key[T]) render(ctx *context) string {
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

func (k Key[T]) Count() Key[T] {
	return Key[T]{
		CountKeyPart[T]{key: k},
	}
}

func (k Key[T]) Op(op Operator, val any) Filter[T] {
	return filter[T](
		func(ctx *context, _ T) string {
			statement := k.render(ctx) + " " + string(op) + " " + ctx.asVar(val)
			return strings.TrimPrefix(statement, ".")
		},
	)
}
