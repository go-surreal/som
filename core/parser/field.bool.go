package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldBool is a plain Go bool.
type FieldBool struct {
	fieldBase
}

func (f *FieldBool) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Bool
}

func (f *FieldBool) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldBool{fieldBase: newBase(t.Name())}, nil
}
