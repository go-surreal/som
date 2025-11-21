package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type String struct {
	*baseField

	source *parser.FieldString
}

func (f *String) typeGo() jen.Code {
	return jen.Add(f.ptr()).String()
}

func (f *String) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *String) TypeDatabase() string {
	return f.optionWrap("string")
}

func (f *String) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: f.sortDefine,
		sortInit:   f.sortInit,
		sortFunc:   nil,

		convFrom: f.convFrom,
		convTo:   f.convTo,

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,

		fieldDef: f.fieldDef,
	}
}

func (f *String) filterDefine(ctx Context) jen.Code {
	filter := "String"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *String) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewString"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *String) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "StringSort").Types(def.TypeModel)
}

func (f *String) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewStringSort").Types(def.TypeModel).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *String) convFrom(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo()) // TODO: record link inject vulnerability? solved by cbor?
}

func (f *String) convTo(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo())
}

func (f *String) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *String) cborMarshal(_ Context) jen.Code {
	// Simple types: direct assignment
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()),
		)
	}
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())
}

func (f *String) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
