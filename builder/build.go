package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/parser"
	"os"
	"path"
	"strings"
)

const (
	dirBase      = "sdb"
	dirQuery     = "query"
	dirWhere     = "where"
	dirPredicate = "predicate"
	dirBy        = "by"
)

const (
	pkgLibFilter = "github.com/marcbinz/sdb/lib/filter"
	pkgLibSort   = "github.com/marcbinz/sdb/lib/sort"

	pkgQuery     = "github.com/marcbinz/sdb/example/gen/sdb/query"     // TODO
	pkgPredicate = "github.com/marcbinz/sdb/example/gen/sdb/predicate" // TODO
)

func Build(input *parser.Result, outDir string) error {
	basePath := path.Join(outDir, dirBase)
	queryPath := path.Join(basePath, dirQuery)
	wherePath := path.Join(basePath, dirWhere)
	predicatePath := path.Join(basePath, dirPredicate)
	byPath := path.Join(basePath, dirBy)

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(queryPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(wherePath, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(predicatePath, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(byPath, os.ModePerm)
	if err != nil {
		return err
	}

	if err := buildSdbFile(basePath); err != nil {
		return err
	}

	if err := buildBaseWhereFile(wherePath); err != nil {
		return err
	}

	for _, model := range input.Nodes {

		if err := baseFile(basePath, model, input.PkgPath); err != nil {
			return err
		}

		if err := buildQueryFile(input, queryPath, model); err != nil {
			return err
		}

		if err := buildWhereFile(input, wherePath, model, input.PkgPath); err != nil {
			return err
		}

		if err := predicateFile(predicatePath, model, input.PkgPath); err != nil {
			return err
		}

		if err := buildByFile(input, byPath, model); err != nil {
			return err
		}
	}

	return nil
}

//
// -- CLIENT
//

func buildSdbFile(basePath string) error {
	fileName := "0_sdb.go"

	f := jen.NewFile("sdb")

	f.Type().Id("Client").
		Struct(
			jen.Id("db").Op("*").Qual("github.com/surrealdb/surrealdb.go", "DB"),
		)

	if err := f.Save(path.Join(basePath, fileName)); err != nil {
		return err
	}

	return nil
}

//
// -- BASE
//

func baseFile(basePath string, model parser.Node, modelPkg string) error {
	fileName := strings.ToLower(model.Name) + ".go"

	f := jen.NewFile("sdb")

	f.Var().Id(model.Name).Id(strings.ToLower(model.Name))

	f.Type().Id(strings.ToLower(model.Name)).Struct()

	f.Func().
		Params(jen.Id(strings.ToLower(model.Name))).
		Id("Query").Params().
		Op("*").Qual(pkgQuery, model.Name).
		Block(
			jen.Return(jen.Op("&").Qual(pkgQuery, model.Name).Values()),
		)

	f.Func().
		Params(jen.Id(strings.ToLower(model.Name))).
		Id("Create").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower(model.Name)).Op("*").Qual(modelPkg, model.Name),
		).
		Error().
		Block(
			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id(strings.ToLower(model.Name))).
		Id("Update").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower(model.Name)).Op("*").Qual(modelPkg, model.Name),
		).
		Error().
		Block(
			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id(strings.ToLower(model.Name))).
		Id("Delete").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower(model.Name)).Op("*").Qual(modelPkg, model.Name),
		).
		Error().
		Block(
			jen.Return(jen.Nil()),
		)

	if err := f.Save(path.Join(basePath, fileName)); err != nil {
		return err
	}

	return nil
}

//
// -- PREDICATE
//

func predicateFile(predicatePath string, model parser.Node, modelPkg string) error {
	fileName := strings.ToLower(model.Name) + ".go"

	f := jen.NewFile("predicate")

	f.Type().Id(model.Name).Struct()

	if err := f.Save(path.Join(predicatePath, fileName)); err != nil {
		return err
	}

	return nil
}
