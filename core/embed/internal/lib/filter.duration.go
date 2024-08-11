//go:build embed

package lib

import (
	"github.com/go-surreal/sdbc"
	"time"
)

type Duration[M any] struct {
	*Base[M, time.Duration]
	*Comparable[M, time.Duration]
}

func NewDuration[M any](key Key[M]) *Duration[M] {
	conv := func(val time.Duration) any {
		return sdbc.Duration{val}
	}

	return &Duration[M]{
		Base:       NewBaseConv[M, time.Duration](key, conv),
		Comparable: NewComparableConv[M, time.Duration](key, conv),
	}
}

func (d *Duration[M]) Days() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.key.fn("duration::days"))
}

func (d *Duration[M]) Hours() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.key.fn("duration::hours"))
}

func (d *Duration[M]) Micros() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.key.fn("duration::micros"))
}

func (d *Duration[M]) Millis() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.key.fn("duration::millis"))
}

func (d *Duration[M]) Mins() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.key.fn("duration::mins"))
}

func (d *Duration[M]) Nanos() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.key.fn("duration::nanos"))
}

func (d *Duration[M]) Secs() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.key.fn("duration::secs"))
}

func (d *Duration[M]) Weeks() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.key.fn("duration::weeks"))
}

func (d *Duration[M]) Years() *Numeric[M, int] {
	return NewNumeric[M, int](d.Base.key.fn("duration::years"))
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
