package parser

import (
	"strings"

	"github.com/wzshiming/gotype"
)

// EnumValue is a declared constant of an Enum type. It is also its own
// TypeHandler.
type EnumValue struct {
	Enum     string
	Variable string
	Value    string
}

func (e *EnumValue) Match(t gotype.Type, _ *TypeContext) bool {
	return t.Kind() == gotype.Declaration
}

func (e *EnumValue) Handle(t gotype.Type, ctx *TypeContext) error {
	ctx.Output.EnumValues = append(ctx.Output.EnumValues, &EnumValue{
		Enum:     t.Declaration().Name(),
		Variable: t.Name(),
		Value:    strings.Trim(t.Value(), "\""),
	})
	return nil
}

func (e *EnumValue) Validate(_ *TypeContext) error {
	return nil
}
