package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/parser"
)

type Duration struct {
	*baseField

	source *parser.FieldDuration
}

func (f *Duration) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual("time", "Duration")
}

func (f *Duration) typeConv() jen.Code {
	return jen.Add(f.ptr()).String()
}

func (f *Duration) TypeDatabase() string {
	return f.optionWrap("duration") + " VALUE type::duration($value)"
}

func (f *Duration) CodeGen() *CodeGen {
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

func (f *Duration) filterDefine(ctx Context) jen.Code {
	filter := "Duration"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(jen.Id("T"))
}

func (f *Duration) filterInit(ctx Context) jen.Code {
	filter := "NewDuration"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(jen.Id("T")).
		Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Duration) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(jen.Id("T"))
}

func (f *Duration) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Duration) convFrom(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.Id("durationPtr").Call(jen.Id("data").Dot(f.NameGo()))
	}

	return jen.Id("data").Dot(f.NameGo()).Dot("String").Call()
}

func (f *Duration) convTo(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.Id("ptrFunc").Call(jen.Id("parseDuration")).Call(jen.Id("data").Dot(f.NameGo()))
	}

	return jen.Id("parseDuration").Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Duration) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
