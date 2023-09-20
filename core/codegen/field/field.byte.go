package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
)

type Byte struct {
	*baseField

	source *parser.FieldByte
}

func (f *Byte) typeGo() jen.Code {
	return jen.Add(f.ptr()).Byte()
}

func (f *Byte) typeConv() jen.Code {
	return f.typeGo()
}

func (f *Byte) TypeDatabase() string {
	if f.source.Pointer() {
		return "string"
	}
	return "string ASSERT $value != NULL"
}

func (f *Byte) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   nil, // Bool does not need a filter function.

		sortDefine: nil, // TODO: should bool be sortable?
		sortInit:   nil,
		sortFunc:   nil, // Bool does not need a sort function.

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *Byte) filterDefine(ctx Context) jen.Code {
	filter := "Base"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(jen.Byte(), jen.Id("T"))
}

func (f *Byte) filterInit(ctx Context) jen.Code {
	filter := "NewBase"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(jen.Byte(), jen.Id("T")).
		Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Byte) convFrom(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Byte) convTo(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Byte) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
