package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type defineBuilder struct {
	*baseBuilder
}

func newDefineBuilder(input *input, basePath, basePkg, pkgName string) *defineBuilder {
	return &defineBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *defineBuilder) build() error {
	if err := b.createDir(); err != nil {
		return err
	}

	if err := b.embedStaticFiles(); err != nil {
		return err
	}

	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *defineBuilder) embedStaticFiles() error {
	files, err := embed.Define()
	if err != nil {
		return err
	}

	for _, file := range files {
		content := string(file.Content)
		content = strings.Replace(content, embedComment, codegenComment, 1)

		err := os.WriteFile(filepath.Join(b.path(), file.Path), []byte(content), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *defineBuilder) buildFile(node *field.NodeTable) error {
	pkgLib := b.subPkg(def.PkgLib)

	f := jen.NewFile(b.pkgName)

	f.PackageComment(codegenComment)

	nodeType := "node" + node.NameGo()
	nodeTypeLive := "Node" + node.NameGo()
	nodeTypeNoLive := "Node" + node.NameGo() + "NoLive"

	f.Line()
	f.Type().Id(nodeType).Struct(
		jen.Id("db").Id("Database"),
		jen.Id("query").Qual(pkgLib, "Query").Types(b.SourceQual(node.Name)),
		jen.Id("unmarshal").Func().Params(jen.Id("buf").Index().Byte(), jen.Id("val").Any()).Error(),
	)

	f.Line()
	f.Type().Id(nodeTypeLive).Struct(
		jen.Id(nodeType),
	)

	f.Line()
	f.Type().Id(nodeTypeNoLive).Struct(
		jen.Id(nodeType),
	)

	f.Line()
	f.Func().Id("New"+node.Name).
		Params(
			jen.Id("db").Id("Database"),
			jen.Id("unmarshal").Func().Params(jen.Id("buf").Index().Byte(), jen.Id("val").Any()).Error(),
		).
		Id(nodeTypeLive).
		Block(
			jen.Return(
				jen.Id(nodeTypeLive).Values(
					jen.Id(nodeType).Values(jen.Dict{
						jen.Id("db"):        jen.Id("db"),
						jen.Id("query"):     jen.Qual(pkgLib, "NewQuery").Types(b.SourceQual(node.Name)).Call(jen.Lit(node.NameDatabase())),
						jen.Id("unmarshal"): jen.Id("unmarshal"),
					}),
				),
			),
		)

	functions := []jen.Code{
		b.buildQueryFuncFilter(node),
		b.buildQueryFuncOrder(node),
		b.buildQueryFuncOrderRandom(node),
	}

	for _, fn := range functions {
		f.Line()
		f.Add(fn)
	}

	if err := f.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *defineBuilder) buildQueryFuncFilter(node *field.NodeTable) jen.Code {
	pkgLib := b.subPkg(def.PkgLib)

	nodeType := "node" + node.NameGo()
	nodeTypeLive := "Node" + node.NameGo()

	return jen.
		Add(comment(`
Filter adds a where statement to the query to
select records based on the given conditions.

Use where.All to chain multiple conditions
together that all need to match.
Use where.Any to chain multiple conditions
together where at least one needs to match.
		`)).
		Func().Params(jen.Id("q").Id(nodeType)).
		Id("Filter").Params(jen.Id("filters").Op("...").Qual(pkgLib, "Filter").Types(b.SourceQual(node.Name))).
		Id(nodeTypeLive).
		Block(
			jen.Id("q").Dot("query").Dot("Where").Op("=").
				Append(jen.Id("q").Dot("query").Dot("Where"), jen.Id("filters").Op("...")),
			jen.Return(jen.Id(nodeTypeLive).Values(jen.Id("q"))),
		)
}

func (b *defineBuilder) buildQueryFuncOrder(node *field.NodeTable) jen.Code {
	pkgLib := b.subPkg(def.PkgLib)

	nodeType := "node" + node.NameGo()
	nodeTypeNoLive := "Node" + node.NameGo() + "NoLive"

	return jen.
		Add(comment(`
Order sorts the returned records based on the given conditions.
If multiple conditions are given, they are applied one after the other.
Note: If OrderRandom is used within the same query,
it would override the sort conditions.
		`)).
		Func().Params(jen.Id("q").Id(nodeType)).
		Id("Order").Params(jen.Id("by").Op("...").Op("*").Qual(pkgLib, "Sort").Types(b.SourceQual(node.Name))).
		Id(nodeTypeNoLive).
		Block(
			jen.For(jen.Id("_").Op(",").Id("s").Op(":=").Range().Id("by")).
				Block(
					jen.Id("q").Dot("query").Dot("Sort").Op("=").
						Append(jen.Id("q").Dot("query").Dot("Sort"), jen.Parens(jen.Op("*").Qual(pkgLib, "SortBuilder")).Parens(jen.Id("s"))),
				),
			jen.Return(jen.Id(nodeTypeNoLive).Values(jen.Id("q"))),
		)
}

func (b *defineBuilder) buildQueryFuncOrderRandom(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()
	nodeTypeNoLive := "Node" + node.NameGo() + "NoLive"

	return jen.
		Add(comment(`
OrderRandom sorts the returned records in a random order.
Note: OrderRandom takes precedence over Order.
		`)).
		Func().Params(jen.Id("q").Id(nodeType)).
		Id("OrderRandom").Params().
		Id(nodeTypeNoLive).
		Block(
			jen.Id("q").Dot("query").Dot("SortRandom").Op("=").True(),
			jen.Return(jen.Id(nodeTypeNoLive).Values(jen.Id("q"))),
		)
}
