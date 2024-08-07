package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type URL struct {
	*baseField

	source *parser.FieldURL
}

func (f *URL) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(def.PkgURL, "URL")
}

func (f *URL) typeConv(_ Context) jen.Code {
	return jen.Add(f.ptr()).String()
}

func (f *URL) TypeDatabase() string {
	if f.source.Pointer() {
		return "option<string | null> ASSERT $value == NONE OR $value == NULL OR string::is::url($value)"
		// TODO: should field be omitted (omitempty) if value is null (instead of being set to null)?
	}

	return `string ASSERT $value == "" OR string::is::url($value)`
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
	filter := "URL"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(jen.Id("T"))
}

func (f *URL) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewURL"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(jen.Id("T")),
		jen.Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *URL) convFrom(_ Context) (jen.Code, jen.Code) {
	fromFunc := "fromURL"

	if f.source.Pointer() {
		fromFunc += "Ptr"
	}

	return jen.Id(fromFunc),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *URL) convTo(_ Context) (jen.Code, jen.Code) {
	toFunc := "toURL"

	if f.source.Pointer() {
		toFunc += "Ptr"
	}

	return jen.Id(toFunc),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *URL) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{"json": f.NameDatabase()})
}
