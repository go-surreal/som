package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
	"path"
	"path/filepath"
	"strings"
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
	if err := b.embedStaticFiles(); err != nil {
		return err
	}

	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *queryBuilder) embedStaticFiles() error {
	tmpl := &embed.Template{
		GenerateOutPath: b.subPkg(""),
	}

	files, err := embed.Query(tmpl)
	if err != nil {
		return err
	}

	for _, file := range files {
		content := string(file.Content)
		content = strings.Replace(content, embedComment, codegenComment, 1)

		b.fs.Write(filepath.Join(b.path(), file.Path), []byte(content))
	}

	return nil
}

func (b *queryBuilder) buildFile(node *field.NodeTable) error {
	pkgLib := b.subPkg(def.PkgLib)

	f := jen.NewFile(b.pkgName)

	f.PackageComment(codegenComment)

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
								jen.Id("query"):    jen.Qual(pkgLib, "NewQuery").Types(b.SourceQual(node.Name)).Call(jen.Lit(node.NameDatabase())),
								jen.Id("convFrom"): jen.Qual(b.subPkg(def.PkgConv), "From"+node.NameGo()),
								jen.Id("convTo"):   jen.Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()),
							}),
					),
			),
		)

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
