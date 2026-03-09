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

	var fetchableFields []*field.Node
	for _, fld := range node.GetFields() {
		if nodeField, ok := fld.(*field.Node); ok {
			fetchableFields = append(fetchableFields, nodeField)
		}
	}

	if len(fetchableFields) > 0 {
		f.Line()
		f.Const().DefsFunc(func(g *jen.Group) {
			for i, nodeField := range fetchableFields {
				constName := typeName + "Fetched" + nodeField.NameGo()
				if i == 0 {
					g.Id(constName).Uint64().Op("=").Lit(1).Op("<<").Iota()
				} else {
					g.Id(constName)
				}
			}
		})
	}

	f.Line()
	f.Func().Id(node.NameGo() + "FetchBit").Params(jen.Id("field").String()).Uint64().Block(
		jen.Switch(jen.Id("field")).BlockFunc(func(g *jen.Group) {
			for _, nodeField := range fetchableFields {
				g.Case(jen.Lit(nodeField.NameDatabase())).Block(
					jen.Return(jen.Id(typeName + "Fetched" + nodeField.NameGo())),
				)
			}
			g.Default().Block(
				jen.Return(jen.Lit(0)),
			)
		}),
	)

	f.Line()
	f.Func().Id(node.NameGo() + "FetchFields").Params(jen.Id("bits").Uint64()).Index().String().BlockFunc(func(g *jen.Group) {
		g.Var().Id("fields").Index().String()
		for _, nodeField := range fetchableFields {
			g.If(jen.Id("bits").Op("&").Id(typeName+"Fetched"+nodeField.NameGo()).Op("!=").Lit(0)).Block(
				jen.Id("fields").Op("=").Append(jen.Id("fields"), jen.Lit(nodeField.NameDatabase())),
			)
		}
		g.Return(jen.Id("fields"))
	})

	f.Line()
	f.Func().Id(node.NameGo() + "SetFetched").Params(
		jen.Id("m").Op("*").Add(b.SourceQual(node.NameGo())),
		jen.Id("bits").Uint64(),
	).Block(
		jen.Id("m").Dot("Node").Dot("SetFetched").Call(jen.Id("bits")),
	)

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
