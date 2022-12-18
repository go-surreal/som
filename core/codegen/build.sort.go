package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/field"
	"os"
	"path"
)

type sortBuilder struct {
	*baseBuilder
}

func newSortBuilder(input *input, basePath, basePkg, pkgName string) *sortBuilder {
	return &sortBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *sortBuilder) build() error {
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

func (b *sortBuilder) buildBaseFile() error {
	content := `

package by

func keyed(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}
`

	data := []byte(codegenComment + content)

	err := os.WriteFile(path.Join(b.path(), "sort.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *sortBuilder) buildFile(node *field.NodeTable) error {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		Table:     node,
	}

	f := jen.NewFile(b.pkgName)

	f.PackageComment(codegenComment)

	f.Var().Id(node.Name).Op("=").Id("new" + node.Name).Types(b.SourceQual(node.NameGo())).Call(jen.Lit(""))

	f.Add(b.byNew(node))

	f.Type().Id(strcase.ToLowerCamel(node.Name)).
		Types(jen.Id("T").Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Id("key").String())
			for _, f := range node.Fields {
				if code := f.CodeGen().SortDefine(fieldCtx); code != nil {
					g.Add(code)
				}
			}
		})

	for _, fld := range node.GetFields() {
		if code := fld.CodeGen().SortFunc(fieldCtx); code != nil {
			f.Add(code)
		}
	}

	if err := f.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *sortBuilder) byNew(node *field.NodeTable) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		Table:     node,
	}

	return jen.Func().Id("new" + node.Name).
		Types(jen.Id("T").Any()).
		Params(jen.Id("key").String()).
		Id(strcase.ToLowerCamel(node.Name)).Types(jen.Id("T")).
		Block(
			jen.Return(
				jen.Id(strcase.ToLowerCamel(node.Name)).Types(jen.Id("T")).
					Values(jen.DictFunc(func(d jen.Dict) {
						d[jen.Id("key")] = jen.Id("key")
						for _, f := range node.Fields {
							if code := f.CodeGen().SortInit(fieldCtx); code != nil {
								d[jen.Id(f.NameGo())] = code
							}
						}
					})),
			),
		)
}
