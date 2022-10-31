package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/parser"
	"path"
)

func buildBaseConvFile(convPath string) error {
	fileName := "conv.go"

	f := jen.NewFile("conv")

	f.Func().Id("prepareID").
		Params(jen.Id("node").String(), jen.Id("id").Any()).
		String().
		Block(
			jen.Return(
				jen.Qual("strings", "TrimPrefix").Call(
					jen.Id("id").Op(".").Parens(jen.String()),
					jen.Id("node").Op("+").Lit(":"),
				),
			),
		)

	// strings.TrimPrefix(data["id"].(string), "user:")

	f.Func().Id("parseTime").Params(jen.Id("val").Any()).Qual("time", "Time").
		Block(
			jen.Id("res").Op(",").Err().Op(":=").
				Qual("time", "Parse").Call(jen.Qual("time", "RFC3339"), jen.Id("val").Op(".").Parens(jen.String())),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("time", "Time").Values()),
			),
			jen.Return(jen.Id("res")),
		)

	f.Func().Id("parseUUID").Params(jen.Id("val").Any()).Qual(pkgUUID, "UUID").
		Block(
			jen.Id("res").Op(",").Err().Op(":=").
				Qual(pkgUUID, "Parse").Call(jen.Id("val").Op(".").Parens(jen.String())),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual(pkgUUID, "UUID").Values()),
			),
			jen.Return(jen.Id("res")),
		)

	// func extract[T any](val any, to func(map[string]any) T) T {
	//	var t T
	//	if val == nil {
	//		return t
	//	}
	//	return to(val.(map[string]any))
	// }

	if err := f.Save(path.Join(convPath, fileName)); err != nil {
		return err
	}

	return nil
}

func buildConvFile(input *parser.Result, convPath string, name string, fields []parser.Field, node *parser.Node) error {
	prefix := "struct."
	if node != nil {
		prefix = "node."
	}

	fileName := prefix + strcase.ToSnake(name) + ".go"

	f := jen.NewFile("conv")

	f.Add(buildConvFromModel(input, name, fields))
	f.Add(buildConvToModel(input, name, node, fields))

	if err := f.Save(path.Join(convPath, fileName)); err != nil {
		return err
	}

	return nil
}

func buildConvFromModel(input *parser.Result, name string, fields []parser.Field) jen.Code {
	return jen.Func().
		Id("From" + name).
		Params(jen.Id("data").Qual(input.PkgPath, name)).
		Map(jen.String()).Any().
		Block(
			jen.Return(jen.Map(jen.String()).Any().Values(jen.DictFunc(func(d jen.Dict) {
				for _, field := range fields {
					ok, key, value := buildConvFromField(field)
					if ok {
						d[key] = value
					}
				}
			}))))
}

func buildConvToModel(input *parser.Result, name string, node *parser.Node, fields []parser.Field) jen.Code {
	return jen.Func().
		Id("To"+name).
		Params(jen.Id("data").Map(jen.String()).Any()).
		Qual(input.PkgPath, name).
		Block(
			jen.Return(jen.Qual(input.PkgPath, name).Values(jen.DictFunc(func(d jen.Dict) {
				for _, field := range fields {
					ok, key, value := buildConvToField(field, node)
					if ok {
						d[key] = value
					}
				}
			}))))
}

func buildConvFromField(field parser.Field) (bool, jen.Code, jen.Code) {
	switch f := field.(type) {
	case parser.FieldID:
		// ID is never set, because the database is providing them, not the application.
		return false, nil, nil
	case parser.FieldString:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("data").Dot(f.Name)
	case parser.FieldInt:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("data").Dot(f.Name)
	case parser.FieldInt32:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("data").Dot(f.Name)
	case parser.FieldInt64:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("data").Dot(f.Name)
	case parser.FieldFloat32:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("data").Dot(f.Name)
	case parser.FieldFloat64:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("data").Dot(f.Name)
	case parser.FieldBool:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("data").Dot(f.Name)
	case parser.FieldTime:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("data").Dot(f.Name)
	case parser.FieldUUID:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("data").Dot(f.Name)
	case parser.FieldNode:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("From" + f.Node).Call(jen.Id("data").Dot(f.Name))
	case parser.FieldStruct:
		return true, jen.Lit(strcase.ToSnake(f.Name)), jen.Id("From" + f.Struct).Call(jen.Id("data").Dot(f.Name))
	}

	return false, nil, nil
}

func buildConvToField(field parser.Field, node *parser.Node) (bool, jen.Code, jen.Code) {
	switch f := field.(type) {
	case parser.FieldID:
		return true, jen.Id(f.Name), jen.Id("prepareID").Call(jen.Lit(strcase.ToSnake(node.Name)), jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))))
	case parser.FieldString:
		return true, jen.Id(f.Name), jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))).Op(".").Parens(jen.String())
	case parser.FieldInt:
		return true, jen.Id(f.Name), jen.Int().Call(jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))).Op(".").Parens(jen.Float64()))
	case parser.FieldInt32:
		return true, jen.Id(f.Name), jen.Int32().Call(jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))).Op(".").Parens(jen.Float64()))
	case parser.FieldInt64:
		return true, jen.Id(f.Name), jen.Int64().Call(jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))).Op(".").Parens(jen.Float64()))
	case parser.FieldFloat32:
		return true, jen.Id(f.Name), jen.Float32().Call(jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))).Op(".").Parens(jen.Float64()))
	case parser.FieldFloat64:
		return true, jen.Id(f.Name), jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))).Op(".").Parens(jen.Float64())
	case parser.FieldBool:
		return true, jen.Id(f.Name), jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))).Op(".").Parens(jen.Bool())
	case parser.FieldTime:
		return true, jen.Id(f.Name), jen.Id("parseTime").Call(jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))))
	case parser.FieldUUID:
		return true, jen.Id(f.Name), jen.Id("parseUUID").Call(jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))))
	case parser.FieldNode:
		return true, jen.Id(f.Name), jen.Id("To" + f.Node).Call(jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))).Op(".").Parens(jen.Map(jen.String()).Any()))
	case parser.FieldStruct:
		return true, jen.Id(f.Name), jen.Id("To" + f.Struct).Call(jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.Name))).Op(".").Parens(jen.Map(jen.String()).Any()))
	}

	return false, nil, nil
}
