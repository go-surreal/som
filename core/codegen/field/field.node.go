package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/parser"
)

type Node struct {
	*baseField

	source *parser.FieldNode
}

// TODO: cool to expose just like that?
func (f *Node) NodeName() string {
	return f.source.Node
}

func (f *Node) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: nil,
		filterInit:   nil,
		filterFunc:   f.filterFunc,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   f.sortFunc,

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *Node) filterFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Elem.NameGoLower()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(strcase.ToLowerCamel(f.source.Node)).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Node).Types(jen.Id("T")).
				Params(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))))
}

func (f *Node) sortFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Elem.NameDatabase()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(strcase.ToLowerCamel(f.source.Node)).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Node).Types(jen.Id("T")).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))))
}

func (f *Node) convFrom(ctx Context) jen.Code {
	return jen.Id("to" + f.source.Node + "Field").Call(jen.Op("&").Id("data").Dot(f.NameGo()))
}

func (f *Node) convTo(ctx Context) jen.Code {
	return jen.Op("*").Id("from" + f.source.Node + "Field").Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Node) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Id(f.source.Node + "Field").
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}
