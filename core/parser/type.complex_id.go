package parser

import (
	"github.com/wzshiming/gotype"
)

// ComplexIDStructHandler claims the key structs used as the type argument of a
// som.Node embed, e.g. PersonID for som.Node[PersonID]. They are resolved by
// the node that references them (see ParseComplexIDFields) and must not end up
// in Output.Structs, so this handler matches them and produces nothing.
type ComplexIDStructHandler struct{}

func (h *ComplexIDStructHandler) Match(t gotype.Type, ctx *TypeContext) bool {
	return IsComplexIDStruct(t, ctx.OutPkg)
}

func (h *ComplexIDStructHandler) Handle(_ gotype.Type, _ *TypeContext) error {
	return nil
}

func (h *ComplexIDStructHandler) Validate(_ *TypeContext) error {
	return nil
}

func IsComplexIDStruct(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}
	nf := t.NumField()
	for i := range nf {
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
