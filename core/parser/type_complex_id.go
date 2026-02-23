package parser

import "github.com/wzshiming/gotype"

type complexIDStructHandler struct{}



func (h *complexIDStructHandler) Match(t gotype.Type, ctx *TypeContext) bool {
	return isComplexIDStruct(t, ctx.OutPkg)
}

func (h *complexIDStructHandler) Handle(_ gotype.Type, _ *TypeContext) error {
	return nil
}

func (h *complexIDStructHandler) Validate(_ *TypeContext) error { return nil }

func isComplexIDStruct(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}
	nf := t.NumField()
	for i := 0; i < nf; i++ {
		f := t.Field(i)
		if !f.IsAnonymous() {
			continue
		}
		if (f.Name() == "ArrayID" || f.Name() == "ObjectID") && f.Elem().PkgPath() == outPkg {
			return true
		}
	}
	return false
}
