package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type Bool struct {
	*baseField

	source *parser.FieldBool
}

func (f *Bool) typeGo() jen.Code {
	return jen.Bool()
}

func (f *Bool) CodeGen() *CodeGen {
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

func (f *Bool) filterDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Bool").Types(jen.Id("T"))
}

func (f *Bool) filterInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewBool").Types(jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Bool) convFrom(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Bool) convTo(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Bool) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Bool().
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"}) // TODO: store "false" (no omitempty)?
}
