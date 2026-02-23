package structtype

import (
	"strings"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type EnumValueHandler struct{}

func (h *EnumValueHandler) Match(t gotype.Type, _ *parser.TypeContext) bool {
	return t.Kind() == gotype.Declaration
}

func (h *EnumValueHandler) Handle(t gotype.Type, ctx *parser.TypeContext) error {
	ctx.Output.EnumValues = append(ctx.Output.EnumValues, &parser.EnumValue{
		Enum:     t.Declaration().Name(),
		Variable: t.Name(),
		Value:    strings.Trim(t.Value(), "\""),
	})
	return nil
}

func (h *EnumValueHandler) Validate(_ *parser.TypeContext) error { return nil }
