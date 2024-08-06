package field

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
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
		return "" // TODO: this is invalid, no?
	}

	if _, ok := f.element.(*Byte); ok {
		return "option<bytes | null>"
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
		filterDefine: f.filterDefine,
		filterInit:   f.filterInit,
		filterFunc:   f.filterFunc,

		sortDefine: nil,
		sortInit:   nil,
		sortFunc:   nil, // TODO

		convFrom: f.convFrom,
		convTo:   f.convTo,
		fieldDef: f.fieldDef,
	}
}

func (f *Slice) filterDefine(ctx Context) jen.Code {
	elemFilter := f.element.CodeGen().filterDefine.Exec(ctx)

	switch element := f.element.(type) {

	case *Node:
		return jen.Id(element.table.NameGoLower() + "Slice").Types(jen.Id("T"))

	case *Edge, *Struct:
		return nil // handled by filterFunc

	case *Byte:
		return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "ByteSlice").Types(jen.Id("T"))

	case *Enum:
		return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "Slice").Types(
			jen.Id("T"),
			jen.Qual(ctx.SourcePkg, element.model.NameGo()),
			elemFilter,
		)

	default:
		return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "Slice").Types(
			jen.Id("T"),
			element.typeGo(),
			elemFilter,
		)
	}
}

func (f *Slice) filterInit(ctx Context) (jen.Code, jen.Code) {
	elemFilter := f.element.CodeGen().filterDefine.Exec(ctx)

	var makeElemFilter jen.Code
	if f.element.CodeGen().filterInit != nil {
		makeElemFilter, _ = f.element.CodeGen().filterInit(ctx)
	}

	if makeElemFilter == nil {
		fmt.Printf("no filter init for %T\n", f.element)
	}

	switch element := f.element.(type) {

	case *Node, *Edge, *Struct:
		return nil, nil // handled by filterFunc

	case *Byte:
		return jen.Qual(ctx.pkgLib(), "NewByteSlice").Types(jen.Id("T")),
			jen.Call(
				jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
			)

	case *Enum:
		return jen.Qual(ctx.pkgLib(), "NewSlice").Types(jen.Id("T"), jen.Qual(ctx.SourcePkg, element.model.NameGo())),
			jen.Call(
				jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
				makeElemFilter,
			)

	case *Slice:
		return jen.Qual(ctx.pkgLib(), "NewSliceMaker").Types(jen.Id("T"), element.typeGo(), elemFilter).
				Call(makeElemFilter),
			jen.Call(
				jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
				//jen.Qual(ctx.pkgLib(), "NewSliceMaker").Types(jen.Id("T"), element.element.typeGo(), elemElemFilter).
				//	Call(makeElemElemFilter),
			)

	default:
		return jen.Qual(ctx.pkgLib(), "NewSliceMaker").Types(jen.Id("T"), element.typeGo(), elemFilter).
				Call(makeElemFilter),
			jen.Call(
				jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
			)
	}
}

func (f *Slice) filterFunc(ctx Context) jen.Code {
	elemFilter := f.element.CodeGen().filterDefine.Exec(ctx)

	var makeElemFilter jen.Code
	if f.element.CodeGen().filterInit != nil {
		makeElemFilter, _ = f.element.CodeGen().filterInit(ctx)
	} else {
		fmt.Printf("no filter init for %T\n", f.element)
	}

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
								jen.Qual(ctx.pkgLib(), "NewSlice").
									Types(
										jen.Id("T"),
										jen.Qual(ctx.SourcePkg, element.table.NameGo()),
										elemFilter,
									).
									Call(
										jen.Id("key"),
										makeElemFilter,
									),
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

	case *Struct:
		{
			return jen.Func().
				Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(jen.Id("T"))).
				Id(f.NameGo()).Params().
				Op("*").Qual(ctx.pkgLib(), "Slice").
				Types(jen.Id("T"), element.typeGo(), elemFilter).
				Block(
					jen.Return(
						jen.Qual(ctx.pkgLib(), "NewSlice").Types(jen.Id("T"), element.typeGo(), elemFilter).
							Call(
								jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("n").Dot("key"), jen.Lit(f.NameDatabase())),
								makeElemFilter,
							),
					),
				)
		}

	default:
		return nil
	}
}

