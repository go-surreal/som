package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
	"path"
)

type fetchBuilder struct {
	*baseBuilder
}

func newFetchBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *fetchBuilder {
	return &fetchBuilder{
		baseBuilder: newBaseBuilder(input, fs, basePkg, pkgName),
	}
}

func (b *fetchBuilder) build() error {
	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *fetchBuilder) buildFile(node *field.NodeTable) error {
	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	f.Line()
	f.Var().Id(node.Name).Op("=").Id(node.NameGoLower()).Types(b.SourceQual(node.NameGo())).Call(jen.Lit(""))

	f.Line()
	f.Type().Id(node.NameGoLower()).
		Types(jen.Add(def.TypeModel).Any()).
		String()

	f.Line()
	f.Func().
		Params(jen.Id("n").Id(node.NameGoLower()).Types(def.TypeModel)).
		Id("fetch").Params(def.TypeModel).Block()

	for _, fld := range node.GetFields() {
		if nodeField, ok := fld.(*field.Node); ok {
			f.Line()
			f.Func().
				Params(jen.Id("n").Id(node.NameGoLower()).Types(def.TypeModel)).
				Id(nodeField.NameGo()).Params().
				Id(nodeField.Table().NameGoLower()).Types(def.TypeModel).
				Block(
					jen.Return(jen.Id(nodeField.Table().NameGoLower()).Types(def.TypeModel).
						Params(jen.Id("keyed").Call(jen.Id("n"), jen.Lit(nodeField.NameDatabase())))))
		}
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
