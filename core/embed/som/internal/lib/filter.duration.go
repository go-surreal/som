//go:build embed

package lib

import (
	"github.com/go-surreal/sdbc"
	"time"
)

type Duration[M any] struct {
	*Base[M, time.Duration, *Duration[M], *Slice[M, time.Duration, *Duration[M]]]
	*Comparable[M, time.Duration, *Duration[M]]
}

func NewDuration[M any](key Key[M]) *Duration[M] {
	conv := func(val time.Duration) any {
		return sdbc.Duration{val}
	}

	return &Duration[M]{
		Base:       NewBaseConv[M, time.Duration, *Duration[M], *Slice[M, time.Duration, *Duration[M]]](key, conv),
		Comparable: NewComparableConv[M, time.Duration, *Duration[M]](key, conv),
	}
}

func (d *Duration[M]) key() Key[M] {
	return d.Base.key()
}

func (d *Duration[M]) Days() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.fn("duration::days"))
}

func (d *Duration[M]) Hours() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.fn("duration::hours"))
}

func (d *Duration[M]) Micros() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.fn("duration::micros"))
}

func (d *Duration[M]) Millis() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.fn("duration::millis"))
}

func (d *Duration[M]) Mins() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.fn("duration::mins"))
}

func (d *Duration[M]) Nanos() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.fn("duration::nanos"))
}

func (d *Duration[M]) Secs() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.fn("duration::secs"))
}

func (d *Duration[M]) Weeks() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.fn("duration::weeks"))
}

func (d *Duration[M]) Years() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.fn("duration::years"))
}

type DurationPtr[M any] struct {
	*Duration[M]
	*Nillable[M]
}

func NewDurationPtr[M any](key Key[M]) *DurationPtr[M] {
	return &DurationPtr[M]{
		Duration: NewDuration[M](key),
		Nillable: NewNillable[M](key),
	}
}
