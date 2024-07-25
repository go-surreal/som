package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Duration struct {
	*baseField

	source *parser.FieldDuration
}

func (f *Duration) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual("time", "Duration")
}

func (f *Duration) typeConv() jen.Code {
	return jen.Add(f.ptr()).Qual(def.PkgSDBC, "Duration")
}

func (f *Duration) TypeDatabase() string {
	return f.optionWrap("duration")
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

func (f *Duration) convFrom(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.Id("fromDurationPtr").Call(jen.Id("data").Dot(f.NameGo()))
	}

	return jen.Qual(def.PkgSDBC, "Duration").Values(jen.Id("data").Dot(f.NameGo()))
}

func (f *Duration) convTo(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.Id("toDurationPtr").Call(jen.Id("data").Dot(f.NameGo()))
	}

	return jen.Id("data").Dot(f.NameGo()).Dot("Duration")
}

func (f *Duration) fieldDef(_ Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
