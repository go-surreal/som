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

	modelInfoVarName := node.NameGoLower() + "ModelInfo"

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

			g.Return(
				jen.Id("Builder").Types(modelType).
					Values(
						jen.Id("builder").Types(modelType).
							Values(jen.Dict{
								jen.Id("db"):    jen.Id("db"),
								jen.Id("query"): jen.Id("q"),
								jen.Id("info"):  jen.Id(modelInfoVarName),
							}),
					),
			)
		})

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
