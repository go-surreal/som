package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/field"
	"github.com/marcbinz/som/core/embed"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type fetchBuilder struct {
	*baseBuilder
}

func newFetchBuilder(input *input, basePath, basePkg, pkgName string) *fetchBuilder {
	return &fetchBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *fetchBuilder) build() error {
	if err := b.createDir(); err != nil {
		return err
	}

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

func (b *fetchBuilder) embedStaticFiles() error {
	files, err := embed.Fetch()
	if err != nil {
		return err
	}

	for _, file := range files {
		content := string(file.Content)
		content = strings.Replace(content, embedComment, codegenComment, 1)

		err := os.WriteFile(filepath.Join(b.path(), file.Path), []byte(content), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *fetchBuilder) buildFile(node *field.NodeTable) error {
	f := jen.NewFile(b.pkgName)

	f.PackageComment(codegenComment)

	f.Line()
	f.Var().Id(node.Name).Op("=").Id(node.NameGoLower()).Types(b.SourceQual(node.NameGo())).Call(jen.Lit(""))

	f.Line()
	f.Type().Id(node.NameGoLower()).
		Types(jen.Id("T").Any()).
		String()

	f.Line()
	f.Func().
		Params(jen.Id("n").Id(node.NameGoLower()).Types(jen.Id("T"))).
		Id("fetch").Params(jen.Id("T")).Block()

	for _, fld := range node.GetFields() {
		if nodeField, ok := fld.(*field.Node); ok {
			f.Line()
			f.Func().
				Params(jen.Id("n").Id(node.NameGoLower()).Types(jen.Id("T"))).
				Id(nodeField.NameGo()).Params().
				Id(nodeField.Table().NameGoLower()).Types(jen.Id("T")).
				Block(
					jen.Return(jen.Id(nodeField.Table().NameGoLower()).Types(jen.Id("T")).
						Params(jen.Id("keyed").Call(jen.Id("n"), jen.Lit(nodeField.NameDatabase())))))
		}
	}

	if err := f.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}
