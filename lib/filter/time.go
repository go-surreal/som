package filter

import (
	"time"
)

type Time[R any] struct {
	*Base[time.Time, R]
	comp *Comparable[time.Time, R]
}

func NewTime[R any](key Key) *Time[R] {
	return &Time[R]{
		Base: &Base[time.Time, R]{key: key},
		comp: &Comparable[time.Time, R]{key: key},
	}
}

func (t *Time[R]) Before(val time.Time) Of[R] {
	return t.comp.LessThan(val)
}

func (t *Time[R]) BeforeOrEqual(val time.Time) Of[R] {
	return t.comp.LessThanEqual(val)
}

func (t *Time[R]) After(val time.Time) Of[R] {
	return t.comp.GreaterThan(val)
}

func (t *Time[R]) AfterOrEqual(val time.Time) Of[R] {
	return t.comp.GreaterThanEqual(val)
}

type TimePtr[R any] struct {
	*Time[R]
	*Nillable[R]
}

func NewTimePtr[R any](key Key) *TimePtr[R] {
	return &TimePtr[R]{
		Time:     NewTime[R](key),
		Nillable: &Nillable[R]{key: key},
	}
}
