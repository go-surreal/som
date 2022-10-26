package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/parser"
	"path"
	"strings"
)

func buildQueryFile(input *parser.Result, queryPath string, model parser.Node) error {
	fileName := strings.ToLower(model.Name) + ".go"

	f := jen.NewFile("query")

	f.Type().Id(model.Name).Struct()

	functions := []jen.Code{
		buildQueryFuncFilter(input, model),
		buildQueryFuncSort(input, model),
		buildQueryFuncOffset(model),
		buildQueryFuncLimit(model),
		buildQueryFuncUnique(model),
		buildQueryFuncTimeout(model),
		buildQueryFuncParallel(model),
		buildQueryFuncCount(model),
		buildQueryFuncExist(model),
		buildQueryFuncAll(input, model),
		buildQueryFuncAllIDs(model),
		buildQueryFuncFirst(input, model),
		buildQueryFuncFirstID(model),
		buildQueryFuncOnly(input, model),
		buildQueryFuncOnlyID(model),
	}

	for _, fn := range functions {
		f.Add(fn)
	}

	f.Add(buildQueryToModel(input, model))
	f.Add(buildQueryFromModel(input, model))

	if err := f.Save(path.Join(queryPath, fileName)); err != nil {
		return err
	}

	return nil
}

func buildQueryFuncFilter(input *parser.Result, model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Filter").Params(jen.Id("filters").Op("...").Op("*").Qual(pkgLibFilter, "Of").Types(jen.Qual(input.PkgPath, model.Name))).
		Op("*").Id(model.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func buildQueryFuncSort(input *parser.Result, model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Sort").Params(jen.Id("by").Op("...").Op("*").Qual(pkgLibSort, "Of").Types(jen.Qual(input.PkgPath, model.Name))).
		Op("*").Id(model.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func buildQueryFuncOffset(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Offset").Params(jen.Id("offset").Int()).
		Op("*").Id(model.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func buildQueryFuncLimit(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Limit").Params(jen.Id("limit").Int()).
		Op("*").Id(model.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func buildQueryFuncUnique(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Unique").Params().
		Op("*").Id(model.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func buildQueryFuncTimeout(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Timeout").Params(jen.Id("timeout").Qual("time", "Duration")).
		Op("*").Id(model.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func buildQueryFuncParallel(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Parallel").Params(jen.Id("parallel").Bool()).
		Op("*").Id(model.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func buildQueryFuncCount(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Count").Params().
		Op("*").Id(model.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func buildQueryFuncExist(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Exist").Params().
		Op("*").Id(model.Name).
		Block(
			jen.Return(jen.Id("q")),
		)
}

func buildQueryFuncAll(input *parser.Result, model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("All").Params().
		Parens(jen.List(jen.Index().Op("*").Qual(input.PkgPath, model.Name), jen.Error())).
		Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)
}

func buildQueryFuncAllIDs(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("AllIDs").Params().
		Parens(jen.List(jen.Index().String(), jen.Error())).
		Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)
}

func buildQueryFuncFirst(input *parser.Result, model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("First").Params().
		Parens(jen.List(jen.Op("*").Qual(input.PkgPath, model.Name), jen.Error())).
		Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)
}

func buildQueryFuncFirstID(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("FirstID").Params().
		Parens(jen.List(jen.String(), jen.Error())).
		Block(
			jen.Return(jen.Lit(""), jen.Nil()),
		)
}

func buildQueryFuncOnly(input *parser.Result, model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("Only").Params().
		Parens(jen.List(jen.Op("*").Qual(input.PkgPath, model.Name), jen.Error())).
		Block(
			jen.Return(jen.Nil(), jen.Nil()),
		)
}

func buildQueryFuncOnlyID(model parser.Node) jen.Code {
	return jen.Func().
		Params(jen.Id("q").Op("*").Id(model.Name)).
		Id("OnlyID").Params().
		Parens(jen.List(jen.String(), jen.Error())).
		Block(
			jen.Return(jen.Lit(""), jen.Nil()),
		)
}

func buildQueryToModel(input *parser.Result, model parser.Node) jen.Code {
	return jen.Func().
		Id("to"+model.Name+"Model").
		Params(jen.Id("data").Map(jen.String()).Any()).
		Qual(input.PkgPath, model.Name).
		Block(
			jen.Return(jen.Qual(input.PkgPath, model.Name).Values()),
		)
}

func buildQueryFromModel(input *parser.Result, model parser.Node) jen.Code {
	return jen.Func().
		Id("from" + model.Name + "Model").
		Params(jen.Id("model").Qual(input.PkgPath, model.Name)).
		Map(jen.String()).Any().
		Block(
			jen.Return(jen.Map(jen.String()).Any().Values()),
		)
}
