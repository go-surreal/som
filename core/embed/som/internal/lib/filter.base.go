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
// E is the type of the field this filter is for.
type Base[M, T any, F, S field[M]] struct {
	Key[M]
	conv func(T) any
}

func NewBase[M, T any, F, S field[M]](key Key[M]) *Base[M, T, F, S] {
	return &Base[M, T, F, S]{Key: key}
}

func NewBaseConv[M, T any, F, S field[M]](key Key[M], conv func(T) any) *Base[M, T, F, S] {
	return &Base[M, T, F, S]{Key: key, conv: conv}
}

func (b *Base[M, T, F, S]) Equal(val T) Filter[M] {
	if b.conv != nil {
		return b.Key.op(OpEqual, b.conv(val))
	}

	return b.Key.op(OpEqual, val)
}

func (b *Base[M, T, F, S]) NotEqual(val T) Filter[M] {
	if b.conv != nil {
		return b.Key.op(OpNotEqual, b.conv(val))
	}

	return b.Key.op(OpNotEqual, val)
}

func (b *Base[M, T, F, S]) In(vals []T) Filter[M] {
	if b.conv != nil {
		mapped := make([]any, len(vals))

		for i, val := range vals {
			mapped[i] = b.conv(val)
		}

		return b.Key.op(OpIn, mapped)
	}

	return b.Key.op(OpIn, vals)
}

func (b *Base[M, T, F, S]) NotIn(vals []T) Filter[M] {
	if b.conv != nil {
		mapped := make([]any, len(vals))

		for i, val := range vals {
			mapped[i] = b.conv(val)
		}

		return b.Key.op(OpNotIn, mapped)
	}

	return b.Key.op(OpNotIn, vals)
}

func (b *Base[M, T, F, S]) Truth() *Bool[M] {
	return NewBool(b.Key.prefix(OpTruth))
}

// TODO: value::diff($value, $value) and value::patch($value, $diff)
// https://github.com/surrealdb/surrealdb/pull/4608

func (b *Base[M, T, F, S]) Zero(is bool) Filter[M] {
	op := OpExactlyEqual

	if !is {
		op = OpNotEqual
	}

	var zero T

	if b.conv != nil {
		return b.Key.op(op, b.conv(zero))
	}

	return b.Key.op(op, zero)
}

type BasePtr[M, T any, F, S field[M]] struct {
	*Base[M, T, F, S]
	*Nillable[M]
}

func NewBasePtr[M, T any, F, S field[M]](key Key[M]) *BasePtr[M, T, F, S] {
	return &BasePtr[M, T, F, S]{
		Base:     NewBase[M, T, F, S](key),
		Nillable: NewNillable[M](key),
	}
}

func NewBasePtrConv[M, T any, F, S field[M]](key Key[M], conv func(T) any) *BasePtr[M, T, F, S] {
	return &BasePtr[M, T, F, S]{
		Base:     NewBaseConv[M, T, F, S](key, conv),
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
