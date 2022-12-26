package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type Bool struct {
	*baseField

	source *parser.FieldBool
}

func (f *Bool) typeGo() jen.Code {
	return jen.Add(f.ptr()).Bool()
}

func (f *Bool) typeConv() jen.Code {
	return f.typeGo()
}

func (f *Bool) TypeDatabase() string {
	return "bool"
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
	filter := "Bool"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, filter).Types(jen.Id("T"))
}

func (f *Bool) filterInit(ctx Context) jen.Code {
	filter := "NewBool"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(def.PkgLibFilter, filter).Types(jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(f.NameDatabase())))
}

func (f *Bool) convFrom(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Bool) convTo(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Bool) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
