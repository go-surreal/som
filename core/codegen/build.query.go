package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/codegen/field"
	"github.com/marcbinz/som/core/embed"
	"os"
	"path"
	"path/filepath"
	"strings"
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

func (b *queryBuilder) embedStaticFiles() error {
	files, err := embed.Query()
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

func (b *queryBuilder) buildFile(node *field.NodeTable) error {
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
		b.buildQueryFuncLive(node),
		// TODO: b.buildQueryFuncLiveDiff(node),
		b.buildQueryFuncDescribe(node), // TODO
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

func (b *queryBuilder) buildQueryFuncFilter(node *field.NodeTable) jen.Code {
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

func (b *queryBuilder) buildQueryFuncOrder(node *field.NodeTable) jen.Code {
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

func (b *queryBuilder) buildQueryFuncOrderRandom(node *field.NodeTable) jen.Code {
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

func (b *queryBuilder) buildQueryFuncOffset(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()
	nodeTypeNoLive := "Node" + node.NameGo() + "NoLive"

	return jen.
		Add(comment(`
Offset skips the first x records for the result set.
		`)).
		Func().Params(jen.Id("q").Id(nodeType)).
		Id("Offset").Params(jen.Id("offset").Int()).
		Id(nodeTypeNoLive).
		Block(
			jen.Id("q").Dot("query").Dot("Offset").Op("=").Id("offset"),
			jen.Return(jen.Id(nodeTypeNoLive).Values(jen.Id("q"))),
		)
}

func (b *queryBuilder) buildQueryFuncLimit(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()
	nodeTypeNoLive := "Node" + node.NameGo() + "NoLive"

	return jen.
		Add(comment(`
Limit restricts the query to return at most x records.
		`)).
		Func().Params(jen.Id("q").Id(nodeType)).
		Id("Limit").Params(jen.Id("limit").Int()).
		Id(nodeTypeNoLive).
		Block(
			jen.Id("q").Dot("query").Dot("Limit").Op("=").Id("limit"),
			jen.Return(jen.Id(nodeTypeNoLive).Values(jen.Id("q"))),
		)
}

func (b *queryBuilder) buildQueryFuncFetch(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()
	nodeTypeLive := "Node" + node.NameGo()

	return jen.
		Add(comment(`
Fetch can be used to return related records.
This works for both records links and edges.
		`)).
		Func().Params(jen.Id("q").Id(nodeType)).
		Id("Fetch").Params(jen.Id("fetch").Op("...").Qual(b.subPkg(def.PkgFetch), "Fetch_").Types(b.SourceQual(node.Name))).
		Id(nodeTypeLive).
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
			jen.Return(jen.Id(nodeTypeLive).Values(jen.Id("q"))),
		)
}

func (b *queryBuilder) buildQueryFuncTimeout(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()
	nodeTypeNoLive := "Node" + node.NameGo() + "NoLive"

	return jen.
		Add(comment(`
Timeout adds an execution time limit to the query.
When exceeded, the query call will return with an error.
		`)).
		Func().Params(jen.Id("q").Id(nodeType)).
		Id("Timeout").Params(jen.Id("timeout").Qual("time", "Duration")).
		Id(nodeTypeNoLive).
		Block(
			jen.Id("q").Dot("query").Dot("Timeout").Op("=").Id("timeout"),
			jen.Return(jen.Id(nodeTypeNoLive).Values(jen.Id("q"))),
		)
}

func (b *queryBuilder) buildQueryFuncParallel(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()
	nodeTypeNoLive := "Node" + node.NameGo() + "NoLive"

	return jen.
		Add(comment(`
Parallel tells SurrealDB that individual parts
of the query can be calculated in parallel.
This could lead to a faster execution.
		`)).
		Func().Params(jen.Id("q").Id(nodeType)).
		Id("Parallel").Params(jen.Id("parallel").Bool()).
		Id(nodeTypeNoLive).
		Block(
			jen.Id("q").Dot("query").Dot("Parallel").Op("=").Id("parallel"),
			jen.Return(jen.Id(nodeTypeNoLive).Values(jen.Id("q"))),
		)
}

func (b *queryBuilder) buildQueryFuncCount(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()

	return jen.Add(

		jen.
			Add(comment(`
Count returns the size of the result set, in other words the
number of records matching the conditions of the query.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("Count").Params(jen.Id("ctx").Qual("context", "Context")).
			Params(jen.Int(), jen.Error()).
			Block(
				jen.Id("req").Op(":=").Id("q").Dot("query").Dot("BuildAsCount").Call(),

				jen.Id("raw").Op(",").Err().Op(":=").Id("q").Dot("db").Dot("Query").Call(
					jen.Id("ctx"),
					jen.Id("req").Dot("Statement"),
					jen.Id("req").Dot("Variables"),
				),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Lit(0), jen.Err()),
				),

				jen.Var().Id("rawCount").Index().Id("queryResult").Types(jen.Id("countResult")),
				jen.Err().Op("=").Id("q").Dot("unmarshal").
					Call(jen.Id("raw"), jen.Op("&").Id("rawCount")),

				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Lit(0), jen.Qual("fmt", "Errorf").Call(jen.Lit("could not count records: %w"), jen.Err())),
				),

				jen.If(
					jen.Len(jen.Id("rawCount")).Op("<").Lit(1).Op("||").
						Len(jen.Id("rawCount").Index(jen.Lit(0)).Dot("Result")).Op("<").Lit(1),
				).Block(
					jen.Return(jen.Lit(0), jen.Nil()),
				),

				jen.Return(jen.Id("rawCount").Index(jen.Lit(0)).Dot("Result").Index(jen.Lit(0)).Dot("Count"), jen.Nil()),
			),

		jen.Line(),

		jen.
			Add(comment(`
CountAsync is the asynchronous version of Count.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("CountAsync").Params(jen.Id("ctx").Qual("context", "Context")).
			Op("*").Id("asyncResult").Types(jen.Int()).
			Block(
				jen.Return(jen.Id("async").Call(jen.Id("ctx"), jen.Id("q").Dot("Count"))),
			),
	)
}

func (b *queryBuilder) buildQueryFuncExists(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()

	return jen.Add(

		jen.
			Add(comment(`
Exists returns whether at least one record for the conditions
of the query exists or not. In other words it returns whether
the size of the result set is greater than 0.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("Exists").Params(jen.Id("ctx").Qual("context", "Context")).
			Params(jen.Bool(), jen.Error()).
			Block(
				jen.List(jen.Id("count"), jen.Err()).Op(":=").Id("q").Dot("Count").Call(jen.Id("ctx")),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.False(), jen.Err()),
				),
				jen.Return(jen.Id("count").Op(">").Lit(0), jen.Nil()),
			),

		jen.Line(),

		jen.
			Add(comment(`
ExistsAsync is the asynchronous version of Exists.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("ExistsAsync").Params(jen.Id("ctx").Qual("context", "Context")).
			Op("*").Id("asyncResult").Types(jen.Bool()).
			Block(
				jen.Return(jen.Id("async").Call(jen.Id("ctx"), jen.Id("q").Dot("Exists"))),
			),
	)
}

func (b *queryBuilder) buildQueryFuncAll(node *field.NodeTable) jen.Code {
	pkgConv := b.subPkg(def.PkgConv)

	nodeType := "node" + node.NameGo()

	resultType := jen.Index().Op("*").Add(b.SourceQual(node.Name))

	return jen.Add(

		jen.
			Add(comment(`
All returns all records matching the conditions of the query.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("All").Params(jen.Id("ctx").Qual("context", "Context")).
			Params(resultType, jen.Error()).
			Block(
				jen.Id("req").Op(":=").Id("q").Dot("query").Dot("BuildAsAll").Call(),

				jen.List(jen.Id("res"), jen.Err()).Op(":=").Id("q").Dot("db").Dot("Query").
					Call(
						jen.Id("ctx"),
						jen.Id("req").Dot("Statement"),
						jen.Id("req").Dot("Variables"),
					),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("could not query records: %w"), jen.Err())),
				),

				jen.Var().Id("rawNodes").Index().Id("queryResult").Types(jen.Op("*").Qual(pkgConv, node.NameGo())),

				jen.Err().Op("=").Id("q").Dot("unmarshal").Call(jen.Id("res"), jen.Op("&").Id("rawNodes")),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("could not unmarshal records: %w"), jen.Err())),
				),

				jen.If(jen.Len(jen.Id("rawNodes")).Op("<").Lit(1)).Block(
					jen.Return(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("empty result"))),
				),

				jen.Var().Id("nodes").Index().Op("*").Add(b.SourceQual(node.NameGo())),
				jen.For(jen.Id("_").Op(",").Id("rawNode").Op(":=").Range().Id("rawNodes").Index(jen.Lit(0)).Dot("Result")).
					Block(
						jen.Id("node").Op(":=").Qual(pkgConv, "To"+node.NameGo()).
							Call(jen.Id("rawNode")),
						jen.Id("nodes").Op("=").Append(jen.Id("nodes"), jen.Id("node")),
					),

				jen.Return(jen.Id("nodes"), jen.Nil()),
			),

		jen.Line(),

		jen.
			Add(comment(`
AllAsync is the asynchronous version of All.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("AllAsync").Params(jen.Id("ctx").Qual("context", "Context")).
			Op("*").Id("asyncResult").Types(resultType).
			Block(
				jen.Return(jen.Id("async").Call(jen.Id("ctx"), jen.Id("q").Dot("All"))),
			),
	)
}

func (b *queryBuilder) buildQueryFuncAllIDs(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()

	return jen.Add(

		jen.
			Add(comment(`
AllIDs returns the IDs of all records matching the conditions of the query.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("AllIDs").Params(jen.Id("ctx").Qual("context", "Context")).
			Parens(jen.List(jen.Index().String(), jen.Error())).
			Block(
				jen.Id("req").Op(":=").Id("q").Dot("query").Dot("BuildAsAllIDs").Call(),

				jen.List(jen.Id("res"), jen.Err()).Op(":=").Id("q").Dot("db").Dot("Query").Call(
					jen.Id("ctx"),
					jen.Id("req").Dot("Statement"),
					jen.Id("req").Dot("Variables"),
				),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("could not query records: %w"), jen.Err())),
				),

				jen.Var().Id("rawNodes").Index().Id("idNode"),
				jen.Err().Op("=").Id("q").Dot("unmarshal").Call(jen.Id("res"), jen.Op("&").Id("rawNodes")),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("could not unmarshal records: %w"), jen.Err())),
				),

				jen.Var().Id("ids").Index().String(),
				jen.For(jen.Id("_").Op(",").Id("rawNode").Op(":=").Range().Id("rawNodes")).
					Block(
						jen.Id("ids").Op("=").Append(jen.Id("ids"), jen.Id("rawNode").Dot("ID")),
					),

				jen.Return(jen.Id("ids"), jen.Nil()),
			),

		jen.Line(),

		jen.
			Add(comment(`
AllIDsAsync is the asynchronous version of AllIDs.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("AllIDsAsync").Params(jen.Id("ctx").Qual("context", "Context")).
			Op("*").Id("asyncResult").Types(jen.Index().String()).
			Block(
				jen.Return(jen.Id("async").Call(jen.Id("ctx"), jen.Id("q").Dot("AllIDs"))),
			),
	)
}

func (b *queryBuilder) buildQueryFuncFirst(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()

	resultType := jen.Op("*").Add(b.SourceQual(node.Name))

	return jen.Add(

		jen.
			Add(comment(`
First returns the first record matching the conditions of the query.
This comes in handy when using a filter for a field with unique values or when
sorting the result set in a specific order where only the first result is relevant.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("First").Params(jen.Id("ctx").Qual("context", "Context")).
			Params(resultType, jen.Error()).
			Block(
				jen.Id("q").Dot("query").Dot("Limit").Op("=").Lit(1),
				jen.List(jen.Id("res"), jen.Err()).Op(":=").Id("q").Dot("All").Call(jen.Id("ctx")),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Nil(), jen.Err()),
				),
				jen.If(jen.Len(jen.Id("res")).Op("<").Lit(1)).Block(
					jen.Return(jen.Nil(), jen.Qual("errors", "New").Call(jen.Lit("empty result"))),
				),
				jen.Return(jen.Id("res").Index(jen.Lit(0)), jen.Nil()),
			),

		jen.Line(),

		jen.
			Add(comment(`
FirstAsync is the asynchronous version of First.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("FirstAsync").Params(jen.Id("ctx").Qual("context", "Context")).
			Op("*").Id("asyncResult").Types(resultType).
			Block(
				jen.Return(jen.Id("async").Call(jen.Id("ctx"), jen.Id("q").Dot("First"))),
			),
	)
}

func (b *queryBuilder) buildQueryFuncFirstID(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()

	return jen.Add(

		jen.
			Add(comment(`
FirstID returns the ID of the first record matching the conditions of the query.
This comes in handy when using a filter for a field with unique values or when
sorting the result set in a specific order where only the first result is relevant.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("FirstID").Params(jen.Id("ctx").Qual("context", "Context")).
			Params(jen.String(), jen.Error()).
			Block(
				jen.Id("q").Dot("query").Dot("Limit").Op("=").Lit(1),
				jen.List(jen.Id("res"), jen.Err()).Op(":=").Id("q").Dot("AllIDs").Call(jen.Id("ctx")),
				jen.If(jen.Err().Op("!=").Nil()).Block(
					jen.Return(jen.Lit(""), jen.Err()),
				),
				jen.If(jen.Len(jen.Id("res")).Op("<").Lit(1)).Block(
					jen.Return(jen.Lit(""), jen.Qual("errors", "New").Call(jen.Lit("empty result"))),
				),
				jen.Return(jen.Id("res").Index(jen.Lit(0)), jen.Nil()),
			),

		jen.Line(),

		jen.
			Add(comment(`
FirstIDAsync is the asynchronous version of FirstID.
		`)).
			Func().Params(jen.Id("q").Id(nodeType)).
			Id("FirstIDAsync").Params(jen.Id("ctx").Qual("context", "Context")).
			Op("*").Id("asyncResult").Types(jen.String()).
			Block(
				jen.Return(jen.Id("async").Call(jen.Id("ctx"), jen.Id("q").Dot("FirstID"))),
			),
	)
}

func (b *queryBuilder) buildQueryFuncLive(node *field.NodeTable) jen.Code {
	nodeTypeLive := "Node" + node.NameGo()

	resultType := jen.Op("*").Add(b.SourceQual(node.Name))

	return jen.
		Add(comment(`
Live registers the constructed query as a live query.
Whenever something in the database changes that matches the 
query conditions, the result channel will receive an update.
If the context is canceled, the result channel will be closed.

Note: If you want both the current result set and live updates,
it is advised to execute the live query first. This is to ensure
data consistency. The other way around there could be missing
updates happening between the initial query and the live query.
		`)).
		Func().Params(jen.Id("q").Id(nodeTypeLive)).
		Id("Live").Params(jen.Id("ctx").Qual("context", "Context")).
		Params(
			jen.Op("<-").Chan().Id("LiveResult").Types(resultType),
			jen.Error(),
		).
		Block(
			jen.Id("req").Op(":=").Id("q").Dot("query").Dot("BuildAsLive").Call(),

			jen.List(jen.Id("resChan"), jen.Err()).Op(":=").Id("q").Dot("db").Dot("Live").Call(
				jen.Id("ctx"),
				jen.Id("req").Dot("Statement"),
				jen.Id("req").Dot("Variables"),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Qual("fmt", "Errorf").Call(jen.Lit("could not query live records: %w"), jen.Err())),
			),

			jen.Return(
				jen.Id("live").Call(
					jen.Id("ctx"),
					jen.Id("resChan"),
					jen.Id("q").Dot("unmarshal"),
					jen.Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()),
				),
				jen.Nil(),
			),
		)
}

// func (b *queryBuilder) buildQueryFuncLiveDiff(node *field.NodeTable) jen.Code {
// 	nodeTypeLive := "Node" + node.NameGo()
//
// 	resultType := jen.Op("*").Add(b.SourceQual(node.Name))
//
// 	return jen.
// 		Add(comment(`
// LiveDiff behaves like Live, but instead of receiving the full result
// set on every change, it only receives the actual changes.
// 		`)).
// 		Func().Params(jen.Id("q").Id(nodeTypeLive)).
// 		Id("LiveDiff").Params(jen.Id("ctx").Qual("context", "Context")).
// 		Op("*").Id("asyncResult").Types(resultType).
// 		Block(
// 			jen.Return(jen.Qual("errors", "New").Call(jen.Lit("not implemented"))),
// 		)
// }

func (b *queryBuilder) buildQueryFuncDescribe(node *field.NodeTable) jen.Code {
	nodeType := "node" + node.NameGo()

	return jen.
		Add(comment(`
Describe returns a string representation of the query.
While this might be a valid SurrealDB query, it
should only be used for debugging purposes.
		`)).
		Func().Params(jen.Id("q").Id(nodeType)).
		Id("Describe").Params().String().
		Block(
			jen.Id("req").Op(":=").Id("q").Dot("query").Dot("BuildAsAll").Call(),
			jen.Return(jen.Qual("strings", "TrimSpace").Call(
				jen.Id("req").Dot("Statement"),
			)),
		)
}

//
// -- HELPER
//

func comment(text string) jen.Code {
	var code jen.Statement

	text = strings.TrimSpace(text)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		code.Comment(line).Line()
	}

	return &code
}
