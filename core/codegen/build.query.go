package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
)

type queryBuilder struct {
	*baseBuilder
}

func newQueryBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *queryBuilder {
	return &queryBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *queryBuilder) build() error {
	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *queryBuilder) buildFile(node *field.NodeTable) error {
	pkgLib := b.relativePkgPath(def.PkgLib)
	pkgConv := b.relativePkgPath(def.PkgConv)
	somPkg := b.relativePkgPath()

	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	modelType := b.SourceQual(node.Name)

	modelInfoVarName := node.NameGoLower() + "ModelInfo"
	rangeFnVarName := node.NameGoLower() + "RangeFn"
	selectTypeName := node.NameGoLower() + "Select"

	convFn := jen.Qual(pkgConv, "To"+node.NameGo()+"Ptr")

	f.Line()
	f.Commentf("%s holds the model-specific unmarshal functions for %s.", modelInfoVarName, node.NameGo())
	f.Var().Id(modelInfoVarName).Op("=").Id("modelInfo").Types(modelType).Values(jen.Dict{
		jen.Id("UnmarshalAll"): jen.Func().Params(
			jen.Id("data").Index().Byte(),
		).Params(jen.Index().Op("*").Add(modelType), jen.Error()).Block(
			jen.Return(jen.Id("unmarshalAll").Call(jen.Id("data"), convFn)),
		),
		jen.Id("UnmarshalOne"): jen.Func().Params(
			jen.Id("data").Index().Byte(),
		).Params(jen.Op("*").Add(modelType), jen.Error()).Block(
			jen.Return(jen.Id("unmarshalOne").Call(jen.Id("data"), convFn)),
		),
		jen.Id("UnmarshalSearchAll"): jen.Func().Params(
			jen.Id("data").Index().Byte(),
			jen.Id("clauses").Index().Qual(pkgLib, "SearchClause"),
		).Params(jen.Index().Qual(pkgLib, "SearchResult").Types(jen.Op("*").Add(modelType)), jen.Error()).Block(
			jen.Return(jen.Id("unmarshalSearchAll").Call(jen.Id("data"), jen.Id("clauses"), convFn)),
		),
	})

	if node.HasComplexID() {
		b.generateRangeFn(f, node, pkgLib, somPkg, modelType, rangeFnVarName)
	} else if node.HasStringID() {
		b.generateStringIDRangeFn(f, node, pkgLib, somPkg, modelType, rangeFnVarName)
	}

	// Generate select struct and field methods
	b.generateSelectStruct(f, node, pkgLib, modelType, selectTypeName)

	// Generate exported type alias for use in repo interfaces
	queryAliasName := node.NameGo() + "Query"
	f.Line()
	f.Commentf("%s is a type alias for the %s query builder.", queryAliasName, node.NameGo())
	f.Type().Id(queryAliasName).Op("=").Id("Builder").Types(modelType, jen.Id(selectTypeName))

	f.Line()
	f.Commentf("New%s creates a new query builder for %s models.", node.NameGo(), node.NameGo())
	f.Func().Id("New"+node.Name).
		Params(
			jen.Id("db").Id("Database"),
		).
		Id("Builder").Types(modelType, jen.Id(selectTypeName)).
		BlockFunc(func(g *jen.Group) {
			g.Id("q").Op(":=").Qual(pkgLib, "NewQuery").Types(modelType).Call(jen.Lit(node.NameDatabase()))

			if node.Source.SoftDelete {
				g.Comment("Automatically exclude soft-deleted records")
				pkgFilter := path.Join(b.basePkg, def.PkgFilter)
				g.Id("q").Dot("SoftDeleteFilter").Op("=").
					Qual(pkgFilter, node.Name).Dot("DeletedAt").Dot("Nil").Call(jen.Lit(true))
			}

			builderDict := jen.Dict{
				jen.Id("db"):    jen.Id("db"),
				jen.Id("query"): jen.Id("q"),
				jen.Id("info"):  jen.Id(modelInfoVarName),
				jen.Id("selectFn"): jen.Func().Params(
					jen.Id("sc").Id("SelectContext"),
				).Id(selectTypeName).Block(
					jen.Return(jen.Id(selectTypeName).Values(jen.Dict{
						jen.Id("SelectContext"): jen.Id("sc"),
					})),
				),
			}

			if node.HasComplexID() || node.HasStringID() {
				builderDict[jen.Id("rangeFn")] = jen.Id(rangeFnVarName)
			}

			g.Return(
				jen.Id("Builder").Types(modelType, jen.Id(selectTypeName)).
					Values(
						jen.Id("builder").Types(modelType, jen.Id(selectTypeName)).
							Values(builderDict),
					),
			)
		})

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}

