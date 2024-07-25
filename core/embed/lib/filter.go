//go:build embed

package lib

import (
	"github.com/go-surreal/sdbc"
	"net/url"
	"strings"
	"time"
)

type Filter[T any] interface {
	build(*context, T) string
}

type filter[T any] func(*context, T) string

//nolint:unused
func (f filter[T]) build(ctx *context, t T) string {
	return f(ctx, t)
}

func KeyFilter[T any](key Key[T]) Filter[T] {
	return filter[T](func(ctx *context, _ T) string {
		return key.render(ctx)
	})
}

//
// -- BASE
//

// TODO: switch T and R
type Base[T any, R any] struct {
	key Key[R]
}

func NewBase[T, R any](key Key[R]) *Base[T, R] {
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

func NewBasePtr[T, R any](key Key[R]) *BasePtr[T, R] {
	return &BasePtr[T, R]{
		Base:     &Base[T, R]{key: key},
		Nillable: &Nillable[R]{key: key},
	}
}

//
// -- COMPARABLE
//

type Comparable[T any, R any] struct {
	key Key[R]
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
	key Key[R]
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
	key  Key[R]
	node string
}

func NewID[R any](key Key[R], node string) *ID[R] {
	return &ID[R]{key: key, node: node}
}

func (b *ID[R]) Equal(val *sdbc.ID) Filter[R] {
	// val = b.node + ":" + val
	return Filter[R](b.key.Op(OpEqual, val))
}

func (b *ID[R]) NotEqual(val *sdbc.ID) Filter[R] {
	// val = b.node + ":" + val
	return Filter[R](b.key.Op(OpNotEqual, val))
}

func (b *ID[R]) In(vals []*sdbc.ID) Filter[R] {
	return Filter[R](b.key.Op(OpInside, vals))
}

func (b *ID[R]) NotIn(vals []*sdbc.ID) Filter[R] {
	return Filter[R](b.key.Op(OpNotInside, vals))
}

//
// -- STRING
//

type String[R any] struct {
	*Base[string, R]
	*Comparable[string, R]
}

func NewString[R any](key Key[R]) *String[R] {
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

func NewStringPtr[R any](key Key[R]) *StringPtr[R] {
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

func NewNumeric[T, R any](key Key[R]) *Numeric[T, R] {
	return &Numeric[T, R]{
		Base:       &Base[T, R]{key: key},
		Comparable: &Comparable[T, R]{key: key},
	}
}

type NumericPtr[T, R any] struct {
	*Numeric[T, R]
	*Nillable[R]
}

func NewNumericPtr[T, R any](key Key[R]) *NumericPtr[T, R] {
	return &NumericPtr[T, R]{
		Numeric:  NewNumeric[T, R](key),
		Nillable: &Nillable[R]{key: key},
	}
}

//
// -- BOOL
//

type Bool[R any] struct {
	key Key[R]
}

func NewBool[R any](key Key[R]) *Bool[R] {
	return &Bool[R]{key: key}
}

func (b *Bool[R]) Is(val bool) Filter[R] {
	return Filter[R](b.key.Op(OpExactlyEqual, val))
}

type BoolPtr[R any] struct {
	*Bool[R]
	*Nillable[R]
}

func NewBoolPtr[R any](key Key[R]) *BoolPtr[R] {
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

func NewTime[R any](key Key[R]) *Time[R] {
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

func NewTimePtr[R any](key Key[R]) *TimePtr[R] {
	return &TimePtr[R]{
		Time:     NewTime[R](key),
		Nillable: &Nillable[R]{key: key},
	}
}

//
// -- DURATION
//

type Duration[R any] struct {
	*Base[time.Duration, R]
	comp *Comparable[time.Time, R]
}

func NewDuration[R any](key Key[R]) *Duration[R] {
	return &Duration[R]{
		Base: &Base[time.Duration, R]{key: key},
		comp: &Comparable[time.Time, R]{key: key},
	}
}

func (t *Duration[R]) Before(val time.Time) Filter[R] {
	return t.comp.LessThan(val)
}

func (t *Duration[R]) BeforeOrEqual(val time.Time) Filter[R] {
	return t.comp.LessThanEqual(val)
}

func (t *Duration[R]) After(val time.Time) Filter[R] {
	return t.comp.GreaterThan(val)
}

func (t *Duration[R]) AfterOrEqual(val time.Time) Filter[R] {
	return t.comp.GreaterThanEqual(val)
}

type DurationPtr[R any] struct {
	*Duration[R]
	*Nillable[R]
}

func NewDurationPtr[R any](key Key[R]) *DurationPtr[R] {
	return &DurationPtr[R]{
		Duration: NewDuration[R](key),
		Nillable: &Nillable[R]{key: key},
	}
}

// -- URL
//

type URL[T any] struct {
	key Key[T]
}

func NewURL[T any](key Key[T]) *URL[T] {
	return &URL[T]{key: key}
}

func (b *URL[T]) Equal(val url.URL) Filter[T] {
	return Filter[T](b.key.Op(OpEqual, val.String()))
}

func (b *URL[T]) NotEqual(val url.URL) Filter[T] {
	return Filter[T](b.key.Op(OpNotEqual, val.String()))
}

func (b *URL[T]) In(vals []url.URL) Filter[T] {
	var mapped []string

	for _, val := range vals {
		mapped = append(mapped, val.String())
	}

	return Filter[T](b.key.Op(OpInside, mapped))
}

func (b *URL[T]) NotIn(vals []url.URL) Filter[T] {
	var mapped []string

	for _, val := range vals {
		mapped = append(mapped, val.String())
	}

	return Filter[T](b.key.Op(OpNotInside, mapped))
}

type URLPtr[T any] struct {
	*URL[T]
	*Nillable[T]
}

func NewURLPtr[T any](key Key[T]) *URLPtr[T] {
	return &URLPtr[T]{
		URL:      &URL[T]{key: key},
		Nillable: &Nillable[T]{key: key},
	}
}

//
// -- SLICE
//

// Slice is a filter that can be used for slice fields.
// T is the type of the outgoing table for the filter statement.
// E is the type of the slice elements.
type Slice[T, E any] struct {
	key Key[T]
}

// NewSlice creates a new slice filter.
func NewSlice[T, E any](key Key[T]) *Slice[T, E] {
	return &Slice[T, E]{
		key: key,
	}
}

func (s *Slice[T, E]) Contains(val E) Filter[T] {
	return Filter[T](s.key.Op(OpContains, val))
}

func (s *Slice[T, E]) ContainsNot(val E) Filter[T] {
	return Filter[T](s.key.Op(OpContainsNot, val))
}

func (s *Slice[T, E]) ContainsAll(vals []E) Filter[T] {
	return Filter[T](s.key.Op(OpContainsAll, vals))
}

func (s *Slice[T, E]) ContainsAny(vals []E) Filter[T] {
	return Filter[T](s.key.Op(OpContainsAny, vals))
}

func (s *Slice[T, E]) ContainsNone(vals []E) Filter[T] {
	return Filter[T](s.key.Op(OpContainsNone, vals))
}

func (s *Slice[T, E]) Count() *Numeric[int, T] {
	return NewNumeric[int, T](s.key.Count())
}

type SlicePtr[T, E any] struct {
	*Slice[T, E]
	*Nillable[T]
}

func NewSlicePtr[T, E, F any](key Key[T]) *SlicePtr[T, E] {
	return &SlicePtr[T, E]{
		Slice:    &Slice[T, E]{key: key},
		Nillable: &Nillable[T]{key: key},
	}
}

//
// -- BYTE SLICE
//

// ByteSlice is a filter that can be used for byte slice fields.
// T is the type of the outgoing table for the filter statement.
type ByteSlice[T any] struct {
	*Base[[]byte, T]
}

// NewSlice creates a new slice filter.
func NewByteSlice[T any](key Key[T]) *ByteSlice[T] {
	return &ByteSlice[T]{
		Base: &Base[[]byte, T]{key: key},
	}
}

type ByteSlicePtr[T any] struct {
	*ByteSlice[T]
	*Nillable[T]
}

func NewByteSlicePtr[T any](key Key[T]) *ByteSlicePtr[T] {
	return &ByteSlicePtr[T]{
		ByteSlice: &ByteSlice[T]{
			Base: &Base[[]byte, T]{key: key},
		},
		Nillable: &Nillable[T]{key: key},
	}
}

//
// -- NODE SLICE
//

// type NodeSlice[T, N any] struct {
// 	key Key[T]
// }
//
// func NewNodeSlice[T, N any](key Key[T]) *NodeSlice[T, N] {
// 	return &NodeSlice[T, N]{key: key}
// }
//
// func (s *NodeSlice[T, N]) build(_ *context, _ T) string {
// 	return "" // TODO
// }
//
// func (s *NodeSlice[T, N]) Count() *Numeric[int, T] {
// 	return NewNumeric[int, T](s.key.Count())
// }

//
// -- ALL | ANY
//

type All[T any] []Filter[T]

func (a All[T]) build(ctx *context, t T) string {
	if len(a) < 1 {
		return ""
	}

	var parts []string
	for _, filter := range a {
		if part := filter.build(ctx, t); part != "" {
			parts = append(parts, strings.TrimPrefix(part, ".")) // TODO: better place to trim?
		}
	}

	if len(parts) < 1 {
		return ""
	}

	return "(" + strings.Join(parts, " "+string(OpAnd)+" ") + ")"
}

type Any[T any] []Filter[T]

//nolint:unused
func (a Any[T]) build(ctx *context, t T) string {
	if len(a) < 1 {
		return ""
	}

	var parts []string
	for _, filter := range a {
		if part := filter.build(ctx, t); part != "" {
			parts = append(parts, strings.TrimPrefix(part, ".")) // TODO: better place to trim?
		}
	}

	if len(parts) < 1 {
		return ""
	}

	return "(" + strings.Join(parts, " "+string(OpOr)+" ") + ")"
}
