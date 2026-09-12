package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldDuration is a time.Duration.
type FieldDuration struct {
	fieldBase
}

func (f *FieldDuration) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Int64 && elem.PkgPath() == "time" && elem.Name() == "Duration"
}

func (f *FieldDuration) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldDuration{fieldBase: newBase(t.Name())}, nil
}
