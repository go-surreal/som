//go:build embed

package lib

type Numeric[M, T any] struct {
	*Base[M, T, *Numeric[M, T], *Slice[M, T, *Numeric[M, T]]]
	*Comparable[M, T, *Numeric[M, T]]
}

func NewNumeric[M, T any](key Key[M]) *Numeric[M, T] {
	return &Numeric[M, T]{
		Base:       NewBase[M, T, *Numeric[M, T], *Slice[M, T, *Numeric[M, T]]](key),
		Comparable: NewComparable[M, T, *Numeric[M, T]](key),
	}
}

func (n *Numeric[M, T]) key() Key[M] {
	return n.Base.key()
}

//
// -- ARITHMETIC
//

func (n *Numeric[M, T]) Add(val T) *Numeric[M, T] {
	return NewNumeric[M, T](n.Base.calc(OpAdd, val))
}

func (n *Numeric[M, T]) Sub(val T) *Numeric[M, T] {
	return NewNumeric[M, T](n.Base.calc(OpSub, val))
}

func (n *Numeric[M, T]) Mul(val T) *Numeric[M, T] {
	return NewNumeric[M, T](n.Base.calc(OpMul, val))
}

func (n *Numeric[M, T]) Div(val T) *Numeric[M, T] {
	return NewNumeric[M, T](n.Base.calc(OpDiv, val))
}

func (n *Numeric[M, T]) Raise(val float64) *Numeric[M, T] {
	return NewNumeric[M, T](n.Base.calc(OpRaise, val))
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

// TODO: math::acos (v2.0.0)

// TODO: math::acot (v2.0.0)

// TODO: math::asin (v2.0.0)

// TODO: math::atan (v2.0.0)

// TODO: math::clamp (v2.0.0)

// TODO: math::cos (v2.0.0)

// TODO: math::cot (v2.0.0)

// TODO: math::deg2rad (v2.0.0)

// TODO: math::e (??)

// TODO: math::inf (??)

// math::lerp
// math::lerpangle
// math::ln
// math::log
// math::log10
// math::log2
// math::neg_inf (??)
// math::pi (??)
// math::rad2deg
// math::sign
// math::sin
// math::tan
// math::tau (??)

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

// TODO: https://surrealdb.com/docs/surrealdb/surrealql/datamodel/numbers#mathematical-constants

//
// -- POINTER
//

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
