package structtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type ComplexIDStructHandler struct{}

func (h *ComplexIDStructHandler) Match(t gotype.Type, ctx *parser.TypeContext) bool {
	return IsComplexIDStruct(t, ctx.OutPkg)
}

func (h *ComplexIDStructHandler) Handle(_ gotype.Type, _ *parser.TypeContext) error {
	return nil
}

func (h *ComplexIDStructHandler) Validate(_ *parser.TypeContext) error { return nil }

func IsComplexIDStruct(t gotype.Type, outPkg string) bool {
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
