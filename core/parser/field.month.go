package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldMonth is a time.Month.
type FieldMonth struct {
	fieldBase
}

func (f *FieldMonth) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Int && elem.PkgPath() == "time" && elem.Name() == "Month"
}

func (f *FieldMonth) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldMonth{fieldBase: newBase(t.Name())}, nil
}
