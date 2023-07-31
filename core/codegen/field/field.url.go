package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type URL struct {
	*baseField

	source *parser.FieldURL
}

func (f *URL) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(def.PkgURL, "URL")
}

func (f *URL) typeConv() jen.Code {
	return jen.Add(f.ptr()).String()
}

func (f *URL) TypeDatabase() string {
	if f.source.Pointer() {
		return "option<string> ASSERT $value == NONE OR $value == NULL OR is::url($value)"
	}

	return `string ASSERT $value == "" OR is::url($value)`
}

func (f *URL) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil,

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *URL) filterDefine(ctx Context) jen.Code {
	filter := "Base"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(jen.Qual(def.PkgURL, "URL"), jen.Id("T"))
}

func (f *URL) filterInit(ctx Context) jen.Code {
	filter := "NewBase"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(jen.Qual(def.PkgURL, "URL"), jen.Id("T")).
		Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *URL) convFrom(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.Id("urlPtr").Call(jen.Id("data").Dot(f.NameGo()))
	}
	return jen.Id("data").Dot(f.NameGo()).Dot("String").Call()
}

func (f *URL) convTo(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.Id("ptrFunc").Call(jen.Id("parseURL")).Call(jen.Id("data").Dot(f.NameGo()))
	}
	return jen.Id("parseURL").Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *URL) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}