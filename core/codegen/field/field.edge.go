package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
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

func (f *Edge) typeConv(_ Context) jen.Code {
	return jen.Op("*").Id(f.table.NameGo())
}

func (f *Edge) TypeDatabase() string {
	return ""
}

func (f *Edge) Table() *EdgeTable {
	return f.table
}

func (f *Edge) CodeGen() *CodeGen {
	return &CodeGen{
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   f.filterFunc,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil, // TODO: f.sortFunc, // edge currently not sortable

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *Edge) filterDefine(_ Context) jen.Code {
	return jen.Id(f.table.NameGoLower()).Types(def.TypeModel)
}

func (f *Edge) filterInit(_ Context) (jen.Code, jen.Code) {
	return jen.Id("new" + f.table.NameGo()).Types(def.TypeModel), nil
}

func (f *Edge) filterFunc(ctx Context) jen.Code {
	receiver := jen.Id(ctx.Table.NameGoLower()).Types(def.TypeModel)
	if ctx.Receiver != nil {
		receiver = ctx.Receiver
	}

	return jen.Func().
		Params(jen.Id("n").Add(receiver)).
		Id(f.NameGo()).Params().
		Add(f.filterDefine(ctx)).
		Block(
			jen.Return(jen.Add(jen.Id("new" + f.table.NameGo()).Types(def.TypeModel)).
				Params(jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("Key"), jen.Lit(f.NameDatabase())))))
}

//nolint:unused // currently not fully implemented
func (f *Edge) sortFunc(ctx Context) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(def.TypeModel)).
		Id(f.NameGo()).Params().
		Id(f.NameGoLower()).Types(def.TypeModel).
		Block(
			jen.Return(jen.Id("new" + f.table.NameGo()).Types(def.TypeModel).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("Key"), jen.Lit(f.NameDatabase())))))
}

func (f *Edge) convFrom(_ Context) (jen.Code, jen.Code) {
	return nil, nil
}

func (f *Edge) convTo(_ Context) (jen.Code, jen.Code) {
	fn := "To" + f.table.NameGo()

	if f.source.Pointer() {
		fn += fnSuffixPtr
	}

	return jen.Id(fn), jen.Call(jen.Id("data").Dot(f.NameGo()))
}

func (f *Edge) fieldDef(ctx Context) jen.Code {
	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + f.omitEmptyIfPtr()})
}
