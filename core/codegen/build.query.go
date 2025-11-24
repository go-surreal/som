package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
	"path"
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

	// Collect password field paths for OMIT clause
	passwordPaths := field.CollectPasswordPaths(node.Fields, "")

	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	// Build the omit slice literal
	var omitArg jen.Code
	if len(passwordPaths) > 0 {
		var omitItems []jen.Code
		for _, p := range passwordPaths {
			omitItems = append(omitItems, jen.Lit(p))
		}
		omitArg = jen.Index().String().Values(omitItems...)
	} else {
		omitArg = jen.Nil()
	}

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
								jen.Id("db"):       jen.Id("db"),
								jen.Id("query"):    jen.Qual(pkgLib, "NewQuery").Types(b.SourceQual(node.Name)).Call(jen.Lit(node.NameDatabase()), omitArg),
								jen.Id("convFrom"): jen.Qual(b.subPkg(def.PkgConv), "From"+node.NameGo()+"Ptr"),
								jen.Id("convTo"):   jen.Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()+"Ptr"),
							}),
					),
			),
		)

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
