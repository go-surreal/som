package lib

import (
	"strings"
)

type KeyPart interface {
	render(ctx *context) string
}

type BaseKeyPart struct {
	name      string
	separator string
	filters   []Filter[any]
}

func (p BaseKeyPart) render(ctx *context) string {
	if len(p.filters) < 1 {
		return p.separator + p.name
	}

	return p.separator + "(" + p.name + " WHERE " + All(ToWhere(p.filters))(ctx) + ")"
}

type CountKeyPart struct {
	key Key
}

func (p CountKeyPart) render(ctx *context) string {
	return "count(" + strings.TrimPrefix(p.key.render(ctx), ".") + ")"
}

type Key []KeyPart

func NewKey() Key {
	return Key{}
}

func (k Key) render(ctx *context) string {
	var statement string

	for _, part := range k {
		statement += part.render(ctx)
	}

	return statement
}

func (k Key) Field(name string) Key {
	return append(k, BaseKeyPart{
		name:      name,
		separator: ".",
	})
}

func (k Key) Node(name string, filters []Filter[any]) Key {
	return append(k, BaseKeyPart{
		name:      name,
		separator: ".",
		filters:   filters,
	})
}

func (k Key) EdgeIn(name string, filters []Filter[any]) Key {
	return append(k, BaseKeyPart{
		name:      name,
		separator: "->",
		filters:   filters,
	})
}

func (k Key) EdgeOut(name string, filters []Filter[any]) Key {
	return append(k, BaseKeyPart{
		name:      name,
		separator: "<-",
		filters:   filters,
	})
}

func (k Key) Count() Key {
	return Key{
		CountKeyPart{key: k},
	}
}

func (k Key) Op(op Operator, val any) Where {
	return func(ctx *context) string {
		statement := k.render(ctx) + " " + string(op) + " " + ctx.asVar(val)
		return strings.TrimPrefix(statement, ".")
	}
}
