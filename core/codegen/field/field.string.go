package field

import (
	"github.com/dave/jennifer/jen"
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
		fieldDef: f.fieldDef,
	}
}

func (f *String) filterDefine(ctx Context) jen.Code {
	filter := "String"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(typeModel)
}

func (f *String) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewString"
	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(typeModel),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *String) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "StringSort").Types(typeModel)
}

func (f *String) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewStringSort").Types(typeModel).
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
		Tag(map[string]string{convTag: f.NameDatabase()})
}
