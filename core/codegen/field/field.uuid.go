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

func (f *UUID) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(def.PkgUUID, "UUID")
}

func (f *UUID) typeConv() jen.Code {
	return jen.Add(f.ptr()).Id("UUID")
}

func (f *UUID) TypeDatabase() string {
	if f.source.Pointer() {
		return "option<uuid | null>"
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
		fieldDef: f.fieldDef,
	}
}

func (f *UUID) filterDefine(ctx Context) jen.Code {
	filter := "Base"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(jen.Qual(def.PkgUUID, "UUID"), jen.Id("T"))
}

func (f *UUID) filterInit(ctx Context) jen.Code {
	filter := "NewBase"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(jen.Qual(def.PkgUUID, "UUID"), jen.Id("T")).
		Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *UUID) convFrom(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.Parens(jen.Op("*").Id("UUID")).Call(jen.Id("data").Dot(f.NameGo()))
	}

	return jen.Id("UUID").Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *UUID) convTo(_ Context) jen.Code {
	if f.source.Pointer() {
		return jen.Parens(jen.Op("*").Qual(def.PkgUUID, "UUID")).Call(jen.Id("data").Dot(f.NameGo()))
	}

	return jen.Qual(def.PkgUUID, "UUID").Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *UUID) fieldDef(_ Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
