package parser

import "github.com/wzshiming/gotype"

type structHandler struct{}

func init() { RegisterTypeHandler(&structHandler{}) }

func (h *structHandler) Priority() int { return 90 }

func (h *structHandler) Match(t gotype.Type, _ *TypeContext) bool {
	return t.Kind() == gotype.Struct
}

func (h *structHandler) Handle(t gotype.Type, ctx *TypeContext) error {
	str, err := parseStruct(t, ctx.OutPkg)
	if err != nil {
		return err
	}
	ctx.Output.Structs = append(ctx.Output.Structs, str)
	return nil
}

func (h *structHandler) Validate(_ *TypeContext) error { return nil }

func parseStruct(v gotype.Type, outPkg string) (*Struct, error) {
	str := &Struct{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		field, err := parseField(f, outPkg)
		if err != nil {
			return nil, err
		}

		str.Fields = append(str.Fields, field)
	}

	return str, nil
}
