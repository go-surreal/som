package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type Enum struct {
	*baseField

	source *parser.FieldEnum
	model  Model
}

func (f *Enum) typeGo() jen.Code {
	return jen.Qual(f.SourcePkg, f.model.NameGo())
}

func (f *Enum) CodeGen() *CodeGen {
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

func (f *Enum) filterDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Base").Types(jen.Qual(ctx.SourcePkg, f.source.Typ), jen.Id("T"))
}

func (f *Enum) filterInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewBase").Types(jen.Qual(ctx.SourcePkg, f.source.Typ), jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Enum) convFrom(ctx Context) jen.Code {
	return jen.String().Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Enum) convTo(ctx Context) jen.Code {
	return jen.Qual(ctx.SourcePkg, f.source.Typ).Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Enum) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).String(). // TODO: support other enum base types (atomic)
						Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}
