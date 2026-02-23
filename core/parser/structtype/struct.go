package structtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type StructHandler struct{}

func (h *StructHandler) Match(t gotype.Type, _ *parser.TypeContext) bool {
	return t.Kind() == gotype.Struct
}

func (h *StructHandler) Handle(t gotype.Type, ctx *parser.TypeContext) error {
	str, err := parseStruct(t, ctx.OutPkg)
	if err != nil {
		return err
	}
	ctx.Output.Structs = append(ctx.Output.Structs, str)
	return nil
}

func (h *StructHandler) Validate(_ *parser.TypeContext) error { return nil }

func parseStruct(v gotype.Type, outPkg string) (*parser.Struct, error) {
	str := &parser.Struct{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		field, err := parser.ParseField(f, outPkg)
		if err != nil {
			return nil, err
		}

		str.Fields = append(str.Fields, field)
	}

	return str, nil
}
