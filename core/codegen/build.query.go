package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/core/codegen/dbtype"
	"github.com/marcbinz/sdb/core/codegen/def"
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

	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	return nil
}

func (b *queryBuilder) buildFile(node *dbtype.Node) error {
	f := jen.NewFile(b.pkgName)

	f.Type().Id(node.Name).Struct(
		jen.Id("build").Op("*").Qual(def.PkgLibBuilder, "Query"),
	)

	f.Func().Id("New" + node.Name).Params().
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Op("&").Id(node.Name).Values(jen.Dict{
				jen.Id("build"): jen.Qual(def.PkgLibBuilder, "NewQuery").Call(),
			})),
		)

	functions := []jen.Code{
		b.buildQueryFuncFilter(node),
		b.buildQueryFuncSort(node),
		b.buildQueryFuncOffset(node),
		b.buildQueryFuncLimit(node),
		b.buildQueryFuncUnique(node),
		b.buildQueryFuncTimeout(node),
		b.buildQueryFuncParallel(node),
		b.buildQueryFuncCount(node),
		b.buildQueryFuncExist(node),
		b.buildQueryFuncAll(node),
		b.buildQueryFuncAllIDs(node),
		b.buildQueryFuncFirst(node),
		b.buildQueryFuncFirstID(node),
		b.buildQueryFuncOnly(node),
		b.buildQueryFuncOnlyID(node),
	}

	for _, fn := range functions {
		f.Add(fn)
	}

	if err := f.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *queryBuilder) buildQueryFuncFilter(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Filter").Params(jen.Id("filters").Op("...").Qual(def.PkgLibFilter, "Of").Types(b.SourceQual(node.Name))).
		Op("*").Id(node.Name).
		Block(
			jen.For(jen.Id("_").Op(",").Id("f").Op(":=").Range().Id("filters")).
				Block(
					jen.Id("q").Dot("build").Dot("Where").Op("=").
						Append(jen.Id("q").Dot("build").Dot("Where"), jen.Qual(def.PkgLibBuilder, "Where").Call(jen.Id("f"))),
				),
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncSort(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Sort").Params(jen.Id("by").Op("...").Op("*").Qual(def.PkgLibSort, "Of").Types(b.SourceQual(node.Name))).
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncOffset(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Offset").Params(jen.Id("offset").Int()).
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncLimit(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Limit").Params(jen.Id("limit").Int()).
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncUnique(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Unique").Params().
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncTimeout(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Timeout").Params(jen.Id("timeout").Qual("time", "Duration")).
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncParallel(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Parallel").Params(jen.Id("parallel").Bool()).
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncCount(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Count").Params().
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncExist(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Exist").Params().
		Op("*").Id(node.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func (b *queryBuilder) buildQueryFuncAll(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("All").Params().
		Parens(jen.List(jen.Index().Op("*").Add(b.SourceQual(node.Name)), jen.Error())).
		Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncAllIDs(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("AllIDs").Params().
		Parens(jen.List(jen.Index().String(), jen.Error())).
		Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncFirst(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("First").Params().
		Parens(jen.List(jen.Op("*").Add(b.SourceQual(node.Name)), jen.Error())).
		Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncFirstID(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("FirstID").Params().
		Parens(jen.List(jen.String(), jen.Error())).
		Block(
			jen.Return(jen.Lit(""), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncOnly(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("Only").Params().
		Parens(jen.List(jen.Op("*").Add(b.SourceQual(node.Name)), jen.Error())).
		Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)
}

func (b *queryBuilder) buildQueryFuncOnlyID(node *dbtype.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(node.Name)).
		Id("OnlyID").Params().
		Parens(jen.List(jen.String(), jen.Error())).
		Block(
			jen.Return(jen.Lit(""), jen.Nil()),
		)
}
