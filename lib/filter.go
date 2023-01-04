package lib

import (
	"strings"
	"time"
)

type Filter[T any] Where

func Filters[T any](filters []Filter[T]) []Filter[any] {
	var mapped []Filter[any]

	for _, filter := range filters {
		mapped = append(mapped, Filter[any](filter))
	}

	return mapped
}

func ToWhere[T any](filters []Filter[T]) []Where {
	var mapped []Where

	for _, filter := range filters {
		mapped = append(mapped, Where(filter))
	}

	return mapped
}

//
// -- BASE
//

type Base[T any, R any] struct {
	key     Key
	isCount bool
}

func NewBase[T, R any](key Key) *Base[T, R] {
	return &Base[T, R]{key: key}
}

func (b *Base[T, R]) Equal(val T) Filter[R] {
	return Filter[R](b.key.Op(OpEqual, val))
}

func (b *Base[T, R]) NotEqual(val T) Filter[R] {
	return Filter[R](b.key.Op(OpNotEqual, val))
}

func (b *Base[T, R]) In(vals []T) Filter[R] {
	return Filter[R](b.key.Op(OpInside, vals))
}

func (b *Base[T, R]) NotIn(vals []T) Filter[R] {
	return Filter[R](b.key.Op(OpNotInside, vals))
}

type BasePtr[T, R any] struct {
	*Base[T, R]
	*Nillable[R]
}

func NewBasePtr[T, R any](key Key) *BasePtr[T, R] {
	return &BasePtr[T, R]{
		Base:     &Base[T, R]{key: key},
		Nillable: &Nillable[R]{key: key},
	}
}

//
// -- COMPARABLE
//

type Comparable[T any, R any] struct {
	key     Key
	isCount bool
}

func (c *Comparable[T, R]) LessThan(val T) Filter[R] {
	return Filter[R](c.key.Op(OpLessThan, val))
}

func (c *Comparable[T, R]) LessThanEqual(val T) Filter[R] {
	return Filter[R](c.key.Op(OpLessThanEqual, val))
}

func (c *Comparable[T, R]) GreaterThan(val T) Filter[R] {
	return Filter[R](c.key.Op(OpGreaterThan, val))
}

func (c *Comparable[T, R]) GreaterThanEqual(val T) Filter[R] {
	return Filter[R](c.key.Op(OpGreaterThanEqual, val))
}

//
// -- NILLABLE
//

type Nillable[R any] struct {
	key Key
}

func (n *Nillable[R]) Nil() Filter[R] {
	return Filter[R](n.key.Op(OpExactlyEqual, nil))
}

func (n *Nillable[R]) NotNil() Filter[R] {
	return Filter[R](n.key.Op(OpNotEqual, nil))
}

//
// -- ID
//

type ID[R any] struct {
	key  Key
	node string
}

func NewID[R any](key Key, node string) *ID[R] {
	return &ID[R]{key: key, node: node}
}

func (b *ID[R]) Equal(val string) Filter[R] {
	val = b.node + ":" + val
	return Filter[R](b.key.Op(OpEqual, val))
}

func (b *ID[R]) NotEqual(val string) Filter[R] {
	val = b.node + ":" + val
	return Filter[R](b.key.Op(OpNotEqual, val))
}

func (b *ID[R]) In(vals []string) Filter[R] {
	var mapped []string

	for _, val := range vals {
		mapped = append(mapped, b.node+":"+val)
	}

	return Filter[R](b.key.Op(OpInside, mapped))
}

func (b *ID[R]) NotIn(vals []string) Filter[R] {
	var mapped []string

	for _, val := range vals {
		mapped = append(mapped, b.node+":"+val)
	}

	return Filter[R](b.key.Op(OpNotInside, mapped))
}

//
// -- STRING
//

type String[R any] struct {
	*Base[string, R]
	*Comparable[string, R]
}

func NewString[R any](key Key) *String[R] {
	return &String[R]{
		Base:       &Base[string, R]{key: key},
		Comparable: &Comparable[string, R]{key: key},
	}
}

func (s *String[R]) FuzzyMatch(val string) Filter[R] {
	return Filter[R](s.Base.key.Op(OpFuzzyMatch, val))
}

func (s *String[R]) NotFuzzyMatch(val string) Filter[R] {
	return Filter[R](s.Base.key.Op(OpFuzzyNotMatch, val))
}

type StringPtr[R any] struct {
	*String[R]
	*Nillable[R]
}

