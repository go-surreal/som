package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
)

type Node struct {
	*baseField

	source *parser.FieldNode
	table  *NodeTable
}

func (f *Node) Source() *parser.FieldNode {
	return f.source
}

func (f *Node) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(f.SourcePkg, f.table.NameGo())
}

func (f *Node) typeConv(_ Context) jen.Code {
	return jen.Op("*").Id(f.table.NameGoLower() + "Link")
}

func (f *Node) TypeDatabase() string {
	// Linked records are always considered optional.
	return fmt.Sprintf("option<record<%s>>", f.table.NameDatabase())
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
	return jen.Id(f.table.NameGoLower()).Types(def.TypeModel)
}

func (f *Node) filterInit(_ Context) (jen.Code, jen.Code) {
	return jen.Id("new" + f.table.NameGo()).Types(def.TypeModel), nil
}

func (f *Node) filterFunc(ctx Context) jen.Code {
	receiver := jen.Id(ctx.Table.NameGoLower()).Types(def.TypeModel)
	if ctx.Receiver != nil {
		receiver = ctx.Receiver
	}

	return jen.Func().
		Params(jen.Id("n").Add(receiver)).
		Id(f.NameGo()).Params().
		Add(f.filterDefine(ctx)).
		Block(
			jen.Return(jen.Add(f.filterInit(ctx)).
				Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("Key"), jen.Lit(f.NameDatabase())))))
}

func (f *Node) sortFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(def.TypeModel)).
		Id(f.NameGo()).Params().
		Id(f.table.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(jen.Id("new" + f.table.NameGo()).Types(def.TypeModel).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())))))
}

func (f *Node) convFrom(_ Context) (jen.Code, jen.Code) {
	funcName := "to" + f.table.NameGo() + "Link"

	if f.source.Pointer() {
		funcName += fnSuffixPtr
	}

	return jen.Id(funcName),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Node) convTo(_ Context) (jen.Code, jen.Code) {
	funcName := "from" + f.table.NameGo() + "Link"

	if f.source.Pointer() {
		funcName += fnSuffixPtr
	}

	return jen.Id(funcName),
		jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Node) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + ",omitempty"})
}
