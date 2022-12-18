package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/parser"
)

type Struct struct {
	*baseField

	source *parser.FieldStruct
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
		Params(jen.Id("n").Id(ctx.Elem.NameGoLower()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(strcase.ToLowerCamel(f.source.Struct)).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Struct).Types(jen.Id("T")).
				Params(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))))
}

func (f *Struct) convFrom(ctx Context) jen.Code {
	return jen.Op("*").Id("From" + f.source.Struct).Call(jen.Op("&").Id("data").Dot(f.NameGo()))
}

func (f *Struct) convTo(ctx Context) jen.Code {
	return jen.Op("*").Id("To" + f.source.Struct).Call(jen.Op("&").Id("data").Dot(f.NameGo()))
}

func (f *Struct) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Id(f.source.Struct).
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}
