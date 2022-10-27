package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/parser"
	"path"
	"strings"
)

func buildBaseWhereFile(wherePath string) error {
	fileName := "base.go"

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

func buildFilterNodeFile(input *parser.Result, wherePath string, node parser.Node) error {
	return buildFilterFile(input, wherePath, node.Name, node.Fields, true)
}

func buildFilterStructFile(input *parser.Result, wherePath string, str parser.Struct) error {
	return buildFilterFile(input, wherePath, str.Name, str.Fields, false)
}

func buildFilterFile(input *parser.Result, wherePath string, name string, fields []parser.Field, isNode bool) error {
	prefix := "struct."
	if isNode {
		prefix = "node."
	}

	fileName := prefix + strcase.ToSnake(name) + ".go"

	f := jen.NewFile("where")

	if isNode {
		f.Var().Id(name).Op("=").Id("new" + name).Types(jen.Qual(input.PkgPath, name)).Call(jen.Lit(""))
	}

	f.Add(whereNew(input, name, fields))

	f.Type().Id(strings.ToLower(name)).
		Types(jen.Id("T").Any()).
		StructFunc(func(g *jen.Group) {
			for _, field := range fields {
				ok, code := whereField(input, field)
				if ok {
					g.Add(code)
				}
			}
		})

	f.Type().Id(strings.ToLower(name)+"Slice").
		Types(jen.Id("T").Any()).
		Struct(
			jen.Id(strings.ToLower(name)).Types(jen.Id("T")),
			jen.Op("*").Qual(pkgLibFilter, "Slice").Types(jen.Qual(input.PkgPath, name), jen.Id("T")),
		)

	for _, field := range fields {
		ok, code := whereFuncs(input, name, field)
		if ok {
			f.Add(code)
		}
	}

	if err := f.Save(path.Join(wherePath, fileName)); err != nil {
		return err
	}

	return nil
}

func whereNew(input *parser.Result, name string, fields []parser.Field) jen.Code {
	return jen.Func().Id("new" + name).
		Types(jen.Id("T").Any()).
		Params(jen.Id("key").String()).
		Id(strings.ToLower(name)).Types(jen.Id("T")).
		Block(
			jen.Return(
				jen.Id(strings.ToLower(name)).Types(jen.Id("T")).
					Values(jen.DictFunc(func(d jen.Dict) {
						for _, field := range fields {
							ok, key, value := whereFieldInit(input, field)
							if ok {
								d[key] = value
							}
						}
					})),
			),
		)
}

func whereFieldInit(input *parser.Result, field parser.Field) (bool, jen.Code, jen.Code) {
	typeNode := jen.Id("T")

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
	case parser.FieldTime:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewTime").Types(typeNode).Params(jen.Id("key"))
	case parser.FieldUUID:
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewBase").Types(jen.Qual(pkgUUID, "UUID"), typeNode).Params(jen.Id("key"))
	case parser.FieldEnum:
		typeEnum := jen.Qual(input.PkgPath, f.Typ)
		return true, jen.Id(f.Name), jen.Qual(pkgLibFilter, "NewBase").Types(typeEnum, typeNode).Params(jen.Id("key"))
	}

	return false, nil, nil
}

func whereField(input *parser.Result, field parser.Field) (bool, jen.Code) {
	typeNode := jen.Id("T")

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
	case parser.FieldTime:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Time").Types(typeNode)
	case parser.FieldUUID:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Base").Types(jen.Qual(pkgUUID, "UUID"), typeNode)
	case parser.FieldEnum:
		typeEnum := jen.Qual(input.PkgPath, f.Typ)
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibFilter, "Base").Types(typeEnum, typeNode)
	}

	return false, nil
}

func whereFuncs(input *parser.Result, name string, field parser.Field) (bool, jen.Code) {
	switch f := field.(type) {
	case parser.FieldNode:
		return true, jen.Func().
			Params(jen.Id(strings.ToLower(name)).Types(jen.Id("T"))).
			Id(f.Name).Params().
			Id(strings.ToLower(f.Node)).Types(jen.Id("T")).
			Block(
				jen.Return(jen.Id("new" + f.Node).Types(jen.Id("T")).Params(jen.Lit(strcase.ToSnake(f.Name)))),
			)
	case parser.FieldStruct:
		return true, jen.Func().
			Params(jen.Id(strings.ToLower(name)).Types(jen.Id("T"))).
			Id(f.Name).Params().
			Id(strings.ToLower(f.Struct)).Types(jen.Id("T")).
			Block(
				jen.Return(jen.Id("new" + f.Struct).Types(jen.Id("T")).Params(jen.Lit(strcase.ToSnake(f.Name)))),
			)
	case parser.FieldSlice:
		if f.IsNode {
			return true, jen.Func().
				Params(jen.Id(strings.ToLower(name)).Types(jen.Id("T"))).
				Id(f.Name).Params().
				Id(strings.ToLower(f.Value) + "Slice").Types(jen.Id("T")).
				Block(
					jen.Return(
						jen.Id(strings.ToLower(f.Value)+"Slice").Types(jen.Id("T")).
							Values(
								jen.Id("new"+f.Value).Types(jen.Id("T")).Call(jen.Lit(strcase.ToSnake(f.Name))),
								jen.Qual(pkgLibFilter, "NewSlice").Types(jen.Qual(input.PkgPath, f.Value), jen.Id("T")).Call(jen.Lit(strcase.ToSnake(f.Name))),
							),
					),
				)
		} else {
			return true, jen.Func().
				Params(jen.Id(strings.ToLower(name)).Types(jen.Id("T"))).
				Id(f.Name).Params().
				Block()
		}
	case parser.FieldMap:
		return true, jen.Func().
			Params(jen.Id(strings.ToLower(name)).Types(jen.Id("T"))).
			Id(f.Name).Params().
			Block()
	}

	return false, nil
}
