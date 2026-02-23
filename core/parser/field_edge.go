package parser

import "github.com/wzshiming/gotype"

type edgeRefFieldHandler struct{}



func (h *edgeRefFieldHandler) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.Struct && isEdge(elem, ctx.OutPkg)
}

func (h *edgeRefFieldHandler) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldEdge{&fieldAtomic{name: t.Name()}, elem.Name()}, nil
}
