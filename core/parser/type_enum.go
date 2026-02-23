package parser

import "github.com/wzshiming/gotype"

type enumHandler struct{}



func (h *enumHandler) Match(t gotype.Type, ctx *TypeContext) bool {
	return isEnum(t, ctx.OutPkg)
}

func (h *enumHandler) Handle(t gotype.Type, ctx *TypeContext) error {
	ctx.Output.Enums = append(ctx.Output.Enums, &Enum{
		Name: t.Name(),
	})
	return nil
}

func (h *enumHandler) Validate(_ *TypeContext) error { return nil }

func isEnum(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.String {
		return false
	}

	return t.String() != "string" && t.PkgPath() == outPkg
}
