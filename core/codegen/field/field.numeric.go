package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type Numeric struct {
	*baseField

	source *parser.FieldNumeric
}

func (f *Numeric) typeGo() jen.Code {
	switch f.source.Type {
	case parser.NumberInt:
		return jen.Int()
	case parser.NumberInt32:
		return jen.Int32()
	case parser.NumberInt64:
		return jen.Int64()
	case parser.NumberFloat32:
		return jen.Float32()
	case parser.NumberFloat64:
		return jen.Float64()
	}
	return jen.Int() // TODO: okay?
}

func (f *Numeric) CodeGen() *CodeGen {
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

func (f *Numeric) filterDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Numeric").Types(f.typeGo(), jen.Id("T"))
}

func (f *Numeric) filterInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewNumeric").Types(f.typeGo(), jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Numeric) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibSort, "Sort").Types(jen.Id("T"))
}

func (f *Numeric) sortInit(ctx Context) jen.Code {
	return jen.Qual(def.PkgLibSort, "NewSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Numeric) convFrom(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Numeric) convTo(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Numeric) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeGo()).
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}
