package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/parser"
	"path"
	"strings"
)

func buildBaseWhereFile(wherePath string) error {
	fileName := "0_where.go"

	f := jen.NewFile("where")

	f.Func().Id("All").
		Types(jen.Id("T").Any()).
		Params(jen.Id("filters").Op("...").Id("T")).
		Id("T").
		Block(
			jen.Return(jen.Nil()),
		)

	f.Func().Id("Any").
		Types(jen.Id("T").Any()).
		Params(jen.Id("filters").Op("...").Id("T")).
		Id("T").
		Block(
			jen.Return(jen.Nil()),
		)

	f.Func().Id("Count").
		Types(jen.Id("T").Any()).
		Params(jen.Id("what").Id("T")).
		Id("T").
		Block(
			jen.Return(jen.Nil()),
		)

	if err := f.Save(path.Join(wherePath, fileName)); err != nil {
		return err
	}

	return nil
}

func buildWhereFile(input *parser.Result, wherePath string, model parser.Node, modelPkg string) error {
	fileName := strings.ToLower(model.Name) + ".go"

	f := jen.NewFile("where")

	f.Var().Id(model.Name).Op("=").Id("new" + model.Name).Call(jen.Lit(""))

	f.Add(whereNew(input, model))

	f.Type().Id(strings.ToLower(model.Name)).StructFunc(func(g *jen.Group) {
		for _, field := range model.Fields {
			ok, code := whereField(input, model, field)
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

func whereNew(input *parser.Result, model parser.Node) jen.Code {
	return jen.Func().Id("new" + model.Name).
		Params(jen.Id("key").String()).
		Id(strings.ToLower(model.Name)).
		Block(
			jen.Return(
				jen.Id(strings.ToLower(model.Name)).Values(jen.DictFunc(func(d jen.Dict) {
					for _, field := range model.Fields {
						ok, key, value := whereFieldInit(input, model, field)
						if ok {
							d[key] = value
						}
					}
				})),
			),
		)
}

func whereFieldInit(input *parser.Result, node parser.Node, field parser.Field) (bool, jen.Code, jen.Code) {
	typeNode := jen.Qual(input.PkgPath, node.Name)

	switch f := field.(type) {
	case parser.FieldID:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewBase").Types(jen.String(), typeNode).Params(jen.Id("key"))
	case parser.FieldString:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewString").Types(typeNode).Params(jen.Id("key"))
	case parser.FieldInt:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewNumeric").Types(jen.Int(), typeNode).Params(jen.Id("key"))
	case parser.FieldInt32:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewNumeric").Types(jen.Int32(), typeNode).Params(jen.Id("key"))
	case parser.FieldInt64:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewNumeric").Types(jen.Int64(), typeNode).Params(jen.Id("key"))
	case parser.FieldFloat32:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewNumeric").Types(jen.Float32(), typeNode).Params(jen.Id("key"))
	case parser.FieldFloat64:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewNumeric").Types(jen.Float64(), typeNode).Params(jen.Id("key"))
	case parser.FieldBool:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewBool").Types(typeNode).Params(jen.Id("key"))
	case parser.FieldEnum:
		typeEnum := jen.Qual(input.PkgPath, f.Typ)
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewBase").Types(typeEnum, typeNode).Params(jen.Id("key"))
	}

	return false, nil, nil
}

func whereField(input *parser.Result, node parser.Node, field parser.Field) (bool, jen.Code) {
	typeNode := jen.Qual(input.PkgPath, node.Name)

	switch f := field.(type) {
	case parser.FieldID:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Base").Types(jen.String(), typeNode)
	case parser.FieldString:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "String").Types(typeNode)
	case parser.FieldInt:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Numeric").Types(jen.Int(), typeNode)
	case parser.FieldInt32:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Numeric").Types(jen.Int32(), typeNode)
	case parser.FieldInt64:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Numeric").Types(jen.Int64(), typeNode)
	case parser.FieldFloat32:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Numeric").Types(jen.Float32(), typeNode)
	case parser.FieldFloat64:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Numeric").Types(jen.Float64(), typeNode)
	case parser.FieldBool:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Bool").Types(typeNode)
	case parser.FieldEnum:
		typeEnum := jen.Qual(input.PkgPath, f.Typ)
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Base").Types(typeEnum, typeNode)
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
