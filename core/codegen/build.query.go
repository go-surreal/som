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

	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	f.Line()
	f.Func().Id("New"+node.Name).
		Params(
			jen.Id("db").Id("Database"),
		).
		Id("Builder").Types(b.SourceQual(node.Name), jen.Qual(b.subPkg(def.PkgConv), node.Name)).
		BlockFunc(func(g *jen.Group) {
			g.Id("q").Op(":=").Qual(pkgLib, "NewQuery").Types(b.SourceQual(node.Name)).Call(jen.Lit(node.NameDatabase()))

			// Inject soft delete filter if model has SoftDelete enabled
			if node.Source.SoftDelete {
				g.Comment("Automatically exclude soft-deleted records")
				pkgFilter := path.Join(b.basePkg, def.PkgFilter)
				g.Id("q").Dot("SoftDeleteFilter").Op("=").
					Qual(pkgFilter, node.Name).Dot("DeletedAt").Dot("Nil").Call(jen.Lit(true))
			}

			g.Return(
				jen.Id("Builder").Types(b.SourceQual(node.Name), jen.Qual(b.subPkg(def.PkgConv), node.Name)).
					Values(
						jen.Id("builder").Types(b.SourceQual(node.Name), jen.Qual(b.subPkg(def.PkgConv), node.Name)).
							Values(jen.Dict{
								jen.Id("db"):       jen.Id("db"),
								jen.Id("query"):    jen.Id("q"),
								jen.Id("convFrom"): jen.Qual(b.subPkg(def.PkgConv), "From"+node.NameGo()+"Ptr"),
								jen.Id("convTo"):   jen.Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()+"Ptr"),
							}),
					),
			)
		})

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
