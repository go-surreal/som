package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type UUID struct {
	*baseField

	source *parser.FieldUUID
}

func (f *UUID) Source() *parser.FieldUUID {
	return f.source
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

func (f *UUID) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil,

		convFrom: f.convFrom,
		convTo:   f.convTo,

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

func (f *UUID) convFrom(_ Context) (jen.Code, jen.Code) {
	fromFunc := "fromUUID"

	if f.source.Pointer() {
		fromFunc += fnSuffixPtr
	}

	return jen.Id(fromFunc),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *UUID) convTo(_ Context) (jen.Code, jen.Code) {
	toFunc := "toUUID"

	if f.source.Pointer() {
		toFunc += fnSuffixPtr
	}

	return jen.Id(toFunc),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *UUID) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *UUID) cborMarshal(ctx Context) jen.Code {
	helper := "marshalUUID"
	if f.source.Pointer() {
		helper += "Ptr"
	}

	return jen.BlockFunc(func(g *jen.Group) {
		if f.source.Pointer() {
			g.If(jen.Id("m").Dot(f.NameGo()).Op("!=").Nil()).Block(
				jen.Id("val").Op(",").Id("_").Op(":=").Qual(ctx.pkgCBOR(), helper).Call(
					jen.Id("m").Dot(f.NameGo()),
				),
				jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Qual(def.PkgCBOR, "RawMessage").Call(jen.Id("val")),
			)
		} else {
			g.Id("val").Op(",").Id("_").Op(":=").Qual(ctx.pkgCBOR(), helper).Call(
				jen.Id("m").Dot(f.NameGo()),
			)
			g.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Qual(def.PkgCBOR, "RawMessage").Call(jen.Id("val"))
		}
	})
}

func (f *UUID) cborUnmarshal(ctx Context) jen.Code {
	helper := "unmarshalUUID"
	if f.source.Pointer() {
		helper += "Ptr"
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Id("m").Dot(f.NameGo()).Op(",").Id("_").Op("=").Qual(ctx.pkgCBOR(), helper).Call(jen.Id("raw")),
	)
}
