package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
)

type Edge struct {
	*baseField

	source *parser.FieldEdge
	table  *EdgeTable
}

func (f *Edge) typeGo() jen.Code {
	return jen.Add(f.ptr()).Qual(f.SourcePkg, f.table.NameGo())
}

func (f *Edge) typeConv() jen.Code {
	return jen.Add(f.ptr()).Id(f.table.NameGo())
}

func (f *Edge) TypeDatabase() string {
	return ""
}

func (f *Edge) Table() *EdgeTable {
	return f.table
}

func (f *Edge) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: nil,
		filterInit:   nil,
		filterFunc:   f.filterFunc,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil, // TODO: f.sortFunc, // edge currently not sortable

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *Edge) filterFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(f.table.NameGoLower()).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.table.NameGo()).Types(jen.Id("T")).
				Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())))))
}

//nolint:unused // currently not fully implemented
func (f *Edge) sortFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(f.NameGoLower()).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.table.NameGo()).Types(jen.Id("T")).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())))))
}

func (f *Edge) convFrom(_ Context) jen.Code {
	return jen.Id("From" + f.table.NameGo()).Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Edge) convTo(_ Context) jen.Code {
	return jen.Id("To" + f.table.NameGo()).Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Edge) fieldDef(_ Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}
