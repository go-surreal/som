package field

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
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
	return jen.Op("*").Id(f.table.NameGoLower() + "Link")
}

func (f *Node) TypeDatabase() string {
	// Linked records are always considered optional.
	return fmt.Sprintf("option<record<%s> | null>", f.table.NameDatabase())
}

func (f *Node) Table() *NodeTable {
	return f.table
}

func (f *Node) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   f.filterFunc,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   f.sortFunc,

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *Node) filterDefine(_ Context) jen.Code {
	return jen.Id(f.table.NameGoLower()).Types(jen.Id("T"))
}

func (f *Node) filterInit(_ Context) (jen.Code, jen.Code) {
	return jen.Id("new" + f.table.NameGo()).Types(jen.Id("T")), nil
}

func (f *Node) filterFunc(ctx Context) jen.Code {
	receiver := jen.Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))
	if ctx.Receiver != nil {
		receiver = ctx.Receiver
	}

	return jen.Func().
		Params(jen.Id("n").Add(receiver)).
		Id(f.NameGo()).Params().
		Add(f.filterDefine(ctx)).
		Block(
			jen.Return(jen.Add(f.filterInit(ctx)).
				Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())))))
}

func (f *Node) sortFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(f.table.NameGoLower()).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.table.NameGo()).Types(jen.Id("T")).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())))))
}

func (f *Node) convFrom(ctx Context) jen.Code {
	funcName := "to" + f.table.NameGo() + "Link"
	if f.source.Pointer() {
		funcName += "Ptr"
	}

	return jen.Id(funcName).Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Node) convTo(ctx Context) jen.Code {
	funcName := "from" + f.table.NameGo() + "Link"
	if f.source.Pointer() {
		funcName += "Ptr"
	}

	return jen.Id(funcName).Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Node) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase()})
}
