package parser

import (
	"github.com/wzshiming/gotype"
)

type UUIDPackage string

const (
	UUIDPackageGoogle UUIDPackage = "github.com/google/uuid"
	UUIDPackageGofrs  UUIDPackage = "github.com/gofrs/uuid"
	UUIDPackageStd    UUIDPackage = "uuid"
)

// FieldUUID is a UUID from one of the supported UUID packages.
type FieldUUID struct {
	fieldBase
	Package UUIDPackage
}

func (f *FieldUUID) Match(elem gotype.Type, _ *FieldContext) bool {
	if elem.Kind() != gotype.Array {
		return false
	}
	switch UUIDPackage(elem.PkgPath()) {
	case UUIDPackageGoogle, UUIDPackageGofrs, UUIDPackageStd:
		return true
	default:
		return false
	}
}

func (f *FieldUUID) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldUUID{fieldBase: newBase(t.Name()), Package: UUIDPackage(elem.PkgPath())}, nil
}

func (f *FieldUUID) contributeFeatures(features *UsedFeatures) {
	switch f.Package {
	case UUIDPackageGoogle:
		features.UsesGoogleUUID = true
	case UUIDPackageGofrs:
		features.UsesGofrsUUID = true
	case UUIDPackageStd:
		features.UsesStdUUID = true
	}
}
