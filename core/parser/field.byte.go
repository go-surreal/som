package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldByte is a plain Go byte.
type FieldByte struct {
	fieldBase
}

func (f *FieldByte) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Byte
}

func (f *FieldByte) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldByte{fieldBase: newBase(t.Name())}, nil
}
