package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type UUID struct {
	*baseField

	source *parser.FieldUUID
}

func (f *UUID) typeGo() jen.Code {
	return jen.Qual(def.PkgUUID, "UUID")
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
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Base").Types(jen.Qual(def.PkgUUID, "UUID"), jen.Id("T"))
}

func (f *UUID) filterInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewBase").Types(jen.Qual(def.PkgUUID, "UUID"), jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *UUID) convFrom(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo()).Dot("String").Call()
}

func (f *UUID) convTo(ctx Context) jen.Code {
	return jen.Id("parseUUID").Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *UUID) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).String().
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}
