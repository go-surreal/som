package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/parser"
)

type Struct struct {
	*baseField

	source *parser.FieldStruct
	model  Model
}

func (f *Struct) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(f.SourcePkg, f.model.NameGo())
}

func (f *Struct) typeConv() jen.Code {
	return jen.Add(f.ptr()).Id(f.model.NameGoLower())
}

func (f *Struct) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: nil,
		filterInit:   nil,
		filterFunc:   f.filterFunc,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil, // TODO

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *Struct) filterFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(f.model.NameGoLower()).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Struct).Types(jen.Id("T")).
				Params(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(f.NameDatabase())))))
}

func (f *Struct) convFrom(ctx Context) jen.Code {
	code := jen.Id("from" + f.model.NameGo())
	if f.source.Pointer() {
		code = jen.Id("ptrFunc").Call(jen.Id("from" + f.model.NameGo()))
	}
	return code.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Struct) convTo(ctx Context) jen.Code {
	code := jen.Id("to" + f.model.NameGo())
	if f.source.Pointer() {
		code = jen.Id("ptrFunc").Call(jen.Id("to" + f.model.NameGo()))
	}
	return code.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Struct) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
