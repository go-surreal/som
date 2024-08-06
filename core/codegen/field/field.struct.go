package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
)

type Struct struct {
	*baseField

	source *parser.FieldStruct
	table  Table
}

func (f *Struct) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(f.SourcePkg, f.table.NameGo())
}

func (f *Struct) typeConv() jen.Code {
	return jen.Add(f.ptr()).Id(f.table.NameGoLower())
}

func (f *Struct) TypeDatabase() string {
	return f.optionWrap("object")
}

func (f *Struct) Table() Table {
	return f.table
}

func (f *Struct) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   f.filterFunc,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil, // TODO

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *Struct) filterDefine(_ Context) jen.Code {
	return jen.Id(f.table.NameGoLower()).Types(jen.Id("T"))
}

func (f *Struct) filterInit(_ Context) (jen.Code, jen.Code) {
	return jen.Id("new" + f.source.Struct).Types(jen.Id("T")), nil
}

func (f *Struct) filterFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Add(f.filterDefine(ctx)).
		Block(
			jen.Return(jen.Add(f.filterInit(ctx)).
				Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())))))
}

func (f *Struct) convFrom(_ Context) (jen.Code, jen.Code) {
	code := jen.Id("from" + f.table.NameGo())

	if !f.source.Pointer() {
		code = jen.Id("noPtrFunc").Call(jen.Id("from" + f.table.NameGo()))
	}

	return code,
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Struct) convTo(_ Context) (jen.Code, jen.Code) {
	code := jen.Id("to" + f.table.NameGo())

	if !f.source.Pointer() {
		code = jen.Id("noPtrFunc").Call(jen.Id("to" + f.table.NameGo()))
	}
	
	return code,
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Struct) fieldDef(_ Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
