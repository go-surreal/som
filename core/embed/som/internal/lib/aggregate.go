//go:build embed

package lib

import "strings"

// KeyProvider is implemented by every filter field. It exposes the underlying
// field key so the aggregate helpers and the view projection renderer can
// build static SurrealQL expressions from a filter field reference.
//
// It is used only when defining read-only views; aggregate functions are
// intentionally kept out of the filter field method set so they cannot be
// used in a WHERE clause, where aggregates are invalid.
type KeyProvider[M any] interface {
	field[M]
}

// Count renders count(<field>): the number of non-null values in the group.
func Count[M any](f KeyProvider[M]) *Int[M, int] {
	return NewInt[M, int](Fn(f.key(), "count"))
}

// Sum renders math::sum(<field>).
func Sum[M any](f KeyProvider[M]) *Float[M, float64] {
	return NewFloat[M, float64](Fn(f.key(), "math::sum"))
}

// Mean renders math::mean(<field>).
func Mean[M any](f KeyProvider[M]) *Float[M, float64] {
	return NewFloat[M, float64](Fn(f.key(), "math::mean"))
}

// Min renders math::min(<field>).
func Min[M any](f KeyProvider[M]) *Float[M, float64] {
	return NewFloat[M, float64](Fn(f.key(), "math::min"))
}

// Max renders math::max(<field>).
func Max[M any](f KeyProvider[M]) *Float[M, float64] {
	return NewFloat[M, float64](Fn(f.key(), "math::max"))
}

// Variance renders math::variance(<field>).
func Variance[M any](f KeyProvider[M]) *Float[M, float64] {
	return NewFloat[M, float64](Fn(f.key(), "math::variance"))
}

// StdDev renders math::stddev(<field>).
func StdDev[M any](f KeyProvider[M]) *Float[M, float64] {
	return NewFloat[M, float64](Fn(f.key(), "math::stddev"))
}

// RenderProjection renders a field reference or aggregate as a static
// SurrealQL expression (values inlined) for use in a view projection or
// GROUP BY clause.
func RenderProjection[M any](f KeyProvider[M]) string {
	ctx := &context{vars: map[string]any{}, literal: true}
	return strings.TrimPrefix(f.key().render(ctx), ".")
}
