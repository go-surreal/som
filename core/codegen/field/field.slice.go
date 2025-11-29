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
		return "option<bytes>"
	}

	// Go treats empty slices as nil, but the custom marshaling
	// ensures that they are stored as NONE in the database.
	return fmt.Sprintf("option<array<%s>>", f.element.TypeDatabase())
}

func (f *Slice) SchemaStatements(table, prefix string) []string {
	fieldType := f.TypeDatabase()
	if fieldType == "" {
		return nil
	}

	// Generate own DEFINE FIELD statement.
	statements := []string{
		fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, f.TypeDatabase(),
		),
	}

	// Only recurse into struct elements, because for primitive
	// elements the type is already part of the array definition.
	if structElem, ok := f.element.(*Struct); ok {
		nestedPrefix := prefix + f.NameDatabase() + ".*."
		for _, field := range structElem.Table().GetFields() {
			statements = append(statements, field.SchemaStatements(table, nestedPrefix)...)
		}
	}

	return statements
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
		sortFunc:   nil, // sorting by slices has no real use case

		cborMarshal:   f.cborMarshal,
		cborUnmarshal: f.cborUnmarshal,
	}
}

func (f *Slice) filterDefine(ctx Context) jen.Code {
	filter := "Slice"

	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	elemFilter := f.element.CodeGen().filterDefine.Exec(ctx.fromSlice())

	switch element := f.element.(type) {

	case *Node, *Edge:
		{
			if !ctx.isFromSlice {
				return nil // handled by filterFunc
			}
		}

	case *String:
		{
			filter := "String"

			if element.source.Pointer() {
				filter += fnSuffixPtr
			}

			filter += "Slice"

			if f.source.Pointer() {
				filter += fnSuffixPtr
			}

			return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).
				Types(def.TypeModel)
		}

	case *Numeric:
		{
			filter := "Numeric"

			switch element.source.Type {

			case parser.NumberInt, parser.NumberInt8, parser.NumberInt16, parser.NumberInt32, parser.NumberInt64,
				parser.NumberUint8, parser.NumberUint16, parser.NumberUint32, parser.NumberRune:
				{
					filter = "Int"
				}

			case parser.NumberFloat32, parser.NumberFloat64:
				{
					filter = "Float"
				}
			}

			if element.source.Pointer() {
				filter += fnSuffixPtr
			}

			filter += "Slice"

			if f.source.Pointer() {
				filter += fnSuffixPtr
			}

			return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).
				Types(def.TypeModel, element.typeGo())
		}

	case *Byte:
		{
			// TODO: pointers
			return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), "ByteSlice").Types(def.TypeModel)
		}

	case *Enum:
		return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(
			def.TypeModel,
			jen.Qual(ctx.SourcePkg, element.model.NameGo()),
			elemFilter,
		)
	}

	return jen.Id(f.NameGo()).Op("*").Qual(ctx.pkgLib(), filter).Types(
		def.TypeModel,
		f.element.typeGo(),
		elemFilter,
	)
}

func (f *Slice) filterInit(ctx Context) (jen.Code, jen.Code) {
	filter := "NewSlice"

	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

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
		}

	case *String:
		{
			filter := "NewString"

			if element.source.Pointer() {
				filter += fnSuffixPtr
			}

			filter += "Slice"

			if f.source.Pointer() {
				filter += fnSuffixPtr
			}

			return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel),
				jen.Call(
					jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
				)
		}

	case *Numeric:
		{
			filter := "NewNumericSlice"

			switch element.source.Type {

			case parser.NumberInt, parser.NumberInt8, parser.NumberInt16, parser.NumberInt32, parser.NumberInt64,
				parser.NumberUint8, parser.NumberUint16, parser.NumberUint32, parser.NumberRune:
				{
					filter = "NewInt"
				}

			case parser.NumberFloat32, parser.NumberFloat64:
				{
					filter = "NewFloat"
				}
			}

			if element.source.Pointer() {
				filter += fnSuffixPtr
			}

			filter += "Slice"

			if f.source.Pointer() {
				filter += fnSuffixPtr
			}

			return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel, element.typeGo()),
				jen.Call(
					jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
				)
		}

	case *Struct:
		{
			//if !ctx.isFromSlice {
			//	return nil, nil // handled by filterFunc
			//}
		}

	case *Byte:
		return jen.Qual(ctx.pkgLib(), "NewByteSlice").Types(def.TypeModel),
			jen.Call(
				jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
			)

	case *Enum:
		return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel, jen.Qual(ctx.SourcePkg, element.model.NameGo())),
			jen.Call(
				jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
				makeElemFilter,
			)
	}

	filter = "NewSliceMaker"

	if f.source.Pointer() {
		filter += fnSuffixPtr
	}

	return jen.Qual(ctx.pkgLib(), filter).Types(def.TypeModel, f.element.typeGo(), elemFilter).
			Call(makeElemFilter),
		jen.Call(
			jen.Qual(ctx.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(f.NameDatabase())),
		)
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
							jen.Id("n").Dot("Key"),
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
										jen.Id("n").Dot("Key"),
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
										jen.Id("n").Dot("Key"),
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

