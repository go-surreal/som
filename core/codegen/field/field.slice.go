package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
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
}

func (f *Slice) TypeDatabase() string {
	if f.element.TypeDatabase() == "" {
		return ""
	}

	// Note: No "ASSERT $value != NULL" used here,
	// because the zero value of a slice is nil.
	return "array"
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
					jen.Id("filters").Op("...").Qual(def.PkgLib, "Filter").
						Types(jen.Qual(f.SourcePkg, element.table.NameGo())),
				).
				Id(element.table.NameGoLower()+"Slice").Types(jen.Id("T")).
				Block(
					jen.Id("key").Op(":=").Qual(def.PkgLib, "Node").
						Call(
							jen.Id("n").Dot("key"),
							jen.Lit(f.NameDatabase()),
							jen.Id("filters"),
						),
					jen.Return(
						jen.Id(element.table.NameGoLower()+"Slice").Types(jen.Id("T")).
							Values(
								jen.Qual(def.PkgLib, "KeyFilter").Types(jen.Id("T")).
									Call(jen.Id("key")),
								jen.Qual(def.PkgLib, "NewSlice").Types(jen.Id("T"), jen.Qual(ctx.SourcePkg, element.table.NameGo())).
									Call(jen.Id("key")),
							),
					),
				)
		}

	case *Edge:
		{
			if tableEqual(ctx.Table, element.table.In.table) {
				return jen.Func().
					Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).Id(f.NameGo()).
					Params(
						jen.Id("filters").Op("...").Qual(def.PkgLib, "Filter").
							Types(jen.Qual(f.SourcePkg, element.table.NameGo())),
					).
					Params(jen.Id(element.table.NameGoLower() + "In").Index(jen.Id("T"))).
					Block(
						jen.Return(
							jen.Id("new" + element.table.NameGo() + "In").Index(jen.Id("T")).
								Call(
									jen.Qual(def.PkgLib, "EdgeIn").Call(
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
					Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).Id(f.NameGo()).
					Params(
						jen.Id("filters").Op("...").Qual(def.PkgLib, "Filter").
							Types(jen.Qual(f.SourcePkg, element.table.NameGo())),
					).
					Params(jen.Id(element.table.NameGoLower() + "Out").Index(jen.Id("T"))).
					Block(
						jen.Return(
							jen.Id("new" + element.table.NameGo() + "Out").Index(jen.Id("T")).
								Call(
									jen.Qual(def.PkgLib, "EdgeOut").Call(
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
				Op("*").Qual(def.PkgLib, "Slice").Types(jen.Id("T"), jen.Qual(ctx.SourcePkg, element.model.NameGo())).
				Block(
					jen.Return(
						jen.Qual(def.PkgLib, "NewSlice").Types(jen.Id("T"), jen.Qual(ctx.SourcePkg, element.model.NameGo())).
							Call(
								jen.Qual(def.PkgLib, "Field").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())),
							),
					),
				)
		}

	default:
		{
			return jen.Func().
				Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
				Id(f.NameGo()).Params().
				Op("*").Qual(def.PkgLib, "Slice").Types(jen.Id("T"), element.typeGo()).
				Block(
					jen.Return(
						jen.Qual(def.PkgLib, "NewSlice").Types(jen.Id("T"), element.typeGo()).
							Call(
								jen.Qual(def.PkgLib, "Field").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())),
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
			mapFn := "mapSlice"
			if f.source.Pointer() {
				mapFn = "mapSlicePtr"
			}

			if element.source.Pointer() {
				mapFn = "mapPtrSlice"
				if f.source.Pointer() {
					mapFn = "mapPtrSlicePtr"
				}
			}

			return jen.Id(mapFn).Call(
				jen.Id("data").Dot(f.NameGo()),
				jen.Id("to"+element.table.NameGo()+"Link"),
			)
		}

	case *Struct:
		{
			mapFn := "mapSlice"
			if f.source.Pointer() {
				mapFn = "mapSlicePtr"
			}

			if element.source.Pointer() {
				mapFn = "mapPtrSlice"
				if f.source.Pointer() {
					mapFn = "mapPtrSlicePtr"
				}
			}

			return jen.Id(mapFn).Call(
				jen.Id("data").Dot(f.NameGo()),
				jen.Id("from"+element.table.NameGo()),
			)
		}

	case *Edge:
		{
			return nil // TODO: should an edge really not be addable like that?
		}

	case *Enum:
		{
			mapEnumFn := jen.Id("mapEnum").Types(jen.Qual(f.SourcePkg, element.model.NameGo()), jen.String())
			if element.source.Pointer() {
				mapEnumFn = jen.Id("ptrFunc").Call(mapEnumFn)
			}

			return jen.Id("mapSlice").Call(jen.Id("data").Dot(f.NameGo()), mapEnumFn)
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
			mapFn := "mapSlice"
			if f.source.Pointer() {
				mapFn = "mapSlicePtr"
			}

			if element.source.Pointer() {
				mapFn = "mapPtrSlice"
				if f.source.Pointer() {
					mapFn = "mapPtrSlicePtr"
				}
			}

			return jen.Id(mapFn).Call(
				jen.Id("data").Dot(f.NameGo()),
				jen.Id("from"+element.table.NameGo()+"Link"),
			)
		}

	case *Struct:
		{
			mapFn := "mapSlice"
			if f.source.Pointer() {
				mapFn = "mapSlicePtr"
			}

			if element.source.Pointer() {
				mapFn = "mapPtrSlice"
				if f.source.Pointer() {
					mapFn = "mapPtrSlicePtr"
				}
			}

			return jen.Id(mapFn).Call(
				jen.Id("data").Dot(f.NameGo()),
				jen.Id("to"+element.table.NameGo()),
			)
		}

	case *Edge:
		return jen.Id("mapSlice").Call(jen.Id("data").Dot(f.NameGo()), jen.Id("To"+element.table.NameGo()))

	case *Enum:
		{
			mapEnumFn := jen.Id("mapEnum").Types(jen.String(), jen.Qual(f.SourcePkg, element.model.NameGo()))
			if element.source.Pointer() {
				mapEnumFn = jen.Id("ptrFunc").Call(mapEnumFn)
			}

			return jen.Id("mapSlice").Call(jen.Id("data").Dot(f.NameGo()), mapEnumFn)
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