func (b *queryBuilder) generateSelectStruct(
	f *jen.File, node *field.NodeTable,
	pkgLib string, modelType jen.Code, selectTypeName string,
) {
	// Generate the select struct type
	f.Line()
	f.Commentf("%s provides field selection for %s queries.", selectTypeName, node.NameGo())
	f.Type().Id(selectTypeName).Struct(
		jen.Id("SelectContext"),
	)

	b.generateSelectMethods(f, node, pkgLib, modelType, selectTypeName, false)

	// Generate array variant for edge traversal results.
	// Edge traversal is many-to-many, so every field is wrapped in a slice.
	arraySelectTypeName := selectTypeName + "Array"
	f.Line()
	f.Commentf("%s is the array variant of %s for edge traversal results.", arraySelectTypeName, selectTypeName)
	f.Type().Id(arraySelectTypeName).Struct(
		jen.Id("SelectContext"),
	)

	b.generateSelectMethods(f, node, pkgLib, modelType, arraySelectTypeName, true)
}

func (b *queryBuilder) generateSelectMethods(
	f *jen.File, node *field.NodeTable,
	pkgLib string, modelType jen.Code, selectTypeName string, arrayWrap bool,
) {
	for _, fld := range node.Fields {
		if fld.NameDatabase() == "id" {
			continue
		}

		switch typedFld := fld.(type) {
		case *field.Node:
			b.generateSelectNodeMethod(f, typedFld, selectTypeName, arrayWrap)
			continue
		case *field.Edge:
			b.generateSelectEdgeMethod(f, typedFld, fld.NameGo(), selectTypeName)
			continue
		}

		if s, ok := fld.(*field.Slice); ok {
			switch elem := s.Element().(type) {
			case *field.Node:
				_ = elem
				continue
			case *field.Edge:
				b.generateSelectEdgeMethod(f, elem, fld.NameGo(), selectTypeName)
				continue
			}
		}

		b.generateSelectFieldMethod(f, fld, pkgLib, modelType, selectTypeName, arrayWrap)
	}
}

func (b *queryBuilder) generateSelectFieldMethod(
	f *jen.File, fld field.Field,
	pkgLib string, _ jen.Code, selectTypeName string, arrayWrap bool,
) {
	ctx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
	}

	fieldGoType := fld.TypeGo()
	if arrayWrap {
		fieldGoType = jen.Index().Add(fieldGoType)
	}

	fieldDBName := fld.NameDatabase()
	fieldGoName := fld.NameGo()

	dict := jen.Dict{
		jen.Id("db"): jen.Id("s").Dot("DB"),
		jen.Id("buildFn"): jen.Func().Params().Op("*").Qual(pkgLib, "Result").Block(
			jen.Return(jen.Id("s").Dot("BuildFn").Call(jen.Lit(fieldDBName))),
		),
		jen.Id("distFn"): jen.Func().Params().Op("*").Qual(pkgLib, "Result").Block(
			jen.Return(jen.Id("s").Dot("DistFn").Call(jen.Lit(fieldDBName))),
		),
		jen.Id("firstFn"): jen.Func().Params().Op("*").Qual(pkgLib, "Result").Block(
			jen.Return(jen.Id("s").Dot("FirstFn").Call(jen.Lit(fieldDBName))),
		),
	}

	if !arrayWrap {
		if code := fld.CodeGen().SelectDecode(ctx); code != nil {
			dict[jen.Id("decodeFn")] = code
		}

		if code := fld.CodeGen().SelectDistDecode(ctx); code != nil {
			dict[jen.Id("distDecodeFn")] = code
		}
	}

	f.Line()
	f.Commentf("%s returns a SelectField for the %s field.", fieldGoName, fieldDBName)
	f.Func().
		Params(jen.Id("s").Id(selectTypeName)).
		Id(fieldGoName).
		Params().
		Id("SelectField").Types(fieldGoType).
		Block(
			jen.Return(jen.Id("SelectField").Types(fieldGoType).Values(dict)),
		)
}