func NewStringPtr[R any](key Key) *StringPtr[R] {
	return &StringPtr[R]{
		String:   NewString[R](key),
		Nillable: &Nillable[R]{key: key},
	}
}

//
// -- NUMERIC
//

type Numeric[T, R any] struct {
	*Base[T, R]
	*Comparable[T, R]
}

func NewNumeric[T, R any](key Key) *Numeric[T, R] {
	return &Numeric[T, R]{
		Base:       &Base[T, R]{key: key},
		Comparable: &Comparable[T, R]{key: key},
	}
}

type NumericPtr[T, R any] struct {
	*Numeric[T, R]
	*Nillable[R]
}

func NewNumericPtr[T, R any](key Key) *NumericPtr[T, R] {
	return &NumericPtr[T, R]{
		Numeric:  NewNumeric[T, R](key),
		Nillable: &Nillable[R]{key: key},
	}
}

//
// -- BOOL
//

type Bool[R any] struct {
	key Key
}

func NewBool[R any](key Key) *Bool[R] {
	return &Bool[R]{key: key}
}

func (b *Bool[R]) Is(val bool) Filter[R] {
	return Filter[R](b.key.Op(OpExactlyEqual, val))
}

type BoolPtr[R any] struct {
	*Bool[R]
	*Nillable[R]
}

func NewBoolPtr[R any](key Key) *BoolPtr[R] {
	return &BoolPtr[R]{
		Bool:     &Bool[R]{key: key},
		Nillable: &Nillable[R]{key: key},
	}
}

//
// -- TIME
//

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

func (t *Time[R]) Before(val time.Time) Filter[R] {
	return t.comp.LessThan(val)
}

func (t *Time[R]) BeforeOrEqual(val time.Time) Filter[R] {
	return t.comp.LessThanEqual(val)
}

func (t *Time[R]) After(val time.Time) Filter[R] {
	return t.comp.GreaterThan(val)
}

func (t *Time[R]) AfterOrEqual(val time.Time) Filter[R] {
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

//
// -- SLICE
//

// Slice is a filter that can be used for slice fields.
// T is the type of the outgoing table for the filter statement.
// E is the type of the slice elements.
type Slice[T, E any] struct {
	key Key
}

// NewSlice creates a new slice filter.
func NewSlice[T, E any](key Key) *Slice[T, E] {
	return &Slice[T, E]{
		key: key,
	}
}

func (s *Slice[T, E]) Contains(val T) Filter[T] {
	return Filter[T](s.key.Op(OpContains, val))
}

func (s *Slice[T, E]) ContainsNot(val T) Filter[T] {
	return Filter[T](s.key.Op(OpContainsNot, val))
}

func (s *Slice[T, E]) ContainsAll(vals []T) Filter[T] {
	return Filter[T](s.key.Op(OpContainsAll, vals))
}

func (s *Slice[T, E]) ContainsAny(vals []T) Filter[T] {
	return Filter[T](s.key.Op(OpContainsAny, vals))
}

func (s *Slice[T, E]) ContainsNone(vals []T) Filter[T] {
	return Filter[T](s.key.Op(OpContainsNone, vals))
}

func (s *Slice[T, E]) Count() *Numeric[int, T] {
	return NewNumeric[int, T](s.key.Count())
}

type SlicePtr[T, E any] struct {
	*Slice[T, E]
	*Nillable[T]
}

func NewSlicePtr[T, E, F any](key Key) *SlicePtr[T, E] {
	return &SlicePtr[T, E]{
		Slice:    &Slice[T, E]{key: key},
		Nillable: &Nillable[T]{key: key},
	}
}

//
// -- ALL | ANY
//

func All(filters []Where) Where {
	return func(ctx *context) string {
		if len(filters) < 1 {
			return ""
		}

		var parts []string
		for _, filter := range filters {
			if part := filter(ctx); part != "" {
				parts = append(parts, part)
			}
		}

		if len(parts) < 1 {
			return ""
		}

		return "(" + strings.Join(parts, " "+string(OpAnd)+" ") + ")"
	}
}

func Any(filters []Where) Where {
	return func(ctx *context) string {
		if len(filters) < 1 {
			return ""
		}

		var parts []string
		for _, filter := range filters {
			if part := filter(ctx); part != "" {
				parts = append(parts, part)
			}
		}

		if len(parts) < 1 {
			return ""
		}

		return "(" + strings.Join(parts, " "+string(OpOr)+" ") + ")"
	}
}
