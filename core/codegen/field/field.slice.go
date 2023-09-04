package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/parser"
)

type Slice struct {
	*baseField

	source  *parser.FieldSlice
	element Field
}

func (f *Slice) typeGo() jen.Code {
	return jen.Add(f.ptr()).Index().Add(f.element.typeGo())
}

func (f *Slice) typeConv() jen.Code {
	return jen.Add(f.ptr()).Index().Add(f.element.typeConv())
	// return jen.Add(f.ptr()).Id("jsonArray").Types(f.element.typeCo
}

func (f *Slice) TypeDatabase() string {
	if f.element.TypeDatabase() == "" {
		return "" // TODO: this seems invalid, no?
	}

	// Go treats empty slices as nil, so the database needs
	// to accept the json NULL value for any array field.
	return "option<array | null>"
}

func (f *Slice) Element() Field {
	return f.element
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

	case *Node:
		{
			return jen.Func().
				Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).Id(f.NameGo()).
				Params(
					jen.Id("filters").Op("...").Qual(ctx.pkgLib(), "Filter").
						Types(jen.Qual(f.SourcePkg, element.table.NameGo())),
				).
				Id(element.table.NameGoLower()+"Slice").Types(jen.Id("T")).
				Block(
					jen.Id("key").Op(":=").Qual(ctx.pkgLib(), "Node").
						Call(
							jen.Id("n").Dot("key"),
							jen.Lit(f.NameDatabase()),
							jen.Id("filters"),
						),
					jen.Return(
						jen.Id(element.table.NameGoLower()+"Slice").Types(jen.Id("T")).
							Values(
								jen.Qual(ctx.pkgLib(), "KeyFilter").Types(jen.Id("T")).
									Call(jen.Id("key")),
								jen.Qual(ctx.pkgLib(), "NewSlice").Types(jen.Id("T"), jen.Qual(ctx.SourcePkg, element.table.NameGo())).
									Call(jen.Id("key")),
							),
					),
				)
		}

	case *Edge:
		{
			receiver := jen.Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))
			if ctx.Receiver != nil {
				receiver = ctx.Receiver
			}

			if tableEqual(ctx.Table, element.table.In.table) {
				return jen.Func().
					Params(jen.Id("n").Add(receiver)).Id(f.NameGo()).
					Params(
						jen.Id("filters").Op("...").Qual(ctx.pkgLib(), "Filter").
							Types(jen.Qual(f.SourcePkg, element.table.NameGo())),
					).
					Params(jen.Id(element.table.NameGoLower() + "In").Index(jen.Id("T"))).
					Block(
						jen.Return(
							jen.Id("new" + element.table.NameGo() + "In").Index(jen.Id("T")).
								Call(
									jen.Qual(ctx.pkgLib(), "EdgeIn").Call(
										jen.Id("n").Dot("key"),
										jen.Lit(element.table.NameDatabase()),
										jen.Id("filters"),
									),
								),
						),
					)
			}

			if tableEqual(ctx.Table, element.table.Out.table) {
				return jen.Func().
					Params(jen.Id("n").Add(receiver)).Id(f.NameGo()).
					Params(
						jen.Id("filters").Op("...").Qual(ctx.pkgLib(), "Filter").
							Types(jen.Qual(f.SourcePkg, element.table.NameGo())),
					).
					Params(jen.Id(element.table.NameGoLower() + "Out").Index(jen.Id("T"))).
					Block(
						jen.Return(
							jen.Id("new" + element.table.NameGo() + "Out").Index(jen.Id("T")).
								Call(
									jen.Qual(ctx.pkgLib(), "EdgeOut").Call(
										jen.Id("n").Dot("key"),
										jen.Lit(element.table.NameDatabase()),
										jen.Id("filters"),
									),
								),
						),
					)
			}

			return nil
		}

	case *Enum:
		{
			return jen.Func().
				Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
				Id(f.NameGo()).Params().
				Op("*").Qual(ctx.pkgLib(), "Slice").Types(jen.Id("T"), jen.Qual(ctx.SourcePkg, element.model.NameGo())).
				Block(
					jen.Return(
						jen.Qual(ctx.pkgLib(), "NewSlice").Types(jen.Id("T"), jen.Qual(ctx.SourcePkg, element.model.NameGo())).
							Call(
								jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())),
							),
					),
				)
		}

	default:
		{
			return jen.Func().
				Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
				Id(f.NameGo()).Params().
				Op("*").Qual(ctx.pkgLib(), "Slice").Types(jen.Id("T"), element.typeGo()).
				Block(
					jen.Return(
						jen.Qual(ctx.pkgLib(), "NewSlice").Types(jen.Id("T"), element.typeGo()).
							Call(
								jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())),
							),
					),
				)
		}
	}
}

func (f *Slice) convFrom(ctx Context) jen.Code {
	switch element := f.element.(type) {

	case *Node:
		{
			mapperFunc := "mapSlice"
			mapFunc := "to" + element.table.NameGo() + "Link"

			if f.source.Pointer() {
				mapperFunc += "Ptr"
			}

			if element.source.Pointer() {
				mapFunc += "Ptr"
			}

			return jen.Id(mapperFunc).Call(
				jen.Id("data").Dot(f.NameGo()),
				jen.Id(mapFunc),
			)
		}

	case *Struct:
		{
			mapFn := "mapSlice"
			fromFn := jen.Id("from" + element.table.NameGo())

			if f.source.Pointer() {
				mapFn = "mapSlicePtr"
			}

			if !element.source.Pointer() {
				fromFn = jen.Id("noPtrFunc").Call(fromFn)
			}

			return jen.Id(mapFn).Call(
				jen.Id("data").Dot(f.NameGo()),
				fromFn,
			)
		}

	case *Edge:
		{
			return nil // TODO: should an edge really not be addable like that?
		}

	case *Enum:
		{
			return jen.Id("data").Dot(f.NameGo())
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
		{
			mapperFunc := "mapSlice"
			mapFunc := "from" + element.table.NameGo() + "Link"

			if f.source.Pointer() {
				mapperFunc += "Ptr"
			}

			if element.source.Pointer() {
				mapFunc += "Ptr"
			}

			return jen.Id(mapperFunc).Call(
				jen.Id("data").Dot(f.NameGo()),
				jen.Id(mapFunc),
			)
		}

	case *Struct:
		{
			mapFn := "mapSlice"
			toFn := jen.Id("to" + element.table.NameGo())

			if f.source.Pointer() {
				mapFn = "mapSlicePtr"
			}

			if !element.source.Pointer() {
				toFn = jen.Id("noPtrFunc").Call(toFn)
			}

			return jen.Id(mapFn).Call(
				jen.Id("data").Dot(f.NameGo()),
				toFn,
			)
		}

	case *Edge:
		return jen.Id("mapSlice").Call(
			jen.Id("data").Dot(f.NameGo()),
			jen.Id("noPtrFunc").Call(jen.Id("To"+element.table.NameGo())),
		)

	case *Enum:
		{
			return jen.Id("data").Dot(f.NameGo())
		}

	default:
		{
			return jen.Id("data").Dot(f.NameGo())
		}

	}
}

func (f *Slice) fieldDef(ctx Context) jen.Code {
	jsonSuffix := ""
	if _, isEdge := f.element.(*Edge); isEdge {
		jsonSuffix = ",omitempty"
	}

	return jen.Id(f.NameGo()).Add(f.typeConv()).
		Tag(map[string]string{"json": f.NameDatabase() + jsonSuffix})
}