func (b *queryBuilder) generateSelectNodeMethod(
	f *jen.File, fld *field.Node, selectTypeName string, arrayWrap bool,
) {
	fieldDBName := fld.NameDatabase()
	fieldGoName := fld.NameGo()
	linkedSelectType := fld.Table().NameGoLower() + "Select"
	if arrayWrap {
		linkedSelectType += "Array"
	}

	f.Line()
	f.Commentf("%s returns a select builder for traversing the %s record link.", fieldGoName, fieldDBName)
	f.Func().
		Params(jen.Id("s").Id(selectTypeName)).
		Id(fieldGoName).
		Params().
		Id(linkedSelectType).
		Block(
			jen.Return(jen.Id(linkedSelectType).Values(jen.Dict{
				jen.Id("SelectContext"): jen.Id("s").Dot("Prefixed").Call(jen.Lit(fieldDBName + ".")),
			})),
		)
}

func (b *queryBuilder) generateSelectEdgeMethod(
	f *jen.File, fld *field.Edge, fieldGoName, selectTypeName string,
) {
	edge := fld.Table()
	outTable := edge.Out.Table()
	// Edge traversal always produces array results, so use the array variant.
	linkedSelectType := outTable.NameGoLower() + "SelectArray"
	prefix := "->" + edge.NameDatabase() + "->" + outTable.NameDatabase() + "."

	f.Line()
	f.Commentf("%s returns a select builder for traversing the %s edge to %s.", fieldGoName, edge.NameDatabase(), outTable.NameGo())
	f.Func().
		Params(jen.Id("s").Id(selectTypeName)).
		Id(fieldGoName).
		Params().
		Id(linkedSelectType).
		Block(
			jen.Return(jen.Id(linkedSelectType).Values(jen.Dict{
				jen.Id("SelectContext"): jen.Id("s").Dot("Prefixed").Call(jen.Lit(prefix)),
			})),
		)
}

func (b *queryBuilder) generateRangeFn(
	f *jen.File, node *field.NodeTable,
	pkgLib, somPkg string, modelType jen.Code, varName string,
) {
	cid := node.Source.ComplexID
	keyType := b.SourceQual(cid.StructName)

	f.Line()
	f.Var().Id(varName).Op("=").Id("rangeFn").Types(modelType).Call(
		jen.Func().Params(
			jen.Id("q").Op("*").Qual(pkgLib, "Query").Types(modelType),
			jen.Id("from").Qual(somPkg, "RangeFrom"),
			jen.Id("to").Qual(somPkg, "RangeTo"),
		).String().BlockFunc(func(g *jen.Group) {
			g.Id("expr").Op(":=").Lit(":")

			// From bound
			g.If(jen.Op("!").Id("from").Dot("IsOpen").Call()).BlockFunc(func(inner *jen.Group) {
				inner.Id("key").Op(":=").Id("from").Dot("Value").Call().Assert(keyType)
				inner.Id("expr").Op("+=").Add(b.rangeBoundExpr(node, cid, "key"))
			})

			// Operator between bounds
			g.If(jen.Op("!").Id("from").Dot("IsOpen").Call().Op("&&").Op("!").Id("from").Dot("IsInclusive").Call()).Block(
				jen.Id("expr").Op("+=").Lit(">"),
			)
			g.Id("expr").Op("+=").Lit("..")
			g.If(jen.Op("!").Id("to").Dot("IsOpen").Call().Op("&&").Id("to").Dot("IsInclusive").Call()).Block(
				jen.Id("expr").Op("+=").Lit("="),
			)

			// To bound
			g.If(jen.Op("!").Id("to").Dot("IsOpen").Call()).BlockFunc(func(inner *jen.Group) {
				inner.Id("key").Op(":=").Id("to").Dot("Value").Call().Assert(keyType)
				inner.Id("expr").Op("+=").Add(b.rangeBoundExpr(node, cid, "key"))
			})

			g.Return(jen.Id("expr"))
		}),
	)
}

