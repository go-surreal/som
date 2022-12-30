package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/codegen/field"
	"os"
	"path"
)

type filterBuilder struct {
	*baseBuilder
}

func newFilterBuilder(input *input, basePath, basePkg, pkgName string) *filterBuilder {
	return &filterBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *filterBuilder) build() error {
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

	for _, edge := range b.edges {
		if err := b.buildFile(edge); err != nil {
			return err
		}
	}

	for _, object := range b.objects {
		if err := b.buildFile(object); err != nil {
			return err
		}
	}

	return nil
}

func (b *filterBuilder) buildBaseFile() error {
	content := `

package where

import filter "github.com/marcbinz/som/lib/filter"

func All[T any](filters ...filter.Of[T]) filter.Of[T] {
	return filter.All[T](filters)
}

func Any[T any](filters ...filter.Of[T]) filter.Of[T] {
	return filter.Any[T](filters)
}
`

	data := []byte(codegenComment + content)

	err := os.WriteFile(path.Join(b.path(), "where.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *filterBuilder) buildFile(elem field.Element) error {
	file := jen.NewFile(b.pkgName)

	file.PackageComment(codegenComment)

	if edge, ok := elem.(*field.EdgeTable); ok {
		b.buildEdge(file, edge)
	} else {
		b.buildOther(file, elem)
	}

	if err := file.Save(path.Join(b.path(), elem.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *filterBuilder) buildOther(file *jen.File, elem field.Element) {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		Table:     elem,
	}

	if _, ok := elem.(*field.NodeTable); ok {
		file.Var().Id(elem.NameGo()).Op("=").
			Id("new" + elem.NameGo()).Types(b.SourceQual(elem.NameGo())).
			Call(jen.Qual(def.PkgLibFilter, "NewKey").Call())
	}

	file.Add(b.whereNew(elem))

	file.Type().Id(elem.NameGoLower()).
		Types(jen.Id("T").Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Id("key").Qual(def.PkgLibFilter, "Key"))
			for _, f := range elem.GetFields() {
				if code := f.CodeGen().FilterDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	file.Type().Id(elem.NameGoLower()+"Slice").
		Types(jen.Id("T").Any()).
		Struct(
			jen.Id(elem.NameGoLower()).Types(jen.Id("T")),
			jen.Op("*").Qual(def.PkgLibFilter, "Slice").Types(jen.Id("T"), b.SourceQual(elem.NameGo())),
		)

	for _, fld := range elem.GetFields() {
		if code := fld.CodeGen().FilterFunc(fieldCtx); code != nil {
			file.Add(code)
		}
	}
}

func (b *filterBuilder) whereNew(elem field.Element) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		Table:     elem,
	}

	return jen.Func().Id("new" + elem.NameGo()).
		Types(jen.Id("T").Any()).
		Params(jen.Id("key").Qual(def.PkgLibFilter, "Key")).
		Id(elem.NameGoLower()).Types(jen.Id("T")).
		Block(
			jen.Return(
				jen.Id(elem.NameGoLower()).Types(jen.Id("T")).
					Values(jen.DictFunc(func(d jen.Dict) {
						d[jen.Id("key")] = jen.Id("key")
						for _, f := range elem.GetFields() {
							if code := f.CodeGen().FilterInit(fieldCtx); code != nil {
								d[jen.Id(f.NameGo())] = code
							}
						}
					})),
			),
		)
}

func (b *filterBuilder) buildEdge(file *jen.File, edge *field.EdgeTable) {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		Table:     edge,
	}

	file.Func().Id("new" + edge.NameGo() + "In").
		Types(jen.Id("T").Any()).
		Params(jen.Id("key").Qual(def.PkgLibFilter, "Key")).
		Id(edge.NameGoLower() + "In").Types(jen.Id("T")).
		Block(
			jen.Return(
				jen.Id(edge.NameGoLower() + "In").Types(jen.Id("T")).Values(
					jen.Id(edge.NameGoLower()).Types(jen.Id("T")).
						Values(jen.DictFunc(func(d jen.Dict) {
							d[jen.Id("key")] = jen.Id("key")
							for _, f := range edge.GetFields() {
								if code := f.CodeGen().FilterInit(fieldCtx); code != nil {
									d[jen.Id(f.NameGo())] = code
								}
							}
						})),
				),
			),
		)

	file.Type().Id(edge.NameGoLower() + "In").
		Types(jen.Id("T").Any()).
		Struct(
			jen.Id(edge.NameGoLower()).Types(jen.Id("T")),
		)

	// TODO: below! -> node.FilterFunc?
	// if out, ok := edge.Out.(*field.Node); ok {
	file.Func().
		Params(jen.Id("i").Id(edge.NameGoLower() + "In").Types(jen.Id("T"))).
		Id(edge.Out.NameGo()).Params().
		Id(edge.Out.NameGoLower()).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + edge.Out.NameGo()).Types(jen.Id("T")).
				Params(jen.Id("i").Dot("key").Dot("In").Call(jen.Lit(edge.Out.NameDatabase())))))
	// }

	file.Func().Id("new" + edge.NameGo() + "Out").
		Types(jen.Id("T").Any()).
		Params(jen.Id("key").Qual(def.PkgLibFilter, "Key")).
		Id(edge.NameGoLower() + "Out").Types(jen.Id("T")).
		Block(
			jen.Return(
				jen.Id(edge.NameGoLower() + "Out").Types(jen.Id("T")).Values(
					jen.Id(edge.NameGoLower()).Types(jen.Id("T")).
						Values(jen.DictFunc(func(d jen.Dict) {
							d[jen.Id("key")] = jen.Id("key")
							for _, f := range edge.GetFields() {
								if code := f.CodeGen().FilterInit(fieldCtx); code != nil {
									d[jen.Id(f.NameGo())] = code
								}
							}
						})),
				),
			),
		)

	file.Type().Id(edge.NameGoLower() + "Out").
		Types(jen.Id("T").Any()).
		Struct(
			jen.Id(edge.NameGoLower()).Types(jen.Id("T")),
		)

	// TODO: below! -> node.FilterFunc?
	// if in, ok := edge.In.(*field.Node); ok {
	file.Func().
		Params(jen.Id("o").Id(edge.NameGoLower() + "Out").Types(jen.Id("T"))).
		Id(edge.In.NameGo()).Params().
		Id(edge.In.NameGoLower()).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + edge.In.NameGo()).Types(jen.Id("T")).
				Params(jen.Id("o").Dot("key").Dot("Out").Call(jen.Lit(edge.In.NameDatabase())))))
	// }

	file.Type().Id(edge.NameGoLower()).
		Types(jen.Id("T").Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Id("key").Qual(def.PkgLibFilter, "Key"))
			for _, f := range edge.GetFields() {
				if code := f.CodeGen().FilterDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	for _, fld := range edge.GetFields() {
		if code := fld.CodeGen().FilterFunc(fieldCtx); code != nil {
			file.Add(code)
		}
	}
}
