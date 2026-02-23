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
	pkgPath := elem.PkgPath()
	return pkgPath == string(parser.UUIDPackageGoogle) || pkgPath == string(parser.UUIDPackageGofrs)
}

func (h *UUIDHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldUUID(t.Name(), false, parser.UUIDPackage(elem.PkgPath())), nil
}
