package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
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
	pkgLib := b.subPkg(def.PkgLib)
	pkgConv := b.subPkg(def.PkgConv)

	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	modelType := b.SourceQual(node.Name)
	convType := jen.Qual(pkgConv, node.Name)

	modelInfoVarName := node.NameGoLower() + "ModelInfo"

	f.Line()
	f.Commentf("%s holds the model-specific unmarshal functions for %s.", modelInfoVarName, node.NameGo())
	f.Var().Id(modelInfoVarName).Op("=").Id("ModelInfo").Types(modelType).Values(jen.Dict{
		jen.Id("UnmarshalAll"): jen.Func().Params(
			jen.Id("unmarshal").Func().Params(jen.Index().Byte(), jen.Any()).Error(),
			jen.Id("data").Index().Byte(),
		).Params(jen.Index().Op("*").Add(modelType), jen.Error()).Block(
			jen.Var().Id("rawNodes").Index().Id("queryResult").Types(jen.Op("*").Add(convType)),
			jen.If(jen.Err().Op(":=").Id("unmarshal").Call(jen.Id("data"), jen.Op("&").Id("rawNodes")), jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("could not unmarshal records: %w"), jen.Err())),
			),
			jen.If(jen.Len(jen.Id("rawNodes")).Op("<").Lit(1)).Block(
				jen.Return(jen.Nil(), jen.Nil()),
			),
			jen.Id("results").Op(":=").Make(jen.Index().Op("*").Add(modelType), jen.Len(jen.Id("rawNodes").Index(jen.Lit(0)).Dot("Result"))),
			jen.For(jen.List(jen.Id("i"), jen.Id("raw")).Op(":=").Range().Id("rawNodes").Index(jen.Lit(0)).Dot("Result")).Block(
				jen.Id("results").Index(jen.Id("i")).Op("=").Qual(pkgConv, "To"+node.NameGo()+"Ptr").Call(jen.Id("raw")),
			),
			jen.Return(jen.Id("results"), jen.Nil()),
		),
		jen.Id("UnmarshalOne"): jen.Func().Params(
			jen.Id("unmarshal").Func().Params(jen.Index().Byte(), jen.Any()).Error(),
			jen.Id("data").Index().Byte(),
		).Params(jen.Op("*").Add(modelType), jen.Error()).Block(
			jen.Var().Id("raw").Op("*").Add(convType),
			jen.If(jen.Err().Op(":=").Id("unmarshal").Call(jen.Id("data"), jen.Op("&").Id("raw")), jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(jen.Qual(pkgConv, "To"+node.NameGo()+"Ptr").Call(jen.Id("raw")), jen.Nil()),
		),
		jen.Id("UnmarshalLive"): jen.Func().Params(
			jen.Id("unmarshal").Func().Params(jen.Index().Byte(), jen.Any()).Error(),
			jen.Id("data").Index().Byte(),
		).Params(jen.Op("*").Add(modelType), jen.Error()).Block(
			jen.Var().Id("raw").Op("*").Add(convType),
			jen.If(jen.Err().Op(":=").Id("unmarshal").Call(jen.Id("data"), jen.Op("&").Id("raw")), jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Return(jen.Qual(pkgConv, "To"+node.NameGo()+"Ptr").Call(jen.Id("raw")), jen.Nil()),
		),
		jen.Id("UnmarshalSearchAll"): jen.Func().Params(
			jen.Id("unmarshal").Func().Params(jen.Index().Byte(), jen.Any()).Error(),
			jen.Id("data").Index().Byte(),
			jen.Id("clauses").Index().Qual(pkgLib, "SearchClause"),
		).Params(jen.Index().Qual(pkgLib, "SearchResult").Types(jen.Op("*").Add(modelType)), jen.Error()).Block(
			jen.Var().Id("rawNodes").Index().Id("queryResult").Types(jen.Id("searchRawResult").Types(jen.Op("*").Add(convType))),
			jen.If(jen.Err().Op(":=").Id("unmarshal").Call(jen.Id("data"), jen.Op("&").Id("rawNodes")), jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("could not unmarshal search records: %w"), jen.Err())),
			),
			jen.If(jen.Len(jen.Id("rawNodes")).Op("<").Lit(1)).Block(
				jen.Return(jen.Nil(), jen.Nil()),
			),
			jen.Var().Id("results").Index().Qual(pkgLib, "SearchResult").Types(jen.Op("*").Add(modelType)),
			jen.For(jen.List(jen.Id("_"), jen.Id("raw")).Op(":=").Range().Id("rawNodes").Index(jen.Lit(0)).Dot("Result")).Block(
				jen.Id("rec").Op(":=").Qual(pkgConv, "To"+node.NameGo()+"Ptr").Call(jen.Id("raw").Dot("Model")),
				jen.Id("result").Op(":=").Qual(pkgLib, "SearchResult").Types(jen.Op("*").Add(modelType)).Values(jen.Dict{
					jen.Id("Model"):      jen.Id("rec"),
					jen.Id("Scores"):     jen.Id("raw").Dot("Scores"),
					jen.Id("Highlights"): jen.Make(jen.Map(jen.Int()).String()),
					jen.Id("Offsets"):    jen.Make(jen.Map(jen.Int()).Index().Qual(pkgLib, "Offset")),
				}),
				jen.For(jen.List(jen.Id("_"), jen.Id("clause")).Op(":=").Range().Id("clauses")).Block(
					jen.If(jen.Id("clause").Dot("Highlights")).Block(
						jen.If(jen.List(jen.Id("hl"), jen.Id("ok")).Op(":=").Id("raw").Dot("Highlights").Index(jen.Id("clause").Dot("Ref")), jen.Id("ok")).Block(
							jen.Id("result").Dot("Highlights").Index(jen.Id("clause").Dot("Ref")).Op("=").Id("hl"),
						),
					),
					jen.If(jen.Id("clause").Dot("Offsets")).Block(
						jen.If(jen.List(jen.Id("offs"), jen.Id("ok")).Op(":=").Id("raw").Dot("Offsets").Index(jen.Id("clause").Dot("Ref")), jen.Id("ok")).Block(
							jen.Id("libOffsets").Op(":=").Make(jen.Index().Qual(pkgLib, "Offset"), jen.Len(jen.Id("offs"))),
							jen.For(jen.List(jen.Id("i"), jen.Id("off")).Op(":=").Range().Id("offs")).Block(
								jen.Id("libOffsets").Index(jen.Id("i")).Op("=").Qual(pkgLib, "Offset").Values(jen.Dict{
									jen.Id("Start"): jen.Id("off").Dot("Start"),
									jen.Id("End"):   jen.Id("off").Dot("End"),
								}),
							),
							jen.Id("result").Dot("Offsets").Index(jen.Id("clause").Dot("Ref")).Op("=").Id("libOffsets"),
						),
					),
				),
				jen.Id("results").Op("=").Append(jen.Id("results"), jen.Id("result")),
			),
			jen.Return(jen.Id("results"), jen.Nil()),
		),
	})

	f.Line()
	f.Commentf("New%s creates a new query builder for %s models.", node.NameGo(), node.NameGo())
	f.Func().Id("New"+node.Name).
		Params(
			jen.Id("db").Id("Database"),
		).
		Id("Builder").Types(modelType).
		Block(
			jen.Return(
				jen.Id("Builder").Types(modelType).
					Values(
						jen.Id("builder").Types(modelType).
							Values(jen.Dict{
								jen.Id("db"):    jen.Id("db"),
								jen.Id("query"): jen.Qual(pkgLib, "NewQuery").Types(modelType).Call(jen.Lit(node.NameDatabase())),
								jen.Id("info"):  jen.Id(modelInfoVarName),
							}),
					),
			),
		)

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
