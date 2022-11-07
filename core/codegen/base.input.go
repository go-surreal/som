package codegen

import (
	"errors"
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/dbtype"
	"github.com/marcbinz/sdb/core/codegen/field"
	"github.com/marcbinz/sdb/core/parser"
)

type input struct {
	sourcePkgPath string
	nodes         []*dbtype.Node
	objects       []*dbtype.Object
	enums         []*dbtype.Enum
}

func newInput(source *parser.Output) (*input, error) {
	var in input

	in.sourcePkgPath = source.PkgPath

	for _, node := range source.Nodes {
		dbNode := &dbtype.Node{
			Name: node.Name,
		}

		for _, f := range node.Fields {
			dbField, ok := field.Convert(f, strcase.ToSnake) // TODO
			if !ok {
				return nil, fmt.Errorf("could not convert field: %v", f)
			}
			dbNode.Fields = append(dbNode.Fields, dbField)
		}

		in.nodes = append(in.nodes, dbNode)
	}

	for _, str := range source.Structs {
		dbObject := &dbtype.Object{
			Name: str.Name,
		}

		for _, f := range str.Fields {
			dbField, ok := field.Convert(f, strcase.ToSnake) // TODO
			if !ok {
				return nil, errors.New("could not convert field")
			}
			dbObject.Fields = append(dbObject.Fields, dbField)
		}

		in.objects = append(in.objects, dbObject)
	}

	for _, enum := range source.Enums {
		dbEnum := &dbtype.Enum{
			Name: enum.Name,
		}

		in.enums = append(in.enums, dbEnum)
	}

	return &in, nil
}

func (in *input) isEnum(name string) bool {
	for _, enum := range in.enums {
		if enum.Name == name {
			return true
		}
	}
	return false
}

func (in *input) SourceQual(name string) jen.Code {
	return jen.Qual(in.sourcePkgPath, name)
}
