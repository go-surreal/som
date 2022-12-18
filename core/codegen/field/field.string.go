package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type String struct {
	*baseField

	source *parser.FieldString
}

func (f *String) typeGo() jen.Code {
	return jen.String()
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
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "String").Types(jen.Id("T"))
}

func (f *String) filterInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewString").Types(jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *String) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibSort, "String").Types(jen.Id("T"))
}

func (f *String) sortInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibSort, "NewString").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *String) convFrom(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo()) // TODO: vulnerability -> record link could be injected
}

func (f *String) convTo(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *String) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).String().
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}
