package builder

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/parser"
	"path"
	"strings"
)

func buildByFile(input *parser.Result, byPath string, model parser.Node) error {
	fileName := strings.ToLower(model.Name) + ".go"

	f := jen.NewFile("by")

	f.Var().Id(model.Name).Op("=").Id("new" + model.Name).Call(jen.Lit(""))

	f.Add(byNew(input, model))

	f.Type().Id(strings.ToLower(model.Name)).StructFunc(func(g *jen.Group) {
		for _, field := range model.Fields {
			ok, code := byField(input, model, field)
			if ok {
				g.Add(code)
			}
		}
	})

	f.Func().Params(jen.Id(strings.ToLower(model.Name))).
		Id("Random").Params().
		Op("*").Qual(pkgLibSort, "Of").Types(jen.Qual(input.PkgPath, model.Name)).
		Block(
			jen.Return(jen.Nil()),
		)

	if err := f.Save(path.Join(byPath, fileName)); err != nil {
		return err
	}

	return nil
}

func byNew(input *parser.Result, model parser.Node) jen.Code {
	return jen.Func().Id("new" + model.Name).
		Params(jen.Id("key").String()).
		Id(strings.ToLower(model.Name)).
		Block(
			jen.Return(
				jen.Id(strings.ToLower(model.Name)).Values(jen.DictFunc(func(d jen.Dict) {
					for _, field := range model.Fields {
						ok, key, value := byFieldInit(input, model, field)
						if ok {
							d[key] = value
						}
					}
				})),
			),
		)
}

func byFieldInit(input *parser.Result, node parser.Node, field parser.Field) (bool, jen.Code, jen.Code) {
	switch f := field.(type) {
	case parser.FieldID:
		return true, jen.Id(f.Name), jen.Qual(pkgLibSort, "NewSort").Types(jen.Qual(input.PkgPath, node.Name)).Params(jen.Id("key"))
	case parser.FieldString:
		return true, jen.Id(f.Name), jen.Qual(pkgLibSort, "NewString").Types(jen.Qual(input.PkgPath, node.Name)).Params(jen.Id("key"))
	case parser.FieldInt:
		return true, jen.Id(f.Name), jen.Qual(pkgLibSort, "NewSort").Types(jen.Qual(input.PkgPath, node.Name)).Params(jen.Id("key"))
	case parser.FieldInt32:
		return true, jen.Id(f.Name), jen.Qual(pkgLibSort, "NewSort").Types(jen.Qual(input.PkgPath, node.Name)).Params(jen.Id("key"))
	case parser.FieldInt64:
		return true, jen.Id(f.Name), jen.Qual(pkgLibSort, "NewSort").Types(jen.Qual(input.PkgPath, node.Name)).Params(jen.Id("key"))
	case parser.FieldFloat32:
		return true, jen.Id(f.Name), jen.Qual(pkgLibSort, "NewSort").Types(jen.Qual(input.PkgPath, node.Name)).Params(jen.Id("key"))
	case parser.FieldFloat64:
		return true, jen.Id(f.Name), jen.Qual(pkgLibSort, "NewSort").Types(jen.Qual(input.PkgPath, node.Name)).Params(jen.Id("key"))
	case parser.FieldTime:
		return true, jen.Id(f.Name), jen.Qual(pkgLibSort, "NewSort").Types(jen.Qual(input.PkgPath, node.Name)).Params(jen.Id("key"))
	}

	return false, nil, nil
}

func byField(input *parser.Result, node parser.Node, field parser.Field) (bool, jen.Code) {
	switch f := field.(type) {
	case parser.FieldID:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibSort, "Sort").Types(jen.Qual(input.PkgPath, node.Name))
	case parser.FieldString:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibSort, "String").Types(jen.Qual(input.PkgPath, node.Name))
	case parser.FieldInt:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibSort, "Sort").Types(jen.Qual(input.PkgPath, node.Name))
	case parser.FieldInt32:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibSort, "Sort").Types(jen.Qual(input.PkgPath, node.Name))
	case parser.FieldInt64:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibSort, "Sort").Types(jen.Qual(input.PkgPath, node.Name))
	case parser.FieldFloat32:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibSort, "Sort").Types(jen.Qual(input.PkgPath, node.Name))
	case parser.FieldFloat64:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibSort, "Sort").Types(jen.Qual(input.PkgPath, node.Name))
	case parser.FieldTime:
		return true, jen.Id(f.Name).Op("*").Qual(pkgLibSort, "Sort").Types(jen.Qual(input.PkgPath, node.Name))
	}

	return false, nil
}
