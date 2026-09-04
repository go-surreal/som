package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type UUIDHandler struct{}

func (h *UUIDHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	if elem.Kind() != gotype.Array {
		return false
	}
	switch parser.UUIDPackage(elem.PkgPath()) {
	case parser.UUIDPackageGoogle, parser.UUIDPackageGofrs, parser.UUIDPackageStd:
		return true
	default:
		return false
	}
}

func (h *UUIDHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldUUID(t.Name(), parser.UUIDPackage(elem.PkgPath())), nil
}
