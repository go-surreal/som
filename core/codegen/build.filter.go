package codegen

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/core/codegen/dbtype"
	"github.com/marcbinz/sdb/core/codegen/def"
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

	for _, object := range b.objects {
		if err := b.buildFile(object); err != nil {
			return err
		}
	}

	return nil
}

func (b *filterBuilder) buildBaseFile() error {
	content := `package where

import filter "github.com/marcbinz/sdb/lib/filter"

func All[T any](filters ...filter.Of[T]) filter.Of[T] {
	return filter.All[T](filters)
}

func Any[T any](filters ...filter.Of[T]) filter.Of[T] {
	return filter.Any[T](filters)
}

func keyed(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}
`

	err := os.WriteFile(path.Join(b.path(), "where.go"), []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *filterBuilder) buildFile(elem dbtype.Element) error {
	file := jen.NewFile(b.pkgName)

	if _, ok := elem.(*dbtype.Node); ok {
		file.Var().Id(elem.NameGo()).Op("=").Id("new" + elem.NameGo()).Types(b.SourceQual(elem.NameGo())).Call(jen.Lit(""))
	}

	file.Add(b.whereNew(elem))

	file.Type().Id(strings.ToLower(elem.NameGo())).
		Types(jen.Id("T").Any()).
		StructFunc(func(g *jen.Group) {
			g.Add(jen.Id("key").String())
			for _, f := range elem.GetFields() {
				if code := f.FilterDefine(b.sourcePkgPath); code != nil {
					g.Add(code)
				}
			}
		})

	file.Type().Id(strings.ToLower(elem.NameGo())+"Slice").
		Types(jen.Id("T").Any()).
		Struct(
			jen.Id(strings.ToLower(elem.NameGo())).Types(jen.Id("T")),
			jen.Op("*").Qual(def.PkgLibFilter, "Slice").Types(b.SourceQual(elem.NameGo()), jen.Id("T")),
		)

	for _, fld := range elem.GetFields() {
		if code := fld.FilterFunc(b.sourcePkgPath, elem.NameGo()); code != nil {
			file.Add(code)
		}
	}

	if err := file.Save(path.Join(b.path(), elem.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *filterBuilder) whereNew(elem dbtype.Element) jen.Code {
	return jen.Func().Id("new" + elem.NameGo()).
		Types(jen.Id("T").Any()).
		Params(jen.Id("key").String()).
		Id(strings.ToLower(elem.NameGo())).Types(jen.Id("T")).
		Block(
			jen.Return(
				jen.Id(strings.ToLower(elem.NameGo())).Types(jen.Id("T")).
					Values(jen.DictFunc(func(d jen.Dict) {
						d[jen.Id("key")] = jen.Id("key")
						for _, f := range elem.GetFields() {
							if code := f.FilterInit(b.sourcePkgPath, elem.NameGo()); code != nil {
								d[jen.Id(f.NameGo())] = code
							}
						}
					})),
			),
		)
}