func (f *Slice) cborMarshal(_ Context) jen.Code {
	// For struct slices, we need to convert each element through the conv wrapper
	// to get proper snake_case field names in the CBOR output.
	if structElem, ok := f.element.(*Struct); ok {
		convFuncName := "from" + structElem.element.NameGo()
		if structElem.source.Pointer() {
			convFuncName += "Ptr"
		}

		// Determine the slice element type based on whether the struct element is a pointer
		var sliceElemType jen.Code
		if structElem.source.Pointer() {
			sliceElemType = jen.Index().Op("*").Id(structElem.element.NameGoLower())
		} else {
			sliceElemType = jen.Index().Id(structElem.element.NameGoLower())
		}

		// Handle pointer-to-slice case by dereferencing
		srcSlice := jen.Id("c").Dot(f.NameGo())
		if f.source.Pointer() {
			srcSlice = jen.Op("*").Id("c").Dot(f.NameGo())
		}

		return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
			jen.Id("convSlice").Op(":=").Make(
				sliceElemType,
				jen.Len(srcSlice),
			),
			jen.For(
				jen.Id("i").Op(",").Id("v").Op(":=").Range().Add(srcSlice),
			).Block(
				jen.Id("convSlice").Index(jen.Id("i")).Op("=").Id(convFuncName).Call(jen.Id("v")),
			),
			jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("convSlice"),
		)
	}

	return jen.If(jen.Id("c").Dot(f.NameGo()).Op("!=").Nil()).Block(
		jen.Id("data").Index(jen.Lit(f.NameDatabase())).Op("=").Id("c").Dot(f.NameGo()),
	)
}

func (f *Slice) cborUnmarshal(ctx Context) jen.Code {
	// For struct slices, we need to unmarshal into the conv wrapper and then convert back.
	if structElem, ok := f.element.(*Struct); ok {
		convFuncName := "to" + structElem.element.NameGo()
		if structElem.source.Pointer() {
			convFuncName += "Ptr"
		}

		// Determine the slice element type based on whether the struct element is a pointer
		var sliceElemType jen.Code
		if structElem.source.Pointer() {
			sliceElemType = jen.Index().Op("*").Id(structElem.element.NameGoLower())
		} else {
			sliceElemType = jen.Index().Id(structElem.element.NameGoLower())
		}

		// Determine the inner slice type (the model slice type)
		innerSliceType := jen.Index().Add(f.element.typeGo())

		// Build the assignment statement - handle pointer-to-slice case
		var assignStmt jen.Code
		if f.source.Pointer() {
			assignStmt = jen.Block(
				jen.Id("result").Op(":=").Make(innerSliceType, jen.Len(jen.Id("convSlice"))),
				jen.For(
					jen.Id("i").Op(",").Id("v").Op(":=").Range().Id("convSlice"),
				).Block(
					jen.Id("result").Index(jen.Id("i")).Op("=").Id(convFuncName).Call(jen.Id("v")),
				),
				jen.Id("c").Dot(f.NameGo()).Op("=").Op("&").Id("result"),
			)
		} else {
			assignStmt = jen.Block(
				jen.Id("c").Dot(f.NameGo()).Op("=").Make(
					f.typeGo(),
					jen.Len(jen.Id("convSlice")),
				),
				jen.For(
					jen.Id("i").Op(",").Id("v").Op(":=").Range().Id("convSlice"),
				).Block(
					jen.Id("c").Dot(f.NameGo()).Index(jen.Id("i")).Op("=").Id(convFuncName).Call(jen.Id("v")),
				),
			)
		}

		return jen.If(
			jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
			jen.Id("ok"),
		).BlockFunc(func(g *jen.Group) {
			g.Var().Id("convSlice").Add(sliceElemType)
			g.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convSlice"))
			g.Add(assignStmt)
		})
	}

	return jen.If(
		jen.Id("raw").Op(",").Id("ok").Op(":=").Id("rawMap").Index(jen.Lit(f.NameDatabase())),
		jen.Id("ok"),
	).Block(
		jen.Qual(ctx.pkgCBOR(), "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("c").Dot(f.NameGo())),
	)
}
