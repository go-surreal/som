package parser

import (
	"fmt"
	"strings"

	"github.com/wzshiming/gotype"
)

type GeoPackage string

const (
	GeoPackageOrb            GeoPackage = "github.com/paulmach/orb"
	GeoPackageSimplefeatures GeoPackage = "github.com/peterstace/simplefeatures/geom"
	GeoPackageGoGeom         GeoPackage = "github.com/twpayne/go-geom"
)

type GeometryType int

const (
	GeometryPoint GeometryType = iota
	GeometryLineString
	GeometryPolygon
	GeometryMultiPoint
	GeometryMultiLineString
	GeometryMultiPolygon
	GeometryCollection
)

func (g GeometryType) String() string {
	switch g {
	case GeometryPoint:
		return "Point"
	case GeometryLineString:
		return "LineString"
	case GeometryPolygon:
		return "Polygon"
	case GeometryMultiPoint:
		return "MultiPoint"
	case GeometryMultiLineString:
		return "MultiLineString"
	case GeometryMultiPolygon:
		return "MultiPolygon"
	case GeometryCollection:
		return "Collection"
	default:
		return "Unknown"
	}
}

// FieldGeometry is a geometry value from one of the supported geo packages.
type FieldGeometry struct {
	fieldBase
	Package GeoPackage
	Type    GeometryType
}

func (f *FieldGeometry) Match(elem gotype.Type, _ *FieldContext) bool {
	pkg := elem.PkgPath()
	if pkg == string(GeoPackageOrb) ||
		pkg == string(GeoPackageSimplefeatures) ||
		pkg == string(GeoPackageGoGeom) {
		_, ok := geoTypeName(pkg, elem.Name())
		return ok
	}
	if elem.Kind() == gotype.Invalid {
		s := elem.String()
		if name, ok := strings.CutPrefix(s, "orb."); ok {
			_, found := orbTypes[name]
			return found
		}
		if pkg, name, ok := resolveGeomInvalid(s); ok {
			_, ok := geoTypeName(string(pkg), name)
			return ok
		}
	}
	return false
}

func (f *FieldGeometry) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	pkg := elem.PkgPath()
	name := elem.Name()

	if pkg == "" && elem.Kind() == gotype.Invalid {
		s := elem.String()
		if trimmed, ok := strings.CutPrefix(s, "orb."); ok {
			pkg = string(GeoPackageOrb)
			name = trimmed
		} else if p, n, ok := resolveGeomInvalid(s); ok {
			pkg = string(p)
			name = n
		}
	}

	geoType, ok := geoTypeName(pkg, name)
	if !ok {
		return nil, fmt.Errorf("unsupported geometry type: %s.%s", pkg, name)
	}
	return &FieldGeometry{fieldBase: newBase(t.Name()), Package: GeoPackage(pkg), Type: geoType}, nil
}

func (f *FieldGeometry) contributeFeatures(features *UsedFeatures) {
	switch f.Package {
	case GeoPackageOrb:
		features.UsesOrbGeo = true
	case GeoPackageSimplefeatures:
		features.UsesSimplefeaturesGeo = true
	case GeoPackageGoGeom:
		features.UsesGoGeomGeo = true
	}
}

// resolveGeomInvalid attempts to resolve a gotype.Invalid string with a "geom." prefix
// to the correct geo package. Both simplefeatures and go-geom use package name "geom",
// so we check both type maps. Since the type names don't overlap between the two
// packages (e.g. "Collection" is orb-only, "GeometryCollection" is sf/go-geom),
// we can disambiguate by checking which map contains the name.
func resolveGeomInvalid(s string) (GeoPackage, string, bool) {
	name, ok := strings.CutPrefix(s, "geom.")
	if !ok {
		return "", "", false
	}
	if _, ok := sfTypes[name]; ok {
		return GeoPackageSimplefeatures, name, true
	}
	if _, ok := goGeomTypes[name]; ok {
		return GeoPackageGoGeom, name, true
	}
	return "", "", false
}

var orbTypes = map[string]GeometryType{
	"Point":           GeometryPoint,
	"LineString":      GeometryLineString,
	"Polygon":         GeometryPolygon,
	"MultiPoint":      GeometryMultiPoint,
	"MultiLineString": GeometryMultiLineString,
	"MultiPolygon":    GeometryMultiPolygon,
	"Collection":      GeometryCollection,
}

var sfTypes = map[string]GeometryType{
	"Point":              GeometryPoint,
	"LineString":         GeometryLineString,
	"Polygon":            GeometryPolygon,
	"MultiPoint":         GeometryMultiPoint,
	"MultiLineString":    GeometryMultiLineString,
	"MultiPolygon":       GeometryMultiPolygon,
	"GeometryCollection": GeometryCollection,
}

var goGeomTypes = map[string]GeometryType{
	"Point":              GeometryPoint,
	"LineString":         GeometryLineString,
	"Polygon":            GeometryPolygon,
	"MultiPoint":         GeometryMultiPoint,
	"MultiLineString":    GeometryMultiLineString,
	"MultiPolygon":       GeometryMultiPolygon,
	"GeometryCollection": GeometryCollection,
}

func geoTypeName(pkg string, name string) (GeometryType, bool) {
	switch GeoPackage(pkg) {
	case GeoPackageOrb:
		t, ok := orbTypes[name]
		return t, ok
	case GeoPackageSimplefeatures:
		t, ok := sfTypes[name]
		return t, ok
	case GeoPackageGoGeom:
		t, ok := goGeomTypes[name]
		return t, ok
	}
	return 0, false
}
