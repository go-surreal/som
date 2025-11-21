package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Bool struct {
	*baseField

	source *parser.FieldBool
}

func (f *Bool) typeGo() jen.Code {
	return jen.Add(f.ptr()).Bool()
}

func (f *Bool) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *Bool) TypeDatabase() string {
	return f.optionWrap("bool")
}

func (f *Bool) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil, // Bool does not need a filter function.

		sortDefine: nil, // TODO: should bool be sortable?
		sortInit:   nil,
		sortFunc:   nil, // Bool does not need a sort function.

		convFrom:      f.convFrom,
		convTo:        f.convTo,
		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
		fieldDef:      f.fieldDef,
	}
}

func (f *Bool) filterDefine(ctx Context) jen.Code {
	filter := "Bool"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(def.TypeModel)
}

func (f *Bool) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewBool"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Bool) convFrom(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo())
}

func (f *Bool) convTo(_ Context) (jen.Code, jen.Code) {
	return jen.Null(), jen.Id("data").Dot(f.NameGo())
}

func (f *Bool) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}

func (f *Bool) cborMarshal(_ Context) jen.Code {
	// Simple types: direct assignment
	if f.source.Pointer() {
		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()),
		)
	}
	return jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo())
}

func (f *Bool) cborUnmarshal(ctx Context) jen.Code {
	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
