package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/field"
	"os"
	"path"
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

	// Generate the base file.
	if err := b.buildBaseFile(); err != nil {
		return err
	}

	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *fetchBuilder) buildBaseFile() error {
	content := `

package with

// Fetch_ has a suffix of "_" to prevent clashes with node names.
type Fetch_[T any] interface {
	fetch(T)
}

func keyed[S ~string](base S, key string) string {
	if base == "" {
		return key
	}
	return string(base) + "." + key
}
`

	data := []byte(codegenComment + content)

	err := os.WriteFile(path.Join(b.path(), "fetch.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *fetchBuilder) buildFile(node *field.DatabaseNode) error {
	f := jen.NewFile(b.pkgName)

	f.PackageComment(codegenComment)

	f.Var().Id(node.Name).Op("=").Id(strcase.ToLowerCamel(node.Name)).Types(b.SourceQual(node.NameGo())).Call(jen.Lit(""))

	f.Type().Id(strcase.ToLowerCamel(node.Name)).
		Types(jen.Id("T").Any()).
		String()

	f.Func().
		Params(jen.Id("n").Id(strcase.ToLowerCamel(node.NameGo())).Types(jen.Id("T"))).
		Id("fetch").Params(jen.Id("T")).Block()

	for _, fld := range node.GetFields() {
		if nodeField, ok := fld.(*field.Node); ok {
			f.Func().
				Params(jen.Id("n").Id(strcase.ToLowerCamel(node.NameGo())).Types(jen.Id("T"))).
				Id(nodeField.NameGo()).Params().
				Id(strcase.ToLowerCamel(nodeField.NodeName())).Types(jen.Id("T")).
				Block(
					jen.Return(jen.Id(strcase.ToLowerCamel(nodeField.NodeName())).Types(jen.Id("T")).
						Params(jen.Id("keyed").Call(jen.Id("n"), jen.Lit(nodeField.NameDatabase())))))
		}
	}

	if err := f.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}
