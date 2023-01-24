package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type UUID struct {
	*baseField

	source *parser.FieldUUID
}

func (f *UUID) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(def.PkgUUID, "UUID")
}

func (f *UUID) typeConv() jen.Code {
	return jen.Add(f.ptr()).String()
}

func (f *UUID) TypeDatabase() string {
	if f.source.Pointer() {
		return "string"
	}
	return "string ASSERT $value != NULL" // TODO: assert for uuid? (currently fails for unknown reason!)
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

	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLib, filter).Types(jen.Qual(def.PkgUUID, "UUID"), jen.Id("T"))
}

func (f *UUID) filterInit(ctx Context) jen.Code {
	filter := "NewBase"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(def.PkgLib, filter).Types(jen.Qual(def.PkgUUID, "UUID"), jen.Id("T")).
		Params(jen.Qual(def.PkgLib, "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *UUID) convFrom(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.Id("uuidPtr").Call(jen.Id("data").Dot(f.NameGo()))
	}
	return jen.Id("data").Dot(f.NameGo()).Dot("String").Call()
}

func (f *UUID) convTo(ctx Context) jen.Code {
	if f.source.Pointer() {
		return jen.Id("ptrFunc").Call(jen.Id("parseUUID")).Call(jen.Id("data").Dot(f.NameGo()))
	}
	return jen.Id("parseUUID").Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *UUID) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
