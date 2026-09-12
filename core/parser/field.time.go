package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldTime is a time.Time. Besides a user-declared timestamp it also backs
// the som-managed timestamps contributed by the Timestamps, SoftDelete and
// Expiry feature embeds.
type FieldTime struct {
	fieldBase
	IsCreatedAt bool
	IsUpdatedAt bool
	IsDeletedAt bool
	IsExpiresAt bool
	// ExpiresIn holds the TTL duration (SurrealDB literal) for IsExpiresAt fields.
	ExpiresIn string
}

func (f *FieldTime) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Struct && elem.PkgPath() == "time" && elem.Name() == "Time"
}

func (f *FieldTime) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldTime{fieldBase: newBase(t.Name())}, nil
}

func (f *FieldTime) IsInternal() bool {
	return f.IsCreatedAt || f.IsUpdatedAt || f.IsDeletedAt || f.IsExpiresAt
}
