package fieldtype

import (
	"fmt"
	"strings"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type GeometryHandler struct{}

func (h *GeometryHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	pkg := elem.PkgPath()
	if pkg == string(parser.GeoPackageOrb) ||
		pkg == string(parser.GeoPackageSimplefeatures) ||
		pkg == string(parser.GeoPackageGoGeom) {
		_, ok := geoTypeName(pkg, elem.Name())
		return ok
	}
	if elem.Kind() == gotype.Invalid {
		s := elem.String()
		if strings.HasPrefix(s, "orb.") {
			name := strings.TrimPrefix(s, "orb.")
			_, ok := orbTypes[name]
			return ok
		}
		if pkg, name, ok := resolveGeomInvalid(s); ok {
			_, ok := geoTypeName(string(pkg), name)
			return ok
		}
	}
	return false
}

func (h *GeometryHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	pkg := elem.PkgPath()
	name := elem.Name()

	if pkg == "" && elem.Kind() == gotype.Invalid {
		s := elem.String()
		if strings.HasPrefix(s, "orb.") {
			pkg = string(parser.GeoPackageOrb)
			name = strings.TrimPrefix(s, "orb.")
		} else if p, n, ok := resolveGeomInvalid(s); ok {
			pkg = string(p)
			name = n
		}
	}

	geoType, ok := geoTypeName(pkg, name)
	if !ok {
		return nil, fmt.Errorf("unsupported geometry type: %s.%s", pkg, name)
	}
	return parser.NewFieldGeometry(t.Name(), parser.GeoPackage(pkg), geoType), nil
}

// resolveGeomInvalid attempts to resolve a gotype.Invalid string with a "geom." prefix
// to the correct geo package. Both simplefeatures and go-geom use package name "geom",
// so we check both type maps. Since the type names don't overlap between the two
// packages (e.g. "Collection" is orb-only, "GeometryCollection" is sf/go-geom),
// we can disambiguate by checking which map contains the name.
func resolveGeomInvalid(s string) (parser.GeoPackage, string, bool) {
	if !strings.HasPrefix(s, "geom.") {
		return "", "", false
	}
	name := strings.TrimPrefix(s, "geom.")
	if _, ok := sfTypes[name]; ok {
		return parser.GeoPackageSimplefeatures, name, true
	}
	if _, ok := goGeomTypes[name]; ok {
		return parser.GeoPackageGoGeom, name, true
	}
	return "", "", false
}

var orbTypes = map[string]parser.GeometryType{
	"Point":           parser.GeometryPoint,
	"LineString":      parser.GeometryLineString,
	"Polygon":         parser.GeometryPolygon,
	"MultiPoint":      parser.GeometryMultiPoint,
	"MultiLineString": parser.GeometryMultiLineString,
	"MultiPolygon":    parser.GeometryMultiPolygon,
	"Collection":      parser.GeometryCollection,
}

var sfTypes = map[string]parser.GeometryType{
	"Point":              parser.GeometryPoint,
	"LineString":         parser.GeometryLineString,
	"Polygon":            parser.GeometryPolygon,
	"MultiPoint":         parser.GeometryMultiPoint,
	"MultiLineString":    parser.GeometryMultiLineString,
	"MultiPolygon":       parser.GeometryMultiPolygon,
	"GeometryCollection": parser.GeometryCollection,
}

var goGeomTypes = map[string]parser.GeometryType{
	"Point":              parser.GeometryPoint,
	"LineString":         parser.GeometryLineString,
	"Polygon":            parser.GeometryPolygon,
	"MultiPoint":         parser.GeometryMultiPoint,
	"MultiLineString":    parser.GeometryMultiLineString,
	"MultiPolygon":       parser.GeometryMultiPolygon,
	"GeometryCollection": parser.GeometryCollection,
}

func geoTypeName(pkg string, name string) (parser.GeometryType, bool) {
	switch parser.GeoPackage(pkg) {
	case parser.GeoPackageOrb:
		t, ok := orbTypes[name]
		return t, ok
	case parser.GeoPackageSimplefeatures:
		t, ok := sfTypes[name]
		return t, ok
	case parser.GeoPackageGoGeom:
		t, ok := goGeomTypes[name]
		return t, ok
	}
	return 0, false
}
