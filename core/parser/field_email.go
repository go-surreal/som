package parser

import "github.com/wzshiming/gotype"

type emailFieldHandler struct{}



func (h *emailFieldHandler) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.String && elem.PkgPath() == ctx.OutPkg && elem.Name() == "Email"
}

func (h *emailFieldHandler) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldEmail{&fieldAtomic{name: t.Name()}}, nil
}