func (f *Slice) convFrom(ctx Context) (jen.Code, jen.Code) {
	switch element := f.element.(type) {

	case *Slice:
		{
			mapperFunc := "mapSliceFn"
			fromFunc, _ := element.CodeGen().convFrom(ctx.fromSlice())

			if f.source.Pointer() {
				mapperFunc += "Ptr"
			}

			return jen.Id(mapperFunc).Call(fromFunc),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Node:
		{
			mapperFunc := "mapSliceFn"
			fromFunc := "to" + element.table.NameGo() + "Link"

			if f.source.Pointer() {
				mapperFunc += "Ptr"
			}

			if element.source.Pointer() {
				fromFunc += "Ptr"
			}

			return jen.Id(mapperFunc).Call(jen.Id(fromFunc)),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Struct:
		{
			mapperFunc := "mapSliceFn"
			fromFunc := jen.Id("from" + element.table.NameGo())

			if f.source.Pointer() {
				mapperFunc += "Ptr"
			}

			if !element.source.Pointer() {
				fromFunc = jen.Id("noPtrFunc").Call(fromFunc)
			}

			return jen.Id(mapperFunc).Call(fromFunc),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Edge:
		{
			return nil, nil // TODO: should an edge really not be addable like that?
		}

	case *Enum:
		{
			return nil, jen.Id("data").Dot(f.NameGo()) // TODO: correct?
		}

	default:
		{
			element.CodeGen().convFrom(ctx.fromSlice())

			if !ctx.isFromSlice {
				return jen.Null(), jen.Id("data").Dot(f.NameGo())
			}

			typ := jen.Index().Add(element.typeGo())

			if f.source.Pointer() {
				typ = jen.Op("*").Add(typ)
			}

			return jen.Id("noOp").Types(typ), nil // TODO: correct?
		}

	}
}

func (f *Slice) convTo(ctx Context) (jen.Code, jen.Code) {
	switch element := f.element.(type) {

	case *Slice:
		{
			mapperFunc := "mapSliceFn"
			toFunc, _ := element.CodeGen().convTo(ctx.fromSlice())

			if f.source.Pointer() {
				mapperFunc += "Ptr"
			}

			return jen.Id(mapperFunc).Call(toFunc),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Node:
		{
			mapperFunc := "mapSliceFn"
			toFunc := "from" + element.table.NameGo() + "Link"

			if f.source.Pointer() {
				mapperFunc += "Ptr"
			}

			if element.source.Pointer() {
				toFunc += "Ptr"
			}

			return jen.Id(mapperFunc).Call(jen.Id(toFunc)),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Struct:
		{
			mapperFunc := "mapSliceFn"
			toFunc := jen.Id("to" + element.table.NameGo())

			if f.source.Pointer() {
				mapperFunc += "Ptr"
			}

			if !element.source.Pointer() {
				toFunc = jen.Id("noPtrFunc").Call(toFunc)
			}

			return jen.Id(mapperFunc).Call(toFunc),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Edge:
		{
			mapperFunc := "mapSliceFn"
			toFunc := jen.Id("noPtrFunc").Call(jen.Id("To" + element.table.NameGo()))

			if f.source.Pointer() {
				mapperFunc += "Ptr"
			}

			// TODO: Edge can be not a pointer, no?

			return jen.Id(mapperFunc).Call(toFunc),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Enum:
		{
			return nil, jen.Id("data").Dot(f.NameGo()) // TODO: correct?
		}

	default:
		{
			if !ctx.isFromSlice {
				return jen.Null(), jen.Id("data").Dot(f.NameGo())
			}

			typ := jen.Index().Add(element.typeGo())

			if f.source.Pointer() {
				typ = jen.Op("*").Add(typ)
			}

			return jen.Id("noOp").Types(typ), nil // TODO: correct?
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
