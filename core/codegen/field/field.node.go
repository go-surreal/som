package field

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/parser"
)

type Node struct {
	*baseField

	source *parser.FieldNode
	table  *NodeTable
}

func (f *Node) typeGo() jen.Code {
	return jen.Qual(f.SourcePkg, f.table.NameGo())
}

func (f *Node) typeConv() jen.Code {
	return jen.Add(f.ptr()).Id(f.table.NameGoLower() + "Link")
}

func (f *Node) TypeDatabase() string {
	if f.source.Pointer() {
		return fmt.Sprintf("record(%s)", f.table.NameDatabase())
	}
	return fmt.Sprintf("record(%s) ASSERT $value != NULL", f.table.NameDatabase()) // TODO: how does it behave with empty struct?
}

func (f *Node) Table() *NodeTable {
	return f.table
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
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(f.table.NameGoLower()).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.table.NameGo()).Types(jen.Id("T")).
				Params(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(f.NameDatabase())))))
}

func (f *Node) sortFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameDatabase()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(f.table.NameGoLower()).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.table.NameGo()).Types(jen.Id("T")).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())))))
}

func (f *Node) convFrom(ctx Context) jen.Code {
	return jen.Id("to" + f.table.NameGo() + "Link").Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Node) convTo(ctx Context) jen.Code {
	return jen.Id("from" + f.table.NameGo() + "Link").Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Node) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