func (b *queryBuilder) rangeBoundExpr(node *field.NodeTable, cid *parser.FieldComplexID, keyVar string) jen.Code {
	var parts []jen.Code

	if cid.Kind == parser.IDTypeArray {
		parts = append(parts, jen.Lit("["))
		for i, sf := range cid.Fields {
			if i > 0 {
				parts = append(parts, jen.Lit(", "))
			}
			parts = append(parts, b.rangeFieldAsVar(node, sf, keyVar))
		}
		parts = append(parts, jen.Lit("]"))
	} else {
		parts = append(parts, jen.Lit("{"))
		for i, sf := range cid.Fields {
			if i > 0 {
				parts = append(parts, jen.Lit(", "))
			}
			parts = append(parts, jen.Lit(sf.DBName+": "))
			parts = append(parts, b.rangeFieldAsVar(node, sf, keyVar))
		}
		parts = append(parts, jen.Lit("}"))
	}

	result := parts[0]
	for _, p := range parts[1:] {
		result = jen.Add(result).Op("+").Add(p)
	}
	return result
}

func (b *queryBuilder) rangeFieldAsVar(node *field.NodeTable, sf parser.ComplexIDField, keyVar string) jen.Code {
	accessor := jen.Id(keyVar).Dot(sf.Name)
	wrappedValue := fieldValueFrom(b.input, b.basePkg, sf, accessor)
	return jen.Id("q").Dot("AsVar").Call(wrappedValue)
}

func (b *queryBuilder) generateStringIDRangeFn(
	f *jen.File, node *field.NodeTable,
	pkgLib, somPkg string, modelType jen.Code, varName string,
) {
	idTypeName := string(node.Source.IDType)

	f.Line()
	f.Var().Id(varName).Op("=").Id("rangeFn").Types(modelType).Call(
		jen.Func().Params(
			jen.Id("q").Op("*").Qual(pkgLib, "Query").Types(modelType),
			jen.Id("from").Qual(somPkg, "RangeFrom"),
			jen.Id("to").Qual(somPkg, "RangeTo"),
		).String().BlockFunc(func(g *jen.Group) {
			g.Id("expr").Op(":=").Lit(":")

			g.If(jen.Op("!").Id("from").Dot("IsOpen").Call()).Block(
				jen.Id("expr").Op("+=").Id("q").Dot("AsVar").Call(
					jen.Id("from").Dot("Value").Call().Assert(jen.Qual(somPkg, idTypeName)),
				),
			)

			g.If(jen.Op("!").Id("from").Dot("IsOpen").Call().Op("&&").Op("!").Id("from").Dot("IsInclusive").Call()).Block(
				jen.Id("expr").Op("+=").Lit(">"),
			)
			g.Id("expr").Op("+=").Lit("..")
			g.If(jen.Op("!").Id("to").Dot("IsOpen").Call().Op("&&").Id("to").Dot("IsInclusive").Call()).Block(
				jen.Id("expr").Op("+=").Lit("="),
			)

			g.If(jen.Op("!").Id("to").Dot("IsOpen").Call()).Block(
				jen.Id("expr").Op("+=").Id("q").Dot("AsVar").Call(
					jen.Id("to").Dot("Value").Call().Assert(jen.Qual(somPkg, idTypeName)),
				),
			)

			g.Return(jen.Id("expr"))
		}),
	)
}
