package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type Slice struct {
	*baseField

	source     *parser.FieldSlice
	getElement ElemGetter
}

// TODO: cleanup this (temporary) method!
func (f *Slice) Edge() (bool, string, string, string) {
	if !f.source.IsEdge {
		return false, "", "", ""
	}

	rawEdge, _ := f.getElement(f.source.Field.(*parser.FieldEdge).Edge)
	edge := rawEdge.(*DatabaseEdge)
	in := edge.In.(*Node).source.Node
	out := edge.Out.(*Node).source.Node
	field := f.source.Field.(*parser.FieldEdge).Edge

	return true, field, in, out
}

func (f *Slice) CodeGen() *CodeGen {
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

func (f *Slice) filterFunc(ctx Context) jen.Code {
	if f.source.IsEdge {
		rawEdge, _ := f.getElement(f.source.Field.(*parser.FieldEdge).Edge)
		edge := rawEdge.(*DatabaseEdge)
		in := edge.In.(*Node).source.Node
		out := edge.Out.(*Node).source.Node
		field := f.source.Field.(*parser.FieldEdge).Edge

		if ctx.Elem.NameGo() == in {
			return jen.Func().
				Params(jen.Id("n").Id(strcase.ToLowerCamel(ctx.Elem.NameGo())).Types(jen.Id("T"))).
				Id(f.NameGo()).Params().
				Params(jen.Id(strcase.ToLowerCamel(field) + "In").Index(jen.Id("T"))).
				Block(
					jen.Return(
						jen.Id("new" + field + "In").Index(jen.Id("T")).
							Call(jen.Id("n").Dot("key").Dot("In").Call(jen.Lit(edge.NameDatabase()))),
					),
				)
		}

		if ctx.Elem.NameGo() == out {
			return jen.Func().
				Params(jen.Id("n").Id(ctx.Elem.NameGoLower()).Types(jen.Id("T"))).
				Id(f.NameGo()).Params().
				Params(jen.Id(strcase.ToLowerCamel(field) + "Out").Index(jen.Id("T"))).
				Block(
					jen.Return(
						jen.Id("new" + field + "Out").Index(jen.Id("T")).
							Call(jen.Id("n").Dot("key").Dot("Out").Call(jen.Lit(edge.NameDatabase()))),
					),
				)
		}

		return nil
	} else if f.source.IsNode {
		return jen.Func().
			Params(jen.Id("n").Id(ctx.Elem.NameGoLower()).Types(jen.Id("T"))).
			Id(f.NameGo()).Params().
			Id(strcase.ToLowerCamel(f.source.Value)+"Slice").Types(jen.Id("T")).
			Block(
				jen.Id("key").Op(":=").Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))),
				jen.Return(
					jen.Id(strcase.ToLowerCamel(f.source.Value)+"Slice").Types(jen.Id("T")).
						Values(
							jen.Id("new"+f.source.Value).Types(jen.Id("T")).
								Call(jen.Id("key")),
							jen.Qual(def.PkgLibFilter, "NewSlice").Types(jen.Qual(ctx.SourcePkg, f.source.Value), jen.Id("T")).
								Call(jen.Id("key")),
						),
				),
			)
	} else if f.source.IsEnum {
		return jen.Func().
			Params(jen.Id("n").Id(ctx.Elem.NameGoLower()).Types(jen.Id("T"))).
			Id(f.NameGo()).Params().
			Op("*").Qual(def.PkgLibFilter, "Slice").Types(jen.Qual(ctx.SourcePkg, f.source.Value), jen.Id("T")).
			Block(
				jen.Return(
					jen.Qual(def.PkgLibFilter, "NewSlice").Types(jen.Qual(ctx.SourcePkg, f.source.Value), jen.Id("T")).
						Call(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo())))),
				),
			)
	} else {
		return jen.Func().
			Params(jen.Id("n").Id(ctx.Elem.NameGoLower()).Types(jen.Id("T"))).
			Id(f.NameGo()).Params().
			Op("*").Qual(def.PkgLibFilter, "Slice").Types(jen.Id(f.source.Value), jen.Id("T")).
			Block(
				jen.Return(
					jen.Qual(def.PkgLibFilter, "NewSlice").Types(jen.Id(f.source.Value), jen.Id("T")).
						Call(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo())))),
				),
			)
	}
}

func (f *Slice) convFrom(ctx Context) jen.Code {
	if f.source.IsNode {
		return jen.Id("mapRecords").Call(jen.Id("data").Dot(f.NameGo()), jen.Id("to"+f.source.Value+"Field"))
	} else if f.source.IsEdge {
		return nil // TODO: should edges really not be addable like that?
	} else if f.source.IsEnum {
		return jen.Id("convertEnum").Types(jen.Qual(ctx.SourcePkg, f.source.Value), jen.String()).Call(jen.Id("data").Dot(f.NameGo()))
	} else {
		return jen.Id("data").Dot(f.NameGo())
	}
}

func (f *Slice) convTo(ctx Context) jen.Code {
	if f.source.IsNode || f.source.IsEdge {
		return nil
	} else if f.source.IsEnum {
		return jen.Id("convertEnum").Types(jen.String(), jen.Qual(ctx.SourcePkg, f.source.Value)).Call(jen.Id("data").Dot(f.NameGo()))
	} else {
		return jen.Id("data").Dot(f.NameGo())
	}
}

func (f *Slice) fieldDef(ctx Context) jen.Code {
	typ := f.source.Value
	if f.source.IsNode {
		typ = f.source.Value + "Field"
	}
	if f.source.IsEnum {
		typ = "string"
	}
	return jen.Id(f.NameGo()).Index().Id(typ).
		Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
}
