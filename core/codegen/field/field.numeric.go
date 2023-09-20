package field

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
)

type Numeric struct {
	*baseField

	source *parser.FieldNumeric
}

func (f *Numeric) typeGo() jen.Code {
	switch f.source.Type {
	case parser.NumberInt:
		return jen.Add(f.ptr()).Int()
	case parser.NumberInt32:
		return jen.Add(f.ptr()).Int32()
	case parser.NumberInt64:
		return jen.Add(f.ptr()).Int64()
	case parser.NumberFloat32:
		return jen.Add(f.ptr()).Float32()
	case parser.NumberFloat64:
		return jen.Add(f.ptr()).Float64()
	}
	return jen.Empty() // this case can basically not happen ;)
}

func (f *Numeric) typeConv() jen.Code {
	return f.typeGo()
}

func (f *Numeric) TypeDatabase() string {
	switch f.source.Type {

	case parser.NumberInt, parser.NumberInt32, parser.NumberInt64:
		return f.optionWrap("int")

	case parser.NumberFloat32, parser.NumberFloat64:
		return f.optionWrap("float")

	default:
		panic(fmt.Sprintf("unmapped numeric type: %d", f.source.Type))
	}
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
	filter := "Numeric"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(f.typeGo(), jen.Id("T"))
}

func (f *Numeric) filterInit(ctx Context) jen.Code {
	filter := "NewNumeric"
	if f.source.Pointer() {
		filter += "Ptr"
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(f.typeGo(), jen.Id("T")).
		Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Numeric) sortDefine(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "BaseSort").Types(jen.Id("T"))
}

func (f *Numeric) sortInit(ctx Context) jen.Code {
	return jen.Qual(ctx.pkgLib(), "NewBaseSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(f.NameDatabase())))
}

func (f *Numeric) convFrom(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Numeric) convTo(ctx Context) jen.Code {
	return jen.Id("data").Dot(f.NameGo())
}

func (f *Numeric) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
