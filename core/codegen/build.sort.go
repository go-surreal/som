package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/core/codegen/dbtype"
	"github.com/marcbinz/sdb/core/codegen/def"
	"path"
	"strings"
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

	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *sortBuilder) buildFile(node *dbtype.Node) error {
	f := jen.NewFile(b.pkgName)

	f.Var().Id(node.Name).Op("=").Id("new" + node.Name).Call(jen.Lit(""))

	f.Add(b.byNew(node))

	f.Type().Id(strings.ToLower(node.Name)).StructFunc(func(g *jen.Group) {
		for _, f := range node.Fields {
			if code := f.SortDefine(b.SourceQual(node.Name)); code != nil {
				g.Add(code)
			}
		}
	})

	f.Func().Params(jen.Id(strings.ToLower(node.Name))).
		Id("Random").Params().
		Op("*").Qual(def.PkgLibSort, "Of").Types(b.SourceQual(node.Name)).
		Block(
			jen.Return(jen.Nil()),
		)

	if err := f.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *sortBuilder) byNew(node *dbtype.Node) jen.Code {
	return jen.Func().Id("new" + node.Name).
		Params(jen.Id("key").String()).
		Id(strings.ToLower(node.Name)).
		Block(
			jen.Return(
				jen.Id(strings.ToLower(node.Name)).Values(jen.DictFunc(func(d jen.Dict) {
					for _, f := range node.Fields {
						if code := f.SortInit(b.SourceQual(node.Name)); code != nil {
							d[jen.Id(f.NameGo())] = code
						}
					}
				})),
			),
		)
}
