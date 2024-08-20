//go:build embed

package lib

import (
	"github.com/go-surreal/sdbc"
	"time"
)

type Time[M any] struct {
	*Base[M, time.Time, *Time[M], *Slice[M, time.Time, *Time[M]]]
	comp *Comparable[M, time.Time, *Time[M]]
}

func NewTime[M any](key Key[M]) *Time[M] {
	conv := func(val time.Time) any {
		return sdbc.DateTime{val}
	}

	return &Time[M]{
		Base: NewBaseConv[M, time.Time, *Time[M], *Slice[M, time.Time, *Time[M]]](key, conv),
		comp: NewComparableConv[M, time.Time, *Time[M]](key, conv),
	}
}

func (t *Time[M]) Before(val time.Time) Filter[M] {
	return t.comp.LessThan(val)
}

func (t *Time[M]) BeforeOrEqual(val time.Time) Filter[M] {
	return t.comp.LessThanEqual(val)
}

func (t *Time[M]) After(val time.Time) Filter[M] {
	return t.comp.GreaterThan(val)
}

func (t *Time[M]) AfterOrEqual(val time.Time) Filter[M] {
	return t.comp.GreaterThanEqual(val)
}

func (t *Time[M]) Day() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::day"))
}

func (t *Time[M]) Floor(dur time.Duration) *Time[M] {
	return NewTime[M](t.fn("time::floor", sdbc.Duration{dur}))
}

func (t *Time[M]) Format(format string) *String[M] {
	return NewString[M](t.fn("time::format", format))
}

type Group string // TODO!

const (
	GroupYear   Group = "year"
	GroupMonth  Group = "month"
	GroupDay    Group = "day"
	GroupHour   Group = "hour"
	GroupMinute Group = "minute"
	GroupSecond Group = "second"
)

func (t *Time[M]) Group(group Group) *Time[M] {
	return NewTime[M](t.fn("time::group", group))
}

func (t *Time[M]) Hour() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::hour"))
}

func (t *Time[M]) Micros() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::micros"))
}

func (t *Time[M]) Millis() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::millis"))
}

func (t *Time[M]) Minute() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::minute"))
}

func (t *Time[M]) Month() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::month"))
}

func (t *Time[M]) Nano() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::nano")) // TODO: int64? big.Int?
}

func (t *Time[M]) Round(dur time.Duration) *Time[M] {
	return NewTime[M](t.fn("time::round", sdbc.Duration{dur}))
}

func (t *Time[M]) Second() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::second"))
}

func (t *Time[M]) Timezone() *String[M] {
	return NewString[M](t.fn("time::timezone"))
}

func (t *Time[M]) Unix() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::unix"))
}

func (t *Time[M]) Weekday() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::wday")) // TODO: time.Weekday type!
}

func (t *Time[M]) Week() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::week"))
}

func (t *Time[M]) YearDay() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::yday"))
}

func (t *Time[M]) Year() *Numeric[M, int] {
	return NewNumeric[M, int](t.fn("time::year"))
}

type TimePtr[R any] struct {
	*Time[R]
	*Nillable[R]
}

func NewTimePtr[M any](key Key[M]) *TimePtr[M] {
	return &TimePtr[M]{
		Time:     NewTime[M](key),
		Nillable: NewNillable[M](key),
	}
}
