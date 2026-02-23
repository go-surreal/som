package parser

import "github.com/wzshiming/gotype"

type enumFieldHandler struct{}



func (h *enumFieldHandler) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.String && isEnum(elem, ctx.OutPkg)
}

func (h *enumFieldHandler) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldEnum{&fieldAtomic{name: t.Name()}, elem.Name()}, nil
}
