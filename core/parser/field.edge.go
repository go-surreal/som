package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

// FieldEdge is a reference to an edge model.
type FieldEdge struct {
	fieldBase
	Edge string
}

func (f *FieldEdge) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.Struct && IsEdge(elem, ctx.OutPkg)
}

func (f *FieldEdge) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldEdge{fieldBase: newBase(t.Name()), Edge: elem.Name()}, nil
}

func (f *FieldEdge) Verify(out *Output) error {
	if err := f.fieldBase.Verify(out); err != nil {
		return err
	}

	if !edgeExists(f.Edge, out) {
		return fmt.Errorf("field %q references unknown edge %q", f.Name(), f.Edge)
	}

	return nil
}
