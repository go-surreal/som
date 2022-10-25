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
)

const (
	pkgLib = "github.com/marcbinz/sdb/lib"

	pkgQuery = "go.alfnet.dev/service/gampi/gen/sdb/query" // TODO
)

func Build(input *parser.Result, outDir string) error {
	basePath := path.Join(outDir, dirBase)
	queryPath := path.Join(basePath, dirQuery)
	wherePath := path.Join(basePath, dirWhere)
	predicatePath := path.Join(basePath, dirPredicate)

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

	for _, model := range input.Nodes {

		if err := baseFile(basePath, model, input.PkgPath); err != nil {
			return err
		}

		if err := queryFile(queryPath, model, input.PkgPath); err != nil {
			return err
		}

		if err := whereFile(wherePath, model, input.PkgPath); err != nil {
			return err
		}

		if err := predicateFile(predicatePath, model, input.PkgPath); err != nil {
			return err
		}
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
// -- QUERY
//

func queryFile(queryPath string, model parser.Node, modelPkg string) error {
	fileName := strings.ToLower(model.Name) + ".go"

	f := jen.NewFile("query")

	f.Type().Id(model.Name).Struct()

	if err := f.Save(path.Join(queryPath, fileName)); err != nil {
		return err
	}

	return nil
}

//
// -- WHERE
//

func whereFile(wherePath string, model parser.Node, modelPkg string) error {
	fileName := strings.ToLower(model.Name) + ".go"

	f := jen.NewFile("where")

	f.Var().Id(model.Name).Op("=").Id("new" + model.Name).Call(jen.Lit(""))

	f.Add(whereNew(model))

	f.Type().Id(strings.ToLower(model.Name)).StructFunc(func(g *jen.Group) {
		for _, field := range model.Fields {
			ok, code := whereField(field)
			if ok {
				g.Add(code)
			}
		}
	})

	for _, field := range model.Fields {
		ok, code := whereFuncs(model, field)
		if ok {
			f.Add(code)
		}
	}

	if err := f.Save(path.Join(wherePath, fileName)); err != nil {
		return err
	}

	return nil
}

func whereNew(model parser.Node) jen.Code {
	return jen.Func().Id("new" + model.Name).
		Params(jen.Id("origin").String()).
		Id(strings.ToLower(model.Name)).
		Block(
			jen.Return(
				jen.Id(strings.ToLower(model.Name)).Values(jen.DictFunc(func(d jen.Dict) {
					for _, field := range model.Fields {
						ok, key, value := whereFieldInit(field)
						if ok {
							d[key] = value
						}
					}
				})),
			),
		)
}

func whereFieldInit(field parser.Field) (bool, jen.Code, jen.Code) {
	switch f := field.(type) {
	case parser.FieldID:
		return true, jen.Id(f.Name), jen.Qual(pkgLib, "WhereID").Values(jen.Id("origin"))
	case parser.FieldString:
		return true, jen.Id(f.Name), jen.Qual(pkgLib, "WhereString").Values(jen.Id("origin"))
	case parser.FieldInt:
		return true, jen.Id(f.Name), jen.Qual(pkgLib, "WhereInt").Values(jen.Id("origin"))
	case parser.FieldInt32:
		return true, jen.Id(f.Name), jen.Qual(pkgLib, "WhereInt32").Values(jen.Id("origin"))
	case parser.FieldInt64:
		return true, jen.Id(f.Name), jen.Qual(pkgLib, "WhereInt64").Values(jen.Id("origin"))
	case parser.FieldFloat32:
		return true, jen.Id(f.Name), jen.Qual(pkgLib, "WhereFloat32").Values(jen.Id("origin"))
	case parser.FieldFloat64:
		return true, jen.Id(f.Name), jen.Qual(pkgLib, "WhereFloat64").Values(jen.Id("origin"))
	case parser.FieldBool:
		return true, jen.Id(f.Name), jen.Qual(pkgLib, "WhereBool").Values(jen.Id("origin"))
	}

	return false, nil, nil
}

func whereField(field parser.Field) (bool, jen.Code) {
	switch f := field.(type) {
	case parser.FieldID:
		return true, jen.Id(f.Name).Qual(pkgLib, "WhereID")
	case parser.FieldString:
		return true, jen.Id(f.Name).Qual(pkgLib, "WhereString")
	case parser.FieldInt:
		return true, jen.Id(f.Name).Qual(pkgLib, "WhereInt")
	case parser.FieldInt32:
		return true, jen.Id(f.Name).Qual(pkgLib, "WhereInt32")
	case parser.FieldInt64:
		return true, jen.Id(f.Name).Qual(pkgLib, "WhereInt64")
	case parser.FieldFloat32:
		return true, jen.Id(f.Name).Qual(pkgLib, "WhereFloat32")
	case parser.FieldFloat64:
		return true, jen.Id(f.Name).Qual(pkgLib, "WhereFloat64")
	case parser.FieldBool:
		return true, jen.Id(f.Name).Qual(pkgLib, "WhereBool")
	}

	return false, nil
}

func whereFuncs(model parser.Node, field parser.Field) (bool, jen.Code) {
	switch f := field.(type) {
	case parser.FieldStruct:
		return true, jen.Func().
			Params(jen.Id(strings.ToLower(model.Name))).
			Id(f.Name).Params().
			Block()
	case parser.FieldSlice:
		return true, jen.Func().
			Params(jen.Id(strings.ToLower(model.Name))).
			Id(f.Name).Params().
			Block()
	case parser.FieldMap:
		return true, jen.Func().
			Params(jen.Id(strings.ToLower(model.Name))).
			Id(f.Name).Params().
			Block()
	}

	return false, nil
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
