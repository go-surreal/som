package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/core/codegen/dbtype"
	"github.com/marcbinz/sdb/core/codegen/field"
	"os"
	"path"
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
	content := `package with

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

	err := os.WriteFile(path.Join(b.path(), "fetch.go"), []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *fetchBuilder) buildFile(node *dbtype.Node) error {
	f := jen.NewFile(b.pkgName)

	f.Var().Id(node.Name).Op("=").Id(strings.ToLower(node.Name)).Types(b.SourceQual(node.NameGo())).Call(jen.Lit(""))

	f.Type().Id(strings.ToLower(node.Name)).
		Types(jen.Id("T").Any()).
		String()

	f.Func().
		Params(jen.Id("n").Id(strings.ToLower(node.NameGo())).Types(jen.Id("T"))).
		Id("fetch").Params(jen.Id("T")).Block()

	for _, fld := range node.GetFields() {
		if nodeField, ok := fld.(*field.Node); ok {
			f.Func().
				Params(jen.Id("n").Id(strings.ToLower(node.NameGo())).Types(jen.Id("T"))).
				Id(nodeField.NameGo()).Params().
				Id(strings.ToLower(nodeField.NodeName())).Types(jen.Id("T")).
				Block(
					jen.Return(jen.Id(strings.ToLower(nodeField.NodeName())).Types(jen.Id("T")).
						Params(jen.Id("keyed").Call(jen.Id("n"), jen.Lit(nodeField.NameDatabase())))))
		}
	}

	if err := f.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}
