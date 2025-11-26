//go:build embed

package lib

import "strings"

// Nillable is a filter builder for values that can be nil.
// M is the type of the model the filter is for.
type Nillable[M any] struct {
	Key[M]
}

// NewNillable creates a new nillable filter builder.
// M is the type of the model the filter is for.
func NewNillable[M any](key Key[M]) *Nillable[M] {
	return &Nillable[M]{Key: key}
}

// Nil returns a filter that checks if the value is nil.
//
// This checks for both NONE (field absent) and NULL (field explicitly null),
// because in golang both cases are represented as nil.
func (n *Nillable[M]) Nil(is bool) Filter[M] {
	if is {
		return filter[M](func(ctx *context, _ M) string {
			fieldName := strings.TrimPrefix(n.render(ctx), ".")
			return "(" + fieldName + " == NONE OR " + fieldName + " == NULL)"
		})
	}

	return filter[M](func(ctx *context, _ M) string {
		fieldName := strings.TrimPrefix(n.render(ctx), ".")
		return "(" + fieldName + " != NONE AND " + fieldName + " != NULL)"
	})
}
