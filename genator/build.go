package genator

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
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
	dirConv      = "conv"
)

const (
	pkgLibFilter  = "github.com/marcbinz/sdb/lib/filter"
	pkgLibSort    = "github.com/marcbinz/sdb/lib/sort"
	pkgLibBuilder = "github.com/marcbinz/sdb/lib/builder"

	pkgBase      = "github.com/marcbinz/sdb/example/gen/sdb"           // TODO
	pkgQuery     = "github.com/marcbinz/sdb/example/gen/sdb/query"     // TODO
	pkgPredicate = "github.com/marcbinz/sdb/example/gen/sdb/predicate" // TODO
	pkgConv      = "github.com/marcbinz/sdb/example/gen/sdb/conv"      // TODO

	pkgUUID      = "github.com/google/uuid"
	pkgSurrealDB = "github.com/surrealdb/surrealdb.go"
)

func Build(input *parser.Result, outDir string) error {
	basePath := path.Join(outDir, dirBase)
	queryPath := path.Join(basePath, dirQuery)
	wherePath := path.Join(basePath, dirWhere)
	predicatePath := path.Join(basePath, dirPredicate)
	byPath := path.Join(basePath, dirBy)
	convPath := path.Join(basePath, dirConv)

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

	err = os.MkdirAll(convPath, os.ModePerm)
	if err != nil {
		return err
	}

	if err := buildSdbFile(basePath); err != nil {
		return err
	}

	if err := buildBaseWhereFile(wherePath); err != nil {
		return err
	}

	if err := buildBaseConvFile(convPath); err != nil {
		return err
	}

	for _, node := range input.Nodes {

		if err := baseFile(basePath, node, input.PkgPath); err != nil {
			return err
		}

		if err := buildQueryFile(input, queryPath, node); err != nil {
			return err
		}

		if err := buildFilterNodeFile(input, wherePath, node); err != nil {
			return err
		}

		if err := predicateFile(predicatePath, node, input.PkgPath); err != nil {
			return err
		}

		if err := buildByFile(input, byPath, node); err != nil {
			return err
		}

		if err := buildConvFile(input, convPath, node.Name, node.Fields, true); err != nil {
			return err
		}
	}

	for _, str := range input.Structs {
		if err := buildFilterStructFile(input, wherePath, str); err != nil {
			return err
		}

		if err := buildConvFile(input, convPath, str.Name, str.Fields, false); err != nil {
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

	f.Func().Id("NewClient").
		Params(jen.Id("db").Op("*").Qual(pkgSurrealDB, "DB")).
		Op("*").Id("Client").
		Block(
			jen.Return(jen.Id("&").Id("Client").Values(
				jen.Dict{
					jen.Id("db"): jen.Id("db"),
				},
			)),
		)

	f.Func().Params(jen.Id("c").Op("*").Id("Client")).
		Id("Create").Params(jen.Id("node").String(), jen.Id("data").Map(jen.String()).Any()).
		Params(jen.Any(), jen.Error()).
		Block(
			jen.Return(
				jen.Id("c").Dot("db").Dot("Create").Call(jen.Id("node"), jen.Id("data")),
			),
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
			jen.Return(jen.Qual(pkgQuery, "New"+model.Name).Call()),
		)

	f.Func().
		Params(jen.Id(strings.ToLower(model.Name))).
		Id("Create").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("db").Op("*").Id("Client"),
			jen.Id(strings.ToLower(model.Name)).Op("*").Qual(modelPkg, model.Name),
		).
		Error().
		Block(
			jen.If(jen.Id(strings.ToLower(model.Name)).Dot("ID").Op("!=").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID must not be set for a node to be created"))),
				),
			jen.Id("data").Op(":=").Qual(pkgConv, "From"+model.Name).Call(jen.Op("*").Id(strings.ToLower(model.Name))),
			jen.Id("raw").Op(",").Err().Op(":=").Id("db").Dot("Create").Call(jen.Lit(strcase.ToSnake(model.Name)), jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.Id("res").Op(":=").Qual(pkgConv, "To"+model.Name).
				Call(jen.Id("raw").Op(".").Parens(jen.Index().Any()).Index(jen.Lit(0)).Op(".").Parens(jen.Map(jen.String()).Any())),
			jen.Qual("fmt", "Println").Call(jen.Id("res")),
			jen.Return(jen.Nil()),
		)

	// if user.ID != "" {
	//		return errors.New("ID must not be set")
	//	}
	//
	//	data := conv.FromUser(*user)
	//	raw, err := db.Create("user", data)
	//	if err != nil {
	//		return err
	//	}
	//
	//	res := conv.ToUser(raw.([]any)[0].(map[string]any))
	//
	//	fmt.Println(res)
	//	return nil

	f.Func().
		Params(jen.Id(strings.ToLower(model.Name))).
		Id("Read").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("db").Op("*").Id("Client"),
			jen.Id("id").String(),
		).
		Params(jen.Op("*").Qual(modelPkg, model.Name), jen.Error()).
		Block( // TODO: what if id already contains "user:" ?!
			jen.Id("raw").Op(",").Err().Op(":=").Id("db").Dot("db").Dot("Select").Call(jen.Lit(strcase.ToSnake(model.Name)+":").Op("+").Id("id")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Id("res").Op(":=").Qual(pkgConv, "To"+model.Name).Call(jen.Id("raw").Op(".").Parens(jen.Map(jen.String()).Any())),
			jen.Return(jen.Op("&").Id("res"), jen.Nil()),
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
