package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

// FieldNode is a record link to a node model.
type FieldNode struct {
	fieldBase
	Node string
}

func (f *FieldNode) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.Struct && IsNode(elem, ctx.OutPkg)
}

func (f *FieldNode) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldNode{fieldBase: newBase(t.Name()), Node: elem.Name()}, nil
}

func (f *FieldNode) Verify(out *Output) error {
	if err := f.fieldBase.Verify(out); err != nil {
		return err
	}

	if !nodeExists(f.Node, out) {
		return fmt.Errorf("field %q references unknown node %q", f.Name(), f.Node)
	}

	return nil
}
