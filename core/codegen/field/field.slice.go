package field

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
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

func (f *Slice) typeConv(ctx Context) jen.Code {
	return jen.Add(f.ptr()).Index().Add(f.element.typeConv(ctx))
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
	elemFilter := f.element.CodeGen().filterDefine.Exec(ctx.fromSlice())

	switch element := f.element.(type) {

	case *Node, *Edge:
		{
			if !ctx.isFromSlice {
				return nil // handled by filterFunc
			}

			return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "Slice").Types(
				def.TypeModel,
				element.typeGo(), // TODO: no pointers!
				elemFilter,
			)
		}

	case *Byte:
		return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "ByteSlice").Types(def.TypeModel)

	case *Enum:
		return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "Slice").Types(
			def.TypeModel,
			jen.Qual(ctx.SourcePkg, element.model.NameGo()),
			elemFilter,
		)

	default:
		return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "Slice").Types(
			def.TypeModel,
			element.typeGo(), // TODO: no pointers!
			elemFilter,
		)
	}
}

func (f *Slice) filterInit(ctx Context) (jen.Code, jen.Code) {
	elemFilter := f.element.CodeGen().filterDefine.Exec(ctx.fromSlice())

	var makeElemFilter jen.Code
	if f.element.CodeGen().filterInit != nil {
		makeElemFilter, _ = f.element.CodeGen().filterInit(ctx.fromSlice())
	}

	if makeElemFilter == nil {
		fmt.Printf("no filter init for %T\n", f.element)
	}

	switch element := f.element.(type) {

	case *Node, *Edge:
		{
			if !ctx.isFromSlice {
				return nil, nil // handled by filterFunc
			}

			return jen.Qual(ctx.pkgLib(), "NewSliceMaker").Types(def.TypeModel, element.typeGo(), elemFilter).
					Call(makeElemFilter),
				jen.Call(
					jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
				)
		}

	case *Struct:
		{
			//if !ctx.isFromSlice {
			//	return nil, nil // handled by filterFunc
			//}

			return jen.Qual(ctx.pkgLib(), "NewSliceMaker").Types(def.TypeModel, element.typeGo(), elemFilter).
					Call(makeElemFilter),
				jen.Call(
					jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
				)
		}

	case *Byte:
		return jen.Qual(ctx.pkgLib(), "NewByteSlice").Types(def.TypeModel),
			jen.Call(
				jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
			)

	case *Enum:
		return jen.Qual(ctx.pkgLib(), "NewSlice").Types(def.TypeModel, jen.Qual(ctx.SourcePkg, element.model.NameGo())),
			jen.Call(
				jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
				makeElemFilter,
			)

	default:
		return jen.Qual(ctx.pkgLib(), "NewSliceMaker").Types(def.TypeModel, element.typeGo(), elemFilter).
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
				Params(jen.Id("n").Id(ctx.Table.NameGoLower()).Types(def.TypeModel)).Id(f.NameGo()).
				Params(
					jen.Id("filters").Op("...").Qual(ctx.pkgLib(), "Filter").
						Types(jen.Qual(f.SourcePkg, element.table.NameGo())),
				).
				Op("*").Qual(ctx.pkgLib(), "Slice").
				Types(
					def.TypeModel, jen.Qual(f.SourcePkg, element.table.NameGo()), jen.Id(element.table.NameGoLower()).Types(def.TypeModel),
				).
				Block(
					jen.Id("key").Op(":=").Qual(ctx.pkgLib(), "Node").
						Call(
							jen.Id("n").Dot("key"),
							jen.Lit(f.NameDatabase()),
							jen.Id("filters"),
						),
					jen.Return(
						jen.Qual(ctx.pkgLib(), "NewSlice").
							Types(
								def.TypeModel,
								jen.Qual(ctx.SourcePkg, element.table.NameGo()),
								elemFilter,
							).
							Call(
								jen.Id("key"),
								makeElemFilter,
							),
					),
				)
		}

	case *Edge:
		{
			receiver := jen.Id(ctx.Table.NameGoLower()).Types(def.TypeModel)
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
					Params(jen.Id(element.table.NameGoLower() + "In").Index(def.TypeModel)).
					Block(
						jen.Return(
							jen.Id("new" + element.table.NameGo() + "In").Index(def.TypeModel).
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
					Params(jen.Id(element.table.NameGoLower() + "Out").Index(def.TypeModel)).
					Block(
						jen.Return(
							jen.Id("new" + element.table.NameGo() + "Out").Index(def.TypeModel).
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

	default:
		return nil
	}
}

func (f *Slice) convFrom(ctx Context) (jen.Code, jen.Code) {
	switch element := f.element.(type) {

	case *Slice:
		{
			fromFunc, _ := element.CodeGen().convFrom(ctx.fromSlice())

			if fromFunc == nil || isCodeEqual(fromFunc, jen.Null()) {
				return jen.Null(), jen.Id("data").Dot(f.NameGo())
			}

			mapperFunc := "mapSliceFn"

			if f.source.Pointer() {
				mapperFunc += fnSuffixPtr
			}

			return jen.Id(mapperFunc).Call(fromFunc),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Node:
		{
			mapperFunc := "mapSliceFn"
			fromFunc := "to" + element.table.NameGo() + "Link"

			if f.source.Pointer() {
				mapperFunc += fnSuffixPtr
			}

			if element.source.Pointer() {
				fromFunc += fnSuffixPtr
			}

			return jen.Id(mapperFunc).Call(jen.Id(fromFunc)),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Struct:
		{
			mapperFunc := "mapSliceFn"
			fromFunc := jen.Id("from" + element.table.NameGo())

			if f.source.Pointer() {
				mapperFunc += fnSuffixPtr
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
			fromFunc, _ := element.CodeGen().convFrom(ctx.fromSlice())

			if ctx.isFromSlice && isCodeEqual(fromFunc, jen.Null()) {
				return nil, nil // native types do not need conversion
			}

			if fromFunc != nil && !isCodeEqual(fromFunc, jen.Null()) {
				mapperFunc := "mapSliceFn"

				if f.source.Pointer() {
					mapperFunc += fnSuffixPtr
				}

				return jen.Id(mapperFunc).Call(fromFunc), jen.Call(jen.Id("data").Dot(f.NameGo()))
			}

			return jen.Null(), jen.Id("data").Dot(f.NameGo())
		}
	}
}

func (f *Slice) convTo(ctx Context) (jen.Code, jen.Code) {
	switch element := f.element.(type) {

	case *Slice:
		{
			toFunc, _ := element.CodeGen().convTo(ctx.fromSlice())

			if toFunc == nil || isCodeEqual(toFunc, jen.Null()) {
				return jen.Null(), jen.Id("data").Dot(f.NameGo())
			}

			mapperFunc := "mapSliceFn"

			if f.source.Pointer() {
				mapperFunc += fnSuffixPtr
			}

			return jen.Id(mapperFunc).Call(toFunc),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Node:
		{
			mapperFunc := "mapSliceFn"
			toFunc := "from" + element.table.NameGo() + "Link"

			if f.source.Pointer() {
				mapperFunc += fnSuffixPtr
			}

			if element.source.Pointer() {
				toFunc += fnSuffixPtr
			}

			return jen.Id(mapperFunc).Call(jen.Id(toFunc)),
				jen.Call(jen.Id("data").Dot(f.NameGo()))
		}

	case *Struct:
		{
			mapperFunc := "mapSliceFn"
			toFunc := jen.Id("to" + element.table.NameGo())

			if f.source.Pointer() {
				mapperFunc += fnSuffixPtr
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
				mapperFunc += fnSuffixPtr
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
			toFunc, _ := element.CodeGen().convTo(ctx.fromSlice())

			if ctx.isFromSlice && isCodeEqual(toFunc, jen.Null()) {
				return nil, nil // native types do not need conversion
			}

			if toFunc != nil && !isCodeEqual(toFunc, jen.Null()) {
				mapperFunc := "mapSliceFn"

				if f.source.Pointer() {
					mapperFunc += fnSuffixPtr
				}

				return jen.Id(mapperFunc).Call(toFunc), jen.Call(jen.Id("data").Dot(f.NameGo()))
			}

			return jen.Null(), jen.Id("data").Dot(f.NameGo())
		}

	}
}

func (f *Slice) fieldDef(ctx Context) jen.Code {
	jsonSuffix := ""
	if _, isEdge := f.element.(*Edge); isEdge {
		jsonSuffix = ",omitempty"
	}

	return jen.Id(f.NameGo()).Add(f.typeConv(ctx)).
		Tag(map[string]string{convTag: f.NameDatabase() + jsonSuffix})
}
