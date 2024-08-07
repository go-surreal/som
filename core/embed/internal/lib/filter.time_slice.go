//go:build embed

package lib

import "time"

type TimeSlice[M any] struct {
	*Slice[M, time.Time, *Time[M]]
}

func NewTimeSlice[M any](key Key[M]) *TimeSlice[M] {
	return &TimeSlice[M]{
		Slice: NewSlice[M, time.Time, *Time[M]](key, NewTime[M]),
	}
}

func (t *TimeSlice[M]) Min() *Time[M] {
	return NewTime[M](t.key.fn("time::min"))
}

func (t *TimeSlice[M]) Max() *Time[M] {
	return NewTime[M](t.key.fn("time::max"))
}
