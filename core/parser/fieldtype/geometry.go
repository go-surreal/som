package fieldtype

import (
	"strings"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type GeometryHandler struct{}

func (h *GeometryHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	pkg := elem.PkgPath()
	if pkg == string(parser.GeoPackageOrb) || pkg == string(parser.GeoPackageSimplefeatures) {
		_, ok := geoTypeName(pkg, elem.Name())
		return ok
	}
	if elem.Kind() == gotype.Invalid {
		s := elem.String()
		return strings.HasPrefix(s, "orb.") || strings.HasPrefix(s, "geom.")
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
		} else if strings.HasPrefix(s, "geom.") {
			pkg = string(parser.GeoPackageSimplefeatures)
			name = strings.TrimPrefix(s, "geom.")
		}
	}

	geoType, _ := geoTypeName(pkg, name)
	return parser.NewFieldGeometry(t.Name(), parser.GeoPackage(pkg), geoType), nil
}

var orbTypes = map[string]parser.GeometryType{
	"Point":            parser.GeometryPoint,
	"LineString":       parser.GeometryLineString,
	"Polygon":          parser.GeometryPolygon,
	"MultiPoint":       parser.GeometryMultiPoint,
	"MultiLineString":  parser.GeometryMultiLineString,
	"MultiPolygon":     parser.GeometryMultiPolygon,
	"Collection":       parser.GeometryCollection,
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

func geoTypeName(pkg string, name string) (parser.GeometryType, bool) {
	switch parser.GeoPackage(pkg) {
	case parser.GeoPackageOrb:
		t, ok := orbTypes[name]
		return t, ok
	case parser.GeoPackageSimplefeatures:
		t, ok := sfTypes[name]
		return t, ok
	}
	return 0, false
}
