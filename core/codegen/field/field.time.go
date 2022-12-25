package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type Time struct {
	*baseField

	source *parser.FieldTime
}

func (f *Time) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual("time", "Time")
}

func (f *Time) typeConv() jen.Code {
	return f.typeGo()
}

func (f *Time) CodeGen() *CodeGen {
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

func (f *Time) filterDefine(ctx Context) jen.Code {
	filter := "Time"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, filter).Types(jen.Id("T"))
}

func (f *Time) filterInit(ctx Context) jen.Code {
	filter := "NewTime"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(def.PkgLibFilter, filter).Types(jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(f.NameDatabase())))
}

func (f *Time) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibSort, "Sort").Types(jen.Id("T"))
}

func (f *Time) sortInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibSort, "NewSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Time) convFrom(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Time) convTo(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Time) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
