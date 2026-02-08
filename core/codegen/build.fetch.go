package codegen

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/util/fs"
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
		if node.HasComplexID() {
			continue
		}
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *fetchBuilder) buildFile(node *field.NodeTable) error {
	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	typeName := node.NameGoLower()

	f.Line()
	f.Var().Id(node.Name).Op("=").Id(typeName).Types(b.SourceQual(node.NameGo())).Call(jen.Lit(""))

	f.Line()
	f.Type().Id(typeName).
		Types(jen.Add(def.TypeModel).Any()).
		String()

	f.Line()
	f.Func().
		Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
		Id("fetch").Params(def.TypeModel).Block()

	for _, fld := range node.GetFields() {
		if nodeField, ok := fld.(*field.Node); ok {
			relatedTable := nodeField.Table()
			f.Line()
			// Add comment for soft-delete relations
			if relatedTable.Source != nil && relatedTable.Source.SoftDelete {
				f.Comment(nodeField.NameGo() + " returns a fetch accessor for the " + nodeField.NameDatabase() + " relation.")
				f.Comment("Note: Soft-delete filtering does not apply to fetched relations.")
				f.Comment("All related records are returned regardless of their soft-delete status.")
			}
			f.Func().
				Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
				Id(nodeField.NameGo()).Params().
				Id(relatedTable.NameGoLower()).Types(def.TypeModel).
				Block(
					jen.Return(jen.Id(relatedTable.NameGoLower()).Types(def.TypeModel).
						Params(jen.Id("keyed").Call(jen.Id("n"), jen.Lit(nodeField.NameDatabase())))))
		}

		if sliceField, ok := fld.(*field.Slice); ok {
			if nodeElement, ok := sliceField.Element().(*field.Node); ok {
				relatedTable := nodeElement.Table()
				f.Line()
				// Add comment for soft-delete relations
				if relatedTable.Source != nil && relatedTable.Source.SoftDelete {
					f.Comment(sliceField.NameGo() + " returns a fetch accessor for the " + sliceField.NameDatabase() + " slice relation.")
					f.Comment("Note: Soft-delete filtering does not apply to fetched relations.")
					f.Comment("All related records are returned regardless of their soft-delete status.")
				}
				f.Func().
					Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
					Id(sliceField.NameGo()).Params().
					Id(relatedTable.NameGoLower()).Types(def.TypeModel).
					Block(
						jen.Return(jen.Id(relatedTable.NameGoLower()).Types(def.TypeModel).
							Params(jen.Id("keyed").Call(jen.Id("n"), jen.Lit(sliceField.NameDatabase())))))
			}
		}
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
