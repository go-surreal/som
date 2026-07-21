package structtype

import (
	"fmt"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type EnumHandler struct{}

func (h *EnumHandler) Match(t gotype.Type, ctx *parser.TypeContext) bool {
	return IsEnum(t, ctx.OutPkg)
}

func (h *EnumHandler) Handle(t gotype.Type, ctx *parser.TypeContext) error {
	ctx.Output.Enums = append(ctx.Output.Enums, &parser.Enum{
		Name: t.Name(),
	})
	return nil
}

func (h *EnumHandler) Validate(ctx *parser.TypeContext) error {
	names := make([]string, len(ctx.Output.Enums))
	for i, e := range ctx.Output.Enums {
		names[i] = e.Name
	}
	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate enum name %q", dup)
	}
	return nil
}

func IsEnum(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.String {
		return false
	}

	return t.String() != "string" && t.PkgPath() == outPkg
}
