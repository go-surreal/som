//go:build embed

package lib

import (
	"strings"
)

type Filter[M any] interface {
	build(*context, M) string
}

type filter[T any] func(*context, T) string

//nolint:unused
func (f filter[T]) build(ctx *context, t T) string {
	return f(ctx, t)
}

func KeyFilter[M any](key Key[M]) Filter[M] {
	return filter[M](func(ctx *context, _ M) string {
		return key.render(ctx)
	})
}

//
// -- BASE
//

// Base is a filter with basic comparison operations.
// M is the type of the model this filter is for.
// M is the type of the field this filter is for.
type Base[M any, T any] struct {
	key  Key[M]
	conv func(T) any
}

func NewBase[M, T any](key Key[M]) *Base[M, T] {
	return &Base[M, T]{key: key}
}

func NewBaseConv[M, T any](key Key[M], conv func(T) any) *Base[M, T] {
	return &Base[M, T]{key: key, conv: conv}
}

func (b *Base[M, T]) Equal(val T) Filter[M] {
	if b.conv != nil {
		return b.key.op(OpEqual, b.conv(val))
	}

	return b.key.op(OpEqual, val)
}

func (b *Base[M, T]) NotEqual(val T) Filter[M] {
	if b.conv != nil {
		return b.key.op(OpNotEqual, b.conv(val))
	}

	return b.key.op(OpNotEqual, val)
}

func (b *Base[M, T]) In(vals []T) Filter[M] {
	if b.conv != nil {
		var mapped []any

		for _, val := range vals {
			mapped = append(mapped, b.conv(val))
		}

		return b.key.op(OpInside, mapped)
	}

	return b.key.op(OpInside, vals)
}

func (b *Base[M, T]) NotIn(vals []T) Filter[M] {
	if b.conv != nil {
		var mapped []any

		for _, val := range vals {
			mapped = append(mapped, b.conv(val))
		}

		return b.key.op(OpNotInside, mapped)
	}

	return b.key.op(OpNotInside, vals)
}

func (b *Base[M, T]) Zero(is bool) Filter[M] {
	op := OpExactlyEqual

	if !is {
		op = OpNotEqual
	}

	var zero T

	if b.conv != nil {
		return b.key.op(op, b.conv(zero))
	}

	return b.key.op(op, zero)
}

type BasePtr[M, T any] struct {
	*Base[M, T]
	*Nillable[M]
}

func NewBasePtr[M, T any](key Key[M]) *BasePtr[M, T] {
	return &BasePtr[M, T]{
		Base:     NewBase[M, T](key),
		Nillable: NewNillable[M](key),
	}
}

func NewBasePtrConv[M, T any](key Key[M], conv func(T) any) *BasePtr[M, T] {
	return &BasePtr[M, T]{
		Base:     NewBaseConv[M, T](key, conv),
		Nillable: NewNillable[M](key),
	}
}

//
// -- ALL | ANY
//

type All[M any] []Filter[M]

func (a All[M]) build(ctx *context, t M) string {
	if len(a) < 1 {
		return ""
	}

	var parts []string
	for _, filter := range a {
		if part := filter.build(ctx, t); part != "" {
			parts = append(parts, strings.TrimPrefix(part, ".")) // TODO: better place to trim?
		}
	}

	if len(parts) < 1 {
		return ""
	}

	return "(" + strings.Join(parts, " "+string(OpAnd)+" ") + ")"
}

type Any[M any] []Filter[M]

//nolint:unused
func (a Any[M]) build(ctx *context, t M) string {
	if len(a) < 1 {
		return ""
	}

	var parts []string
	for _, filter := range a {
		if part := filter.build(ctx, t); part != "" {
			parts = append(parts, strings.TrimPrefix(part, ".")) // TODO: better place to trim?
		}
	}

	if len(parts) < 1 {
		return ""
	}

	return "(" + strings.Join(parts, " "+string(OpOr)+" ") + ")"
}
