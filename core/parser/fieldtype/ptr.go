package fieldtype

import (
	"fmt"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type PtrHandler struct{}

func (h *PtrHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Ptr
}

func (h *PtrHandler) Parse(t gotype.Type, elem gotype.Type, ctx *parser.FieldContext) (parser.Field, error) {
	field, err := ctx.RecurseParse(elem, elem.Elem(), ctx)
	if err != nil {
		return nil, fmt.Errorf("could not parse elem for ptr field %s: %v", t.Name(), err)
	}

	if t.Name() != "" {
		parser.SetFieldName(field, t.Name())
	}
	parser.SetFieldPointer(field, true)

	return field, nil
}
