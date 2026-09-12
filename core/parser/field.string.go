package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldString is a plain Go string. It is the only scalar that supports a
// fulltext index.
type FieldString struct {
	fieldBase
}

func NewFieldString(name string) *FieldString {
	return &FieldString{fieldBase: newBase(name)}
}

func (f *FieldString) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.String
}

func (f *FieldString) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return NewFieldString(t.Name()), nil
}

func (f *FieldString) Validate() error {
	return nil // string and *string support fulltext
}
