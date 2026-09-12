package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldWeekday is a time.Weekday.
type FieldWeekday struct {
	fieldBase
}

func (f *FieldWeekday) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Int && elem.PkgPath() == "time" && elem.Name() == "Weekday"
}

func (f *FieldWeekday) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldWeekday{fieldBase: newBase(t.Name())}, nil
}
