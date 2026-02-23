package parser

import "github.com/wzshiming/gotype"

type nodeRefFieldHandler struct{}



func (h *nodeRefFieldHandler) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.Struct && isNode(elem, ctx.OutPkg)
}

func (h *nodeRefFieldHandler) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldNode{&fieldAtomic{name: t.Name()}, elem.Name()}, nil
}
