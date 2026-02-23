package parser

import (
	"strings"

	"github.com/wzshiming/gotype"
)

type enumValueHandler struct{}

func init() { RegisterTypeHandler(&enumValueHandler{}) }

func (h *enumValueHandler) Priority() int { return 30 }

func (h *enumValueHandler) Match(t gotype.Type, _ *TypeContext) bool {
	return t.Kind() == gotype.Declaration
}

func (h *enumValueHandler) Handle(t gotype.Type, ctx *TypeContext) error {
	ctx.Output.EnumValues = append(ctx.Output.EnumValues, &EnumValue{
		Enum:     t.Declaration().Name(),
		Variable: t.Name(),
		Value:    strings.Trim(t.Value(), "\""),
	})
	return nil
}

func (h *enumValueHandler) Validate(_ *TypeContext) error { return nil }
