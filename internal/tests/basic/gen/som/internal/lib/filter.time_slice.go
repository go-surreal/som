// Code generated by github.com/go-surreal/som, DO NOT EDIT.

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

type TimeSlicePtr[M any] struct {
	*SlicePtr[M, time.Time, *Time[M]]
}

func NewTimeSlicePtr[M any](key Key[M]) *TimeSlicePtr[M] {
	return &TimeSlicePtr[M]{
		SlicePtr: NewSlicePtr[M, time.Time, *Time[M]](key, NewTime[M]),
	}
}

type TimePtrSlice[M any] struct {
	*Slice[M, *time.Time, *TimePtr[M]]
}

func NewTimePtrSlice[M any](key Key[M]) *TimePtrSlice[M] {
	return &TimePtrSlice[M]{
		Slice: NewSlice[M, *time.Time, *TimePtr[M]](key, NewTimePtr[M]),
	}
}

type TimePtrSlicePtr[M any] struct {
	*SlicePtr[M, *time.Time, *TimePtr[M]]
}

func NewTimePtrSlicePtr[M any](key Key[M]) *TimePtrSlicePtr[M] {
	return &TimePtrSlicePtr[M]{
		SlicePtr: NewSlicePtr[M, *time.Time, *TimePtr[M]](key, NewTimePtr[M]),
	}
}

//
//
//

func (t *TimeSlice[M]) Max() *Time[M] {
	return NewTime[M](t.fn("time::max"))
}

func (t *TimeSlice[M]) Min() *Time[M] {
	return NewTime[M](t.fn("time::min"))
}
