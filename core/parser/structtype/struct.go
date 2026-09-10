package structtype

import (
	"fmt"
	"go/ast"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type StructHandler struct{}

func (h *StructHandler) Match(t gotype.Type, _ *parser.TypeContext) bool {
	return t.Kind() == gotype.Struct
}

func (h *StructHandler) Handle(t gotype.Type, ctx *parser.TypeContext) error {
	str, err := parseStruct(t, ctx)
	if err != nil {
		return err
	}
	ctx.Output.Structs = append(ctx.Output.Structs, str)
	return nil
}

func (h *StructHandler) Validate(ctx *parser.TypeContext) error {
	names := make([]string, len(ctx.Output.Structs))
	for i, s := range ctx.Output.Structs {
		names[i] = s.Name
		if err := validateFields("struct "+s.Name, s.Fields, ctx.Output); err != nil {
			return err
		}
	}
	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate struct name %q", dup)
	}
	return nil
}

func parseStruct(v gotype.Type, ctx *parser.TypeContext) (*parser.Struct, error) {
	str := &parser.Struct{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		field, err := ctx.ParseField(f)
		if err != nil {
			return nil, err
		}

		str.Fields = append(str.Fields, field)
	}

	return str, nil
}
