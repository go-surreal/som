package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

// FieldUnion is a record link that may point to any member of a union.
type FieldUnion struct {
	fieldBase
	Union string
}

func (f *FieldUnion) Match(elem gotype.Type, ctx *FieldContext) bool {
	return IsUnion(elem, ctx.OutPkg)
}

func (f *FieldUnion) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldUnion{fieldBase: newBase(t.Name()), Union: elem.Name()}, nil
}

func (f *FieldUnion) Validate() error {
	if err := f.fieldBase.Validate(); err != nil {
		return err
	}

	if f.Pointer() {
		return fmt.Errorf("field %s: a union field must not be a pointer, an unset link is the nil interface", f.Name())
	}

	return nil
}

func (f *FieldUnion) Verify(out *Output) error {
	if err := f.fieldBase.Verify(out); err != nil {
		return err
	}

	if !unionExists(f.Union, out) {
		return fmt.Errorf("field %q references unknown union %q", f.Name(), f.Union)
	}

	return nil
}
