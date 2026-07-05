package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Geometry struct {
	*baseField

	source *parser.FieldGeometry
}

func (f *Geometry) geoPkg() string {
	switch f.source.Package {
	case parser.GeoPackageSimplefeatures:
		return def.PkgGeoSimplefeatures
	case parser.GeoPackageGoGeom:
		return def.PkgGeoGoGeom
	default:
		return def.PkgGeoOrb
	}
}

// geoGoTypeName returns the Go type name from the external library.
func (f *Geometry) geoGoTypeName() string {
	switch f.source.Package {
	case parser.GeoPackageSimplefeatures, parser.GeoPackageGoGeom:
		if f.source.Type == parser.GeometryCollection {
			return "GeometryCollection"
		}
		return f.source.Type.String()
	default:
		return f.source.Type.String()
	}
}

// geoConvTypeName returns the internal wrapper type name for CBOR marshaling.
func (f *Geometry) geoConvTypeName() string {
	switch f.source.Package {
	case parser.GeoPackageSimplefeatures:
		return f.source.Type.String() + "SF"
	case parser.GeoPackageGoGeom:
		return f.source.Type.String() + "GG"
	default:
		return f.source.Type.String() + "Orb"
	}
}

// geoFilterTypeName returns the filter type name.
func (f *Geometry) geoFilterTypeName() string {
	switch f.source.Package {
	case parser.GeoPackageSimplefeatures:
		return "Geo" + f.source.Type.String() + "SF"
	case parser.GeoPackageGoGeom:
		return "Geo" + f.source.Type.String() + "GG"
	default:
		return "Geo" + f.source.Type.String() + "Orb"
	}
}

func (f *Geometry) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(f.geoPkg(), f.geoGoTypeName())
}

func (f *Geometry) typeConv(ctx Context) jen.Code {
	return jen.Add(f.ptr()).Qual(ctx.pkgTypes(), f.geoConvTypeName())
}

// TypeDatabase returns the SurrealDB type for the geometry field.
func (f *Geometry) TypeDatabase() string {
	var geoType string
	switch f.source.Type {
	case parser.GeometryPoint:
		geoType = "geometry<point>"
	case parser.GeometryLineString:
		geoType = "geometry<line>"
	case parser.GeometryPolygon:
		geoType = "geometry<polygon>"
	case parser.GeometryMultiPoint:
		geoType = "geometry<multipoint>"
	case parser.GeometryMultiLineString:
		geoType = "geometry<multiline>"
	case parser.GeometryMultiPolygon:
		geoType = "geometry<multipolygon>"
	case parser.GeometryCollection:
		geoType = "geometry<collection>"
	default:
		geoType = "geometry"
	}

	if f.source.Pointer() {
		return "option<" + geoType + ">"
	}

	return geoType
}

func (f *Geometry) SchemaStatements(table, prefix string) []string {
	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}
}

func (f *Geometry) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: nil, // geo types are not sortable
		sortInit:   nil,
		sortFunc:   nil,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
	}
}

func (f *Geometry) filterDefine(ctx Context) jen.Code {
	filter := f.geoFilterTypeName()
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Geometry) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "New" + f.geoFilterTypeName()
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(
			jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
		)
}

func (f *Geometry) geoConvExpr(ctx Context, val jen.Code) jen.Code {
	typeName := f.geoConvTypeName()
	if f.source.Package == parser.GeoPackageGoGeom {
		return jen.Qual(ctx.pkgTypes(), typeName).Values(val)
	}
	return jen.Qual(ctx.pkgTypes(), typeName).Call(val)
}

func (f *Geometry) cborMarshal(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).BlockFunc(func(bg *jen.Group) {
			bg.Id("geoVal").Op(":=").Add(f.geoConvExpr(ctx,
				jen.Op("*").Id("c").Dot(f.NameGo()),
			))
			bg.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("geoVal")
		})
	}

	return jen.BlockFunc(func(g *jen.Group) {
		g.Id("geoVal").Op(":=").Add(f.geoConvExpr(ctx,
			jen.Id("c").Dot(f.NameGo()),
		))
		g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("geoVal")
	})
}

func (f *Geometry) cborUnmarshal(ctx Context) jen.Code {
	helper := "Unmarshal" + f.geoConvTypeName()
	if f.source.Pointer() {
		helper += "Ptr"
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Id("c").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(ctx.pkgCBOR(), helper).Call(jen.Id("raw")),
	)
}
