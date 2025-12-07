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

	typeName := node.NameGoLower()

	// Check if this node has soft-delete enabled
	if node.Source != nil && node.Source.SoftDelete {
		// Generate struct-based type with WithDeleted() method
		f.Line()
		f.Var().Id(node.Name).Op("=").Id(typeName).Types(b.SourceQual(node.NameGo())).Values(jen.Dict{
			jen.Id("field"): jen.Lit(""),
		})

		f.Line()
		f.Type().Id(typeName).Types(jen.Add(def.TypeModel).Any()).Struct(
			jen.Id("field").String(),
			jen.Id("withDeleted").Bool(),
		)

		// fetch() method for Fetch_ interface
		f.Line()
		f.Func().
			Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
			Id("fetch").Params(def.TypeModel).Block()

		// String() method for fmt.Sprintf
		f.Line()
		f.Func().
			Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
			Id("String").Params().String().
			Block(jen.Return(jen.Id("n").Dot("field")))

		// IncludesDeleted() method for FetchWithDeleted interface
		f.Line()
		f.Func().
			Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
			Id("IncludesDeleted").Params().Bool().
			Block(jen.Return(jen.Id("n").Dot("withDeleted")))

		// FetchField() method for FetchWithDeleted interface
		f.Line()
		f.Func().
			Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
			Id("FetchField").Params().String().
			Block(jen.Return(jen.Id("n").Dot("field")))

		// WithDeleted() method to get a copy with withDeleted=true
		f.Line()
		f.Func().
			Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
			Id("WithDeleted").Params().Id(typeName).Types(def.TypeModel).
			Block(jen.Return(jen.Id(typeName).Types(def.TypeModel).Values(jen.Dict{
				jen.Id("field"):       jen.Id("n").Dot("field"),
				jen.Id("withDeleted"): jen.True(),
			})))

		// Generate sub-node accessors
		for _, fld := range node.GetFields() {
			if nodeField, ok := fld.(*field.Node); ok {
				f.Line()
				f.Func().
					Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
					Id(nodeField.NameGo()).Params().
					Id(nodeField.Table().NameGoLower()).Types(def.TypeModel).
					Block(
						jen.Return(jen.Id(nodeField.Table().NameGoLower()).Types(def.TypeModel).
							Params(jen.Id("keyedStruct").Call(jen.Id("n").Dot("field"), jen.Lit(nodeField.NameDatabase())))))
			}
		}
	} else {
		// Generate string-based type (original behavior)
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
				f.Line()
				f.Func().
					Params(jen.Id("n").Id(typeName).Types(def.TypeModel)).
					Id(nodeField.NameGo()).Params().
					Id(nodeField.Table().NameGoLower()).Types(def.TypeModel).
					Block(
						jen.Return(jen.Id(nodeField.Table().NameGoLower()).Types(def.TypeModel).
							Params(jen.Id("keyed").Call(jen.Id("n"), jen.Lit(nodeField.NameDatabase())))))
			}
		}
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), node.FileName()))); err != nil {
		return err
	}

	return nil
}
