package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type Slice struct {
	*baseField

	source  *parser.FieldSlice
	element Field
}

func (f *Slice) typeGo() jen.Code {
	return jen.Index().Add(f.element.typeGo())
}

// TODO: cleanup this (temporary) method!
func (f *Slice) Edge() (bool, Field, *Node, *Node) {
	edgeElement, ok := f.element.(*Edge)
	if !ok {
		return false, nil, nil, nil
	}

	in := edgeElement.table.In
	out := edgeElement.table.Out

	return true, f.element, &in, &out
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
	switch element := f.element.(type) {

	case *Edge:
		{
			if tableEqual(ctx.Table, &element.table.In.table) {
				return jen.Func().
					Params(jen.Id("n").Id(strcase.ToLowerCamel(ctx.Table.NameGo())).Types(jen.Id("T"))).
					Id(f.NameGo()).Params().
					Params(jen.Id(element.NameGoLower() + "In").Index(jen.Id("T"))).
					Block(
						jen.Return(
							jen.Id("new" + element.NameGo() + "In").Index(jen.Id("T")).
								Call(jen.Id("n").Dot("key").Dot("In").Call(jen.Lit(element.NameDatabase()))),
						),
					)
			}

			if tableEqual(ctx.Table, &element.table.Out.table) {
				return jen.Func().
					Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
					Id(f.NameGo()).Params().
					Params(jen.Id(element.NameGoLower() + "Out").Index(jen.Id("T"))).
					Block(
						jen.Return(
							jen.Id("new" + element.NameGo() + "Out").Index(jen.Id("T")).
								Call(jen.Id("n").Dot("key").Dot("Out").Call(jen.Lit(element.NameDatabase()))),
						),
					)
			}

			return nil
		}

	case *Node:
		{
			return jen.Func().
				Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
				Id(f.NameGo()).Params().
				Id(element.table.NameGoLower()+"Slice").Types(jen.Id("T")).
				Block(
					jen.Id("key").Op(":=").Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))),
					jen.Return(
						jen.Id(element.table.NameGoLower()+"Slice").Types(jen.Id("T")).
							Values(
								jen.Id("new"+element.table.NameGo()).Types(jen.Id("T")).
									Call(jen.Id("key")),
								jen.Qual(def.PkgLibFilter, "NewSlice").Types(jen.Qual(ctx.SourcePkg, element.table.NameGo()), jen.Id("T")).
									Call(jen.Id("key")),
							),
					),
				)
		}

	case *Enum:
		{
			return jen.Func().
				Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
				Id(f.NameGo()).Params().
				Op("*").Qual(def.PkgLibFilter, "Slice").Types(jen.Qual(ctx.SourcePkg, element.model.NameGo()), jen.Id("T")).
				Block(
					jen.Return(
						jen.Qual(def.PkgLibFilter, "NewSlice").Types(jen.Qual(ctx.SourcePkg, element.model.NameGo()), jen.Id("T")).
							Call(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo())))),
					),
				)
		}

	default:
		{
			return jen.Func().
				Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
				Id(f.NameGo()).Params().
				Op("*").Qual(def.PkgLibFilter, "Slice").Types(element.typeGo(), jen.Id("T")).
				Block(
					jen.Return(
						jen.Qual(def.PkgLibFilter, "NewSlice").Types(element.typeGo(), jen.Id("T")).
							Call(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo())))),
					),
				)
		}
	}
}

func (f *Slice) convFrom(ctx Context) jen.Code {
	switch element := f.element.(type) {

	case *Node:
		{
			return jen.Id("mapRecords").Call(jen.Id("data").Dot(f.NameGo()), jen.Id("to"+element.NameGo()+"Field"))
		}

	case *Edge:
		{
			return nil // TODO: should edges really not be addable like that?
		}

	case *Enum:
		{
			return jen.Id("convertEnum").Types(element.typeGo(), jen.String()).Call(jen.Id("data").Dot(f.NameGo()))
		}

	default:
		{
			return jen.Id("data").Dot(f.NameGo())
		}

	}
}

func (f *Slice) convTo(ctx Context) jen.Code {
	switch element := f.element.(type) {

	case *Node:
		return nil

	case *Edge:
		return nil

	case *Enum:
		{
			return jen.Id("convertEnum").
				Types(jen.String(), element.typeGo()).
				Call(jen.Id("data").Dot(f.NameGo()))
		}

	default:
		{
			return jen.Id("data").Dot(f.NameGo())
		}

	}
}

func (f *Slice) fieldDef(ctx Context) jen.Code {
	switch element := f.element.(type) {

	case *Node:
		{
			return jen.Id(f.NameGo()).Index().Id(element.table.NameGo() + "Field").
				Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
		}

	case *Enum:
		{
			return jen.Id(f.NameGo()).Index().String().
				Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
		}

	default:
		{
			return jen.Id(f.NameGo()).Index().Add(element.typeGo()).
				Tag(map[string]string{"json": f.NameDatabase() + ",omitempty"})
		}
	}
}
