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

func (d *Numeric[M, T]) key() Key[M] {
	return d.Base.key()
}

//func (n *Numeric[M, M]) METHODS(min, max M) Filter[M] {
//	// https://surrealdb.com/docs/surrealdb/surrealql/functions/database/math
//}

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

func (n *Numeric[M, T]) ToInt() *Numeric[M, int] {
	return NewNumeric[M, int](n.Base.prefix(CastInt))
}

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
