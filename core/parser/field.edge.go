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

func (f *FieldEdge) Validate() error {
	// Edges are not real fields in the database schema, so there is no
	// DEFINE FIELD statement a constraint could attach to.
	if err := f.rejectAsserts("an edge does not support constraints"); err != nil {
		return err
	}

	return f.fieldBase.Validate()
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
