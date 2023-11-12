package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"os"
	"path"
	"strings"
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

import "{{libPkg}}"

func All[T any](filters ...lib.Filter[T]) lib.Filter[T] {
	return lib.All[T](filters)
}

func Any[T any](filters ...lib.Filter[T]) lib.Filter[T] {
	return lib.Any[T](filters)
}
`

	content = strings.Replace(content, "{{libPkg}}", b.subPkg(def.PkgLib), 1)
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
	pkgLib := b.subPkg(def.PkgLib)

	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	if _, ok := elem.(*field.NodeTable); ok {
		file.Line()
		file.Var().Id(elem.NameGo()).Op("=").
			Id("new" + elem.NameGo()).Types(b.SourceQual(elem.NameGo())).
			Call(jen.Qual(pkgLib, "NewKey").Types(b.SourceQual(elem.NameGo())).Call())
	}

	file.Line()
	file.Add(b.whereNew(elem))

	file.Type().Id(elem.NameGoLower()).
		Types(jen.Id("T").Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Id("key").Qual(pkgLib, "Key").Types(jen.Id("T")))
			for _, f := range elem.GetFields() {
				if code := f.CodeGen().FilterDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	file.Line()
	file.Type().Id(elem.NameGoLower()+"Edges").
		Types(jen.Id("T").Any()).
		Struct(
			jen.Qual(pkgLib, "Filter").Types(jen.Id("T")),
			jen.Id("key").Qual(pkgLib, "Key").Types(jen.Id("T")),
		)

	fieldCtx.Receiver = jen.Id(elem.NameGoLower() + "Edges").Types(jen.Id("T"))

	for _, fld := range elem.GetFields() {
		isEdge := false

		if _, ok := fld.(*field.Edge); ok {
			isEdge = true
		}

		if slice, ok := fld.(*field.Slice); ok {
			if _, ok := slice.Element().(*field.Edge); ok {
				isEdge = true
			}
		}

		if !isEdge {
			continue
		}

		if code := fld.CodeGen().FilterFunc(fieldCtx); code != nil {
			file.Line()
			file.Add(code)
		}
	}

	file.Line()
	file.Type().Id(elem.NameGoLower()+"Slice").
		Types(jen.Id("T").Any()).
		Struct(
			jen.Qual(pkgLib, "Filter").Types(jen.Id("T")),
			jen.Op("*").Qual(pkgLib, "Slice").Types(jen.Id("T"), b.SourceQual(elem.NameGo())),
		)
}

func (b *filterBuilder) buildEdge(file *jen.File, edge *field.EdgeTable) {
	pkgLib := b.subPkg(def.PkgLib)

	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     edge,
	}

	file.Line()
	file.Var().Id(edge.NameGo()).Op("=").
		Id("new" + edge.NameGo()).Types(b.SourceQual(edge.NameGo())).
		Call(jen.Qual(pkgLib, "NewKey").Types(b.SourceQual(edge.NameGo())).Call())

	file.Line()
	file.Add(b.whereNew(edge))

	file.Line()
	file.Type().Id(edge.NameGoLower()).
		Types(jen.Id("T").Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Id("key").Qual(pkgLib, "Key").Types(jen.Id("T")))
			for _, f := range edge.GetFields() {
				if code := f.CodeGen().FilterDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	for _, fld := range edge.GetFields() {
		if code := fld.CodeGen().FilterFunc(fieldCtx); code != nil {
			file.Line()
			file.Add(code)
		}
	}

	file.Line()
	file.Type().Id(edge.NameGoLower()+"In").
		Types(jen.Id("T").Any()).
		Struct(
			jen.Qual(pkgLib, "Filter").Types(jen.Id("T")),
			jen.Id("key").Qual(pkgLib, "Key").Types(jen.Id("T")),
		)

	file.Line()
	file.Func().Id("new" + edge.NameGo() + "In").
		Types(jen.Id("T").Any()).
		Params(jen.Id("key").Qual(pkgLib, "Key").Types(jen.Id("T"))).
		Id(edge.NameGoLower() + "In").Types(jen.Id("T")).
		Block(
			jen.Return(
				jen.Id(edge.NameGoLower()+"In").Types(jen.Id("T")).Values(
					jen.Qual(pkgLib, "KeyFilter").Call(jen.Id("key")),
					jen.Id("key"),
				),
			),
		)

	file.Line()
	file.Func().
		Params(jen.Id("i").Id(edge.NameGoLower()+"In").Types(jen.Id("T"))).Id(edge.Out.NameGo()).
		Params(
			jen.Id("filters").Op("...").Qual(pkgLib, "Filter").Types(b.SourceQual(edge.Out.Table().NameGo())),
		).
		Id(edge.Out.Table().NameGoLower()+"Edges").Types(jen.Id("T")).
		Block(
			jen.Id("key").Op(":=").Qual(pkgLib, "EdgeIn").Call(
				jen.Id("i").Dot("key"),
				jen.Lit(edge.Out.NameDatabase()),
				jen.Id("filters"),
			),
			jen.Return(jen.Id(edge.Out.Table().NameGoLower()+"Edges").Types(jen.Id("T")).
				Values(
					jen.Qual(pkgLib, "KeyFilter").Call(jen.Id("key")),
					jen.Id("key"),
				),
			))

	file.Line()
	file.Type().Id(edge.NameGoLower()+"Out").
		Types(jen.Id("T").Any()).
		Struct(
			jen.Qual(pkgLib, "Filter").Types(jen.Id("T")),
			jen.Id("key").Qual(pkgLib, "Key").Types(jen.Id("T")),
		)

	file.Line()
	file.Func().Id("new" + edge.NameGo() + "Out").
		Types(jen.Id("T").Any()).
		Params(jen.Id("key").Qual(pkgLib, "Key").Types(jen.Id("T"))).
		Id(edge.NameGoLower() + "Out").Types(jen.Id("T")).
		Block(
			jen.Return(
				jen.Id(edge.NameGoLower()+"Out").Types(jen.Id("T")).Values(
					jen.Qual(pkgLib, "KeyFilter").Call(jen.Id("key")),
					jen.Id("key"),
				),
			),
		)

	file.Line()
	file.Func().
		Params(jen.Id("o").Id(edge.NameGoLower()+"Out").Types(jen.Id("T"))).Id(edge.In.NameGo()).
		Params(
			jen.Id("filters").Op("...").Qual(pkgLib, "Filter").Types(b.SourceQual(edge.In.Table().NameGo())),
		).
		Id(edge.In.Table().NameGoLower()+"Edges").Types(jen.Id("T")).
		Block(
			jen.Id("key").Op(":=").Qual(pkgLib, "EdgeOut").Call(
				jen.Id("o").Dot("key"),
				jen.Lit(edge.In.NameDatabase()),
				jen.Id("filters"),
			),
			jen.Return(jen.Id(edge.In.Table().NameGoLower()+"Edges").Types(jen.Id("T")).
				Values(
					jen.Qual(pkgLib, "KeyFilter").Call(jen.Id("key")),
					jen.Id("key"),
				),
			))
}

func (b *filterBuilder) whereNew(elem field.Element) jen.Code {
	pkgLib := b.subPkg(def.PkgLib)

	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		TargetPkg: b.basePkg,
		Table:     elem,
	}

	return jen.Func().Id("new" + elem.NameGo()).
		Types(jen.Id("T").Any()).
		Params(jen.Id("key").Qual(pkgLib, "Key").Types(jen.Id("T"))).
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
