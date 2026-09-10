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
	mapped := any(val)

	if b.conv != nil {
		mapped = b.conv(val)
	}

	return b.notSet(OpNotEqual, mapped)
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
	mapped := any(vals)

	if b.conv != nil {
		m := make([]any, len(vals))

		for i, val := range vals {
			m[i] = b.conv(val)
		}

		mapped = m
	}

	return b.notSet(OpNotIn, mapped)
}

// notSet builds a negative comparison that also excludes unset (NONE/NULL)
// records. In SurrealDB "NONE != value" and "NONE ∉ set" both evaluate to
// true, so a bare negation would wrongly match records where the field is
// not set. For required fields the extra guard is always true and thus a
// no-op, so this stays correct for both optional and required fields.
func (b *Base[M, T, F, S]) notSet(op Operator, val any) Filter[M] {
	return filter[M](func(ctx *context, _ M) string {
		field := strings.TrimPrefix(b.Key.render(ctx), ".")
		return "(" + field + " " + string(op) + " " + ctx.asVar(val) +
			" AND " + field + " != NONE AND " + field + " != NULL)"
	})
}

func (b *Base[M, T, F, S]) Truth() *Bool[M] {
	return NewBool(b.Key.prefix(OpTruth))
}

// TODO: value::diff($value, $value) and value::patch($value, $diff)
// https://github.com/surrealdb/surrealdb/pull/4608

// Zero compares against the Go zero value of the field type. For a pointer
// field this is the element zero value (e.g. 0 or ""), NOT NONE, so Zero does
// not detect an unset field. Use Nil to check for NONE/NULL instead.
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

type Not[M any] struct {
	Filter Filter[M]
}

//nolint:unused
func (n Not[M]) build(ctx *context, t M) string {
	inner := n.Filter.build(ctx, t)
	if inner == "" {
		return ""
	}

	if !strings.HasPrefix(inner, "(") {
		return "!(" + inner + ")"
	}

	return "!" + inner
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
