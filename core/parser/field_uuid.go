package parser

import "github.com/wzshiming/gotype"

type uuidFieldHandler struct{}

func init() { RegisterFieldHandler(&uuidFieldHandler{}) }

func (h *uuidFieldHandler) Priority() int { return 40 }

func (h *uuidFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	if elem.Kind() != gotype.Array {
		return false
	}
	pkgPath := elem.PkgPath()
	return pkgPath == string(UUIDPackageGoogle) || pkgPath == string(UUIDPackageGofrs)
}

func (h *uuidFieldHandler) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldUUID{
		fieldAtomic: &fieldAtomic{name: t.Name()},
		Package:     UUIDPackage(elem.PkgPath()),
	}, nil
}
