//go:build embed

package lib

type AnyNumeric[M any] interface {
	field[M]
	anyNumeric()
}

type Numeric[M, T any] struct {
	*Base[M, T, AnyNumeric[M], *Slice[M, T, *Numeric[M, T]]]
	*Comparable[M, T, AnyNumeric[M]]
}

func NewNumeric[M, T any](key Key[M]) *Numeric[M, T] {
	return &Numeric[M, T]{
		Base:       NewBase[M, T, AnyNumeric[M], *Slice[M, T, *Numeric[M, T]]](key),
		Comparable: NewComparable[M, T, AnyNumeric[M]](key),
	}
}

type NumericPtr[M, T any] struct {
	*Numeric[M, T]
	*Nillable[M]
}

func NewNumericPtr[M, T any](key Key[M]) *NumericPtr[M, T] {
	return &NumericPtr[M, T]{
		Numeric:  NewNumeric[M, T](key),
		Nillable: NewNillable[M](key),
	}
}

func (n *Numeric[M, T]) key() Key[M] {
	return n.Base.key()
}

func (n *Numeric[M, T]) anyNumeric() {}

//
// -- ARITHMETIC
//

func (n *Numeric[M, T]) Add(val float64) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc(OpAdd, val))
}

func (n *Numeric[M, T]) Sub(val float64) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc(OpSub, val))
}

func (n *Numeric[M, T]) Mul(val float64) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc(OpMul, val))
}

func (n *Numeric[M, T]) Div(val float64) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc(OpDiv, val))
}

func (n *Numeric[M, T]) Raise(val float64) *Float[M, float64] {
	return NewFloat[M, float64](n.Base.calc(OpRaise, val))
}

//
// -- MATH
//

//func (n *Numeric[M, M]) METHODS(min, max M) Filter[M] {
//	// https://surrealdb.com/docs/surrealdb/surrealql/functions/database/math
//}

func (n *Numeric[M, T]) Abs() *Numeric[M, T] {
	return NewNumeric[M, T](n.Base.fn("math::abs"))
}

// TODO: math::acos (surrealdb v2)

// TODO: math::acot (surrealdb v2)

// TODO: math::asin (surrealdb v2)

// TODO: math::atan (surrealdb v2)

// TODO: math::clamp (surrealdb v2)

// TODO: math::cos (surrealdb v2)

// TODO: math::cot (surrealdb v2)

// TODO: math::deg2rad (surrealdb v2)

// TODO: math::lerp (surrealdb v2)

// TODO: math::lerpangle (surrealdb v2)

// TODO: math::ln (surrealdb v2)

// TODO: math::log (surrealdb v2)

// TODO: math::log10 (surrealdb v2)

// TODO: math::log2 (surrealdb v2)

// TODO: math::rad2deg (surrealdb v2)

// TODO: math::sign (surrealdb v2)

// TODO: math::sin (surrealdb v2)

// TODO: math::tan (surrealdb v2)

func (n *Numeric[M, T]) Sqrt() *Float[M, float64] { // TODO: number must not be negative!
	return NewFloat[M, float64](n.Base.fn("math::sqrt"))
}

//
// -- CONVERT DURATION
//

func (n *Numeric[M, T]) AsDurationDays() *Duration[M] {
	return NewDuration[M](n.Base.fn("duration::from::days"))
}

func (n *Numeric[M, T]) AsDurationHours() *Duration[M] {
	return NewDuration[M](n.Base.fn("duration::from::hours"))
}

func (n *Numeric[M, T]) AsDurationMicros() *Duration[M] {
	return NewDuration[M](n.Base.fn("duration::from::micros"))
}

func (n *Numeric[M, T]) AsDurationMillis() *Duration[M] {
	return NewDuration[M](n.Base.fn("duration::from::millis"))
}

func (n *Numeric[M, T]) AsDurationMins() *Duration[M] {
	return NewDuration[M](n.Base.fn("duration::from::mins"))
}

func (n *Numeric[M, T]) AsDurationNanos() *Duration[M] {
	return NewDuration[M](n.Base.fn("duration::from::nanos"))
}

func (n *Numeric[M, T]) AsDurationSecs() *Duration[M] {
	return NewDuration[M](n.Base.fn("duration::from::secs"))
}

func (n *Numeric[M, T]) AsDurationWeeks() *Duration[M] {
	return NewDuration[M](n.Base.fn("duration::from::weeks"))
}

//
// -- CONVERT TIME
//

func (n *Numeric[M, T]) AsTimeMicros() *Time[M] {
	return NewTime[M](n.Base.fn("time::from::micros"))
}

func (n *Numeric[M, T]) AsTimeMillis() *Time[M] {
	return NewTime[M](n.Base.fn("time::from::millis"))
}

func (n *Numeric[M, T]) AsTimeNanos() *Time[M] {
	return NewTime[M](n.Base.fn("time::from::nanos"))
}

func (n *Numeric[M, T]) AsTimeSecs() *Time[M] {
	return NewTime[M](n.Base.fn("time::from::secs"))
}

func (n *Numeric[M, T]) AsTimeUnix() *Time[M] {
	return NewTime[M](n.Base.fn("time::from::unix"))
}
