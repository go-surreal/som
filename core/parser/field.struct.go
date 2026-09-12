package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

// FieldStruct is a nested struct, stored as an embedded object.
//
// It matches any struct and must therefore be registered last, after every
// handler for a struct with a more specific meaning (time.Time, url.URL, node
// and edge references, ...).
type FieldStruct struct {
	fieldBase
	Struct string
}

func (f *FieldStruct) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Struct
}

func (f *FieldStruct) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldStruct{fieldBase: newBase(t.Name()), Struct: elem.Name()}, nil
}

func (f *FieldStruct) Verify(out *Output) error {
	if err := f.fieldBase.Verify(out); err != nil {
		return err
	}

	if !structExists(f.Struct, out) {
		return fmt.Errorf("field %q references unknown struct %q", f.Name(), f.Struct)
	}

	return nil
}
