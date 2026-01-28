package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type UUID struct {
	*baseField

	source *parser.FieldUUID
}

func (f *UUID) uuidPkg() string {
	switch f.source.Package {
	case parser.UUIDPackageGofrs:
		return def.PkgUUIDGofrs
	default:
		return def.PkgUUIDGoogle
	}
}

func (f *UUID) uuidTypeName() string {
	switch f.source.Package {
	case parser.UUIDPackageGofrs:
		return "UUIDGofrs"
	default:
		return "UUIDGoogle"
	}
}

func (f *UUID) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(f.uuidPkg(), "UUID")
}

func (f *UUID) typeConv(ctx Context) jen.Code {
	return jen.Add(f.ptr()).Qual(ctx.pkgTypes(), f.uuidTypeName())
}

func (f *UUID) TypeDatabase() string {
	if f.source.Pointer() {
		return "option<uuid>"
	}

	return "uuid"
}

func (f *UUID) SchemaStatements(table, prefix string) []string {
	return []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}
}

func (f *UUID) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		fieldDefine: f.fieldDefine,
		fieldInit:   f.fieldInit,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
	}
}

func (f *UUID) filterDefine(ctx Context) jen.Code {
	filter := f.uuidTypeName()
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *UUID) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "New" + f.uuidTypeName()
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(
			jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
		)
}

func (f *UUID) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(def.TypeModel)
}

func (f *UUID) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *UUID) fieldDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Qual(ctx.pkgQuery(), "Field").Types(def.TypeModel, jen.Qual(f.uuidPkg(), "UUID"))
}

func (f *UUID) fieldInit(ctx Context) jen.Code {
	factory := "New" + f.uuidTypeName() + "Field"
	if f.source.Pointer() {
		factory = "New" + f.uuidTypeName() + "PtrField"
	}
	return jen.Qual(ctx.pkgQuery(), factory).Types(def.TypeModel).
		Call(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *UUID) cborMarshal(ctx Context) jen.Code {
	typeName := f.uuidTypeName()
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).BlockFunc(func(bg *jen.Group) {
			bg.Id("uuidVal").Op(":=").Qual(ctx.pkgTypes(), typeName).Call(
				jen.Op("*").Id("c").Dot(f.NameGo()),
			)
			bg.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("uuidVal")
		})
	}

	return jen.BlockFunc(func(g *jen.Group) {
		g.Id("uuidVal").Op(":=").Qual(ctx.pkgTypes(), typeName).Call(
			jen.Id("c").Dot(f.NameGo()),
		)
		g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("uuidVal")
	})
}

func (f *UUID) cborUnmarshal(ctx Context) jen.Code {
	helper := "Unmarshal" + f.uuidTypeName()
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
