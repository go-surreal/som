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

	pkgWith := b.subPkg(def.PkgFetch)

	f.Line()
	f.Func().Id("New"+node.Name).
		Params(
			jen.Id("db").Id("Database"),
		).
		Id("Builder").Types(b.SourceQual(node.Name), jen.Qual(b.subPkg(def.PkgConv), node.Name)).
		Block(
			jen.Return(
				jen.Id("Builder").Types(b.SourceQual(node.Name), jen.Qual(b.subPkg(def.PkgConv), node.Name)).
					Values(
						jen.Id("builder").Types(b.SourceQual(node.Name), jen.Qual(b.subPkg(def.PkgConv), node.Name)).
							Values(jen.Dict{
								jen.Id("db"):           jen.Id("db"),
								jen.Id("query"):        jen.Qual(pkgLib, "NewQuery").Types(b.SourceQual(node.Name)).Call(jen.Lit(node.NameDatabase())),
								jen.Id("convFrom"):     jen.Qual(b.subPkg(def.PkgConv), "From"+node.NameGo()+"Ptr"),
								jen.Id("convTo"):       jen.Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()+"Ptr"),
								jen.Id("fetchBitFn"):   jen.Qual(pkgWith, node.NameGo()+"FetchBit"),
								jen.Id("setFetchedFn"): jen.Qual(pkgWith, node.NameGo()+"SetFetched"),
							}),
					),
			),
		)

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
