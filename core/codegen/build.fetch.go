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

	// Collect fetchable fields (Node-type fields)
	var fetchableFields []*field.Node
	for _, fld := range node.GetFields() {
		if nodeField, ok := fld.(*field.Node); ok {
			fetchableFields = append(fetchableFields, nodeField)
		}
	}

	// Generate fetch bit constants using iota (only if there are fetchable fields)
	if len(fetchableFields) > 0 {
		f.Line()
		f.Const().DefsFunc(func(g *jen.Group) {
			for i, nodeField := range fetchableFields {
				constName := node.NameGoLower() + "Fetched" + nodeField.NameGo()
				if i == 0 {
					g.Id(constName).Uint64().Op("=").Lit(1).Op("<<").Iota()
				} else {
					g.Id(constName)
				}
			}
		})
	}

	// Always generate field-to-bit mapping function (exported for use by repo package)
	// Returns 0 for unknown fields or when there are no fetchable fields
	f.Line()
	f.Func().Id(node.NameGo() + "FetchBit").Params(jen.Id("field").String()).Uint64().Block(
		jen.Switch(jen.Id("field")).BlockFunc(func(g *jen.Group) {
			for _, nodeField := range fetchableFields {
				g.Case(jen.Lit(nodeField.NameDatabase())).Block(
					jen.Return(jen.Id(node.NameGoLower() + "Fetched" + nodeField.NameGo())),
				)
			}
			g.Default().Block(
				jen.Return(jen.Lit(0)),
			)
		}),
	)

	// Generate SetFetched function for use by query builder
	f.Line()
	f.Func().Id(node.NameGo() + "SetFetched").Params(
		jen.Id("m").Op("*").Add(b.SourceQual(node.NameGo())),
		jen.Id("bits").Uint64(),
	).Block(
		jen.Id("m").Dot("Node").Dot("SetFetched").Call(jen.Id("bits")),
	)

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

	for _, nodeField := range fetchableFields {
		f.Line()
		f.Func().
			Params(jen.Id("n").Id(node.NameGoLower()).Types(def.TypeModel)).
			Id(nodeField.NameGo()).Params().
			Id(nodeField.Table().NameGoLower()).Types(def.TypeModel).
			Block(
				jen.Return(jen.Id(nodeField.Table().NameGoLower()).Types(def.TypeModel).
					Params(jen.Id("keyed").Call(jen.Id("n"), jen.Lit(nodeField.NameDatabase())))))
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
