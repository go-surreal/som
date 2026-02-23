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

	convFn := jen.Qual(pkgConv, "To"+node.NameGo()+"Ptr")

	f.Line()
	f.Commentf("%s holds the model-specific unmarshal functions for %s.", modelInfoVarName, node.NameGo())
	f.Var().Id(modelInfoVarName).Op("=").Id("modelInfo").Types(modelType).Values(jen.Dict{
		jen.Id("UnmarshalAll"): jen.Func().Params(
			jen.Id("unmarshal").Func().Params(jen.Index().Byte(), jen.Any()).Error(),
			jen.Id("data").Index().Byte(),
		).Params(jen.Index().Op("*").Add(modelType), jen.Error()).Block(
			jen.Return(jen.Id("unmarshalAll").Call(jen.Id("unmarshal"), jen.Id("data"), convFn)),
		),
		jen.Id("UnmarshalOne"): jen.Func().Params(
			jen.Id("unmarshal").Func().Params(jen.Index().Byte(), jen.Any()).Error(),
			jen.Id("data").Index().Byte(),
		).Params(jen.Op("*").Add(modelType), jen.Error()).Block(
			jen.Return(jen.Id("unmarshalOne").Call(jen.Id("unmarshal"), jen.Id("data"), convFn)),
		),
		jen.Id("UnmarshalSearchAll"): jen.Func().Params(
			jen.Id("unmarshal").Func().Params(jen.Index().Byte(), jen.Any()).Error(),
			jen.Id("data").Index().Byte(),
			jen.Id("clauses").Index().Qual(pkgLib, "SearchClause"),
		).Params(jen.Index().Qual(pkgLib, "SearchResult").Types(jen.Op("*").Add(modelType)), jen.Error()).Block(
			jen.Return(jen.Id("unmarshalSearchAll").Call(jen.Id("unmarshal"), jen.Id("data"), jen.Id("clauses"), convFn)),
		),
	})

	if node.HasComplexID() {
		b.generateRangeFn(f, node, pkgLib, somPkg, modelType, rangeFnVarName)
	}

	f.Line()
	f.Commentf("New%s creates a new query builder for %s models.", node.NameGo(), node.NameGo())
	f.Func().Id("New"+node.Name).
		Params(
			jen.Id("db").Id("Database"),
		).
		Id("Builder").Types(modelType).
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
			}

			if node.HasComplexID() {
				builderDict[jen.Id("rangeFn")] = jen.Id(rangeFnVarName)
			}

			g.Return(
				jen.Id("Builder").Types(modelType).
					Values(
						jen.Id("builder").Types(modelType).
							Values(builderDict),
					),
			)
		})

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
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
			g.Var().Id("expr").String()

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
		parts = append(parts, jen.Lit(":["))
		for i, sf := range cid.Fields {
			if i > 0 {
				parts = append(parts, jen.Lit(", "))
			}
			parts = append(parts, b.rangeFieldAsVar(node, sf, keyVar))
		}
		parts = append(parts, jen.Lit("]"))
	} else {
		parts = append(parts, jen.Lit(":{"))
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
