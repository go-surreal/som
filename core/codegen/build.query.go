package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/codegen/field"
	"os"
	"path"
)

type queryBuilder struct {
	*baseBuilder
}

func newQueryBuilder(input *input, basePath, basePkg, pkgName string) *queryBuilder {
	return &queryBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *queryBuilder) build() error {
	if err := b.createDir(); err != nil {
		return err
	}

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

func (b *queryBuilder) buildBaseFile() error {
	content := `

package query

type Database interface {
	Query(statement string, vars any) (any, error)
}

type idNode struct {
	ID string
}

type countResult struct {
	Count int
}
`

	data := []byte(codegenComment + content)

	err := os.WriteFile(path.Join(b.path(), "query.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *queryBuilder) buildFile(node *field.NodeTable) error {
	f := jen.NewFile(b.pkgName)

	f.PackageComment(codegenComment)

	f.Type().Id(node.Name).Struct(
		jen.Id("db").Id("Database"),
		jen.Id("query").Qual(def.PkgLib, "Query").Types(b.SourceQual(node.Name)),
	)

	f.Func().Id("New" + node.Name).Params(jen.Id("db").Id("Database")).
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Op("&").Id(node.Name).Values(jen.Dict{
				jen.Id("db"):    jen.Id("db"),
				jen.Id("query"): jen.Qual(def.PkgLib, "NewQuery").Types(b.SourceQual(node.Name)).Call(jen.Lit(node.NameDatabase())),
			})),
		)

	functions := []jen.Code{
		b.buildQueryFuncFilter(node),
		b.buildQueryFuncOrder(node),
		b.buildQueryFuncOrderRandom(node),
		b.buildQueryFuncOffset(node),
		b.buildQueryFuncLimit(node),
		b.buildQueryFuncFetch(node), // TODO
		b.buildQueryFuncTimeout(node),
		b.buildQueryFuncParallel(node),
		b.buildQueryFuncCount(node),
		b.buildQueryFuncExists(node),
		b.buildQueryFuncAll(node),
		b.buildQueryFuncAllIDs(node),
		b.buildQueryFuncFirst(node),
		b.buildQueryFuncFirstID(node),
		b.buildQueryFuncDescribe(node), // TODO
	}

	for _, fn := range functions {
		f.Add(fn)
	}

	if err := f.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *queryBuilder) buildQueryFuncFilter(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Filter").Params(jen.Id("filters").Op("...").Qual(def.PkgLib, "Filter").Types(b.SourceQual(node.Name))).
		Id(node.Name).
		Block(
			jen.Id("q").Dot("query").Dot("Where").Op("=").
				Append(jen.Id("q").Dot("query").Dot("Where"), jen.Id("filters").Op("...")),
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncOrder(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Order").Params(jen.Id("by").Op("...").Op("*").Qual(def.PkgLib, "Sort").Types(b.SourceQual(node.Name))).
		Id(node.Name).
		Block(
			jen.For(jen.Id("_").Op(",").Id("s").Op(":=").Range().Id("by")).
				Block(
					jen.Id("q").Dot("query").Dot("Sort").Op("=").
						Append(jen.Id("q").Dot("query").Dot("Sort"), jen.Parens(jen.Op("*").Qual(def.PkgLib, "SortBuilder")).Parens(jen.Id("s"))),
				),
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncOrderRandom(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("OrderRandom").Params().
		Id(node.Name).
		Block(
			jen.Id("q").Dot("query").Dot("SortRandom").Op("=").True(),
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncOffset(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Offset").Params(jen.Id("offset").Int()).
		Id(node.Name).
		Block(
			jen.Id("q").Dot("query").Dot("Offset").Op("=").Id("offset"),
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncLimit(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Limit").Params(jen.Id("limit").Int()).
		Id(node.Name).
		Block(
			jen.Id("q").Dot("query").Dot("Limit").Op("=").Id("limit"),
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncFetch(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Fetch").Params(jen.Id("fetch").Op("...").Qual(b.subPkg(def.PkgFetch), "Fetch_").Types(b.SourceQual(node.Name))).
		Id(node.Name).
		Block(
			jen.For(jen.List(jen.Id("_"), jen.Id("f")).Op(":=").Range().Id("fetch")).
				Block(
					jen.If(
						jen.Id("field").Op(":=").Qual("fmt", "Sprintf").Call(jen.Lit("%v"), jen.Id("f")),
						jen.Id("field").Op("!=").Lit(""),
					).
						Block(
							jen.Id("q").Dot("query").Dot("Fetch").Op("=").
								Append(jen.Id("q").Dot("query").Dot("Fetch"), jen.Id("field")),
						),
				),
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncTimeout(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Timeout").Params(jen.Id("timeout").Qual("time", "Duration")).
		Id(node.Name).
		Block(
			jen.Id("q").Dot("query").Dot("Timeout").Op("=").Id("timeout"),
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncParallel(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Parallel").Params(jen.Id("parallel").Bool()).
		Id(node.Name).
		Block(
			jen.Id("q").Dot("query").Dot("Parallel").Op("=").Id("parallel"),
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncCount(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Count").Params().
		Params(jen.Int(), jen.Error()).
		Block(
			jen.Id("res").Op(":=").Id("q").Dot("query").Dot("BuildAsCount").Call(),
			jen.Id("raw").Op(",").Err().Op(":=").Id("q").Dot("db").Dot("Query").Call(
				jen.Id("res").Dot("Statement"),
				jen.Id("res").Dot("Variables"),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Lit(0), jen.Err()),
			),

			jen.Var().Id("rawCount").Id("countResult"),
			jen.List(jen.Id("ok"), jen.Err()).Op(":=").Qual(def.PkgSurrealDB, "UnmarshalRaw").
				Call(jen.Id("raw"), jen.Op("&").Id("rawCount")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Lit(0), jen.Err()),
			),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Lit(0), jen.Nil()),
			),

			jen.Return(jen.Id("rawCount").Dot("Count"), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncExists(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Exists").Params().
		Params(jen.Bool(), jen.Error()).
		Block(
			jen.List(jen.Id("count"), jen.Err()).Op(":=").Id("q").Dot("Count").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.False(), jen.Err()),
			),
			jen.Return(jen.Id("count").Op(">").Lit(0), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncAll(node *field.NodeTable) jen.Code {
	pkgConv := b.subPkg(def.PkgConv)

	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("All").Params().
		Params(jen.Index().Op("*").Add(b.SourceQual(node.Name)), jen.Error()).
		Block(
			jen.Id("res").Op(":=").Id("q").Dot("query").Dot("BuildAsAll").Call(),
			jen.Id("raw").Op(",").Err().Op(":=").
				Id("q").Dot("db").Dot("Query").
				Call(
					jen.Id("res").Dot("Statement"),
					jen.Id("res").Dot("Variables"),
				),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),

			jen.Var().Id("rawNodes").Index().Qual(b.subPkg(def.PkgConv), node.NameGo()),
			jen.List(jen.Id("ok"), jen.Err()).Op(":=").Qual(def.PkgSurrealDB, "UnmarshalRaw").
				Call(jen.Id("raw"), jen.Op("&").Id("rawNodes")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Nil(), jen.Nil()),
			),

			jen.Var().Id("nodes").Index().Op("*").Add(b.SourceQual(node.NameGo())),
			jen.For(jen.Id("_").Op(",").Id("rawNode").Op(":=").Range().Id("rawNodes")).
				Block(
					jen.Id("node").Op(":=").Qual(pkgConv, "To"+node.NameGo()).
						Call(jen.Id("rawNode")),
					jen.Id("nodes").Op("=").Append(jen.Id("nodes"), jen.Op("&").Id("node")),
				),

			jen.Return(jen.Id("nodes"), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncAllIDs(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("AllIDs").Params().
		Parens(jen.List(jen.Index().String(), jen.Error())).
		Block(
			jen.Id("res").Op(":=").Id("q").Dot("query").Dot("BuildAsAllIDs").Call(),
			jen.List(jen.Id("raw"), jen.Err()).Op(":=").Id("q").Dot("db").Dot("Query").Call(
				jen.Id("res").Dot("Statement"),
				jen.Id("res").Dot("Variables"),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),

			jen.Var().Id("rawNodes").Index().Op("*").Id("idNode"),
			jen.List(jen.Id("ok"), jen.Err()).Op(":=").Qual(def.PkgSurrealDB, "UnmarshalRaw").
				Call(jen.Id("raw"), jen.Op("&").Id("rawNodes")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Nil(), jen.Nil()),
			),

			jen.Var().Id("ids").Index().String(),
			jen.For(jen.Id("_").Op(",").Id("rawNode").Op(":=").Range().Id("rawNodes")).
				Block(
					jen.Id("ids").Op("=").Append(jen.Id("ids"), jen.Id("rawNode").Dot("ID")),
				),

			jen.Return(jen.Id("ids"), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncFirst(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("First").Params().
		Params(jen.Op("*").Add(b.SourceQual(node.Name)), jen.Error()).
		Block(
			jen.Id("q").Dot("query").Dot("Limit").Op("=").Lit(1),
			jen.List(jen.Id("res"), jen.Err()).Op(":=").Id("q").Dot("All").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.If(jen.Len(jen.Id("res")).Op("<").Lit(1)).Block(
				jen.Return(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("empty result"))),
			),
			jen.Return(jen.Id("res").Index(jen.Lit(0)), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncFirstID(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("FirstID").Params().
		Params(jen.String(), jen.Error()).
		Block(
			jen.Id("q").Dot("query").Dot("Limit").Op("=").Lit(1),
			jen.List(jen.Id("res"), jen.Err()).Op(":=").Id("q").Dot("AllIDs").Call(),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Lit(""), jen.Err()),
			),
			jen.If(jen.Len(jen.Id("res")).Op("<").Lit(1)).Block(
				jen.Return(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("empty result"))),
			),
			jen.Return(jen.Id("res").Index(jen.Lit(0)), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncDescribe(node *field.NodeTable) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Id(node.Name)).
		Id("Describe").Params().String().
		Block(
			jen.Id("res").Op(":=").Id("q").Dot("query").Dot("BuildAsAll").Call(),
			jen.Return(jen.Id("res").Dot("Statement")),
		)
}
