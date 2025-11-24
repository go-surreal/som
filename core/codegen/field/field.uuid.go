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

func (f *UUID) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(def.PkgUUID, "UUID")
}

func (f *UUID) typeConv(ctx Context) jen.Code {
	return jen.Add(f.ptr()).Qual(ctx.pkgTypes(), "UUID")
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

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,

		fieldDef: f.fieldDef,
	}
}

func (f *UUID) filterDefine(ctx Context) jen.Code {
	filter := "UUID"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *UUID) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewUUID"
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

func (f *UUID) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *UUID) cborMarshal(ctx Context) jen.Code {
	// Using custom types.UUID with MarshalCBOR method.
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).BlockFunc(func(bg *jen.Group) {
			bg.Id("uuidVal").Op(":=").Qual(ctx.pkgTypes(), "UUID").Call(
				jen.Op("*").Id("c").Dot(f.NameGo()),
			)
			bg.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("uuidVal")
		})
	}

	return jen.BlockFunc(func(g *jen.Group) {
		g.Id("uuidVal").Op(":=").Qual(ctx.pkgTypes(), "UUID").Call(
			jen.Id("c").Dot(f.NameGo()),
		)
		g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Op("&").Id("uuidVal")
	})
}

func (f *UUID) cborUnmarshal(ctx Context) jen.Code {
	helper := "UnmarshalUUID"
	if f.source.Pointer() {
		helper = "UnmarshalUUIDPtr"
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Id("c").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(ctx.pkgCBOR(), helper).Call(jen.Id("raw")),
	)
}
