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
	edges         []*dbtype.Edge
	objects       []*dbtype.Object
	enums         []*dbtype.Enum
}

func newInput(source *parser.Output) (*input, error) {
	var in input

	in.sourcePkgPath = source.PkgPath

	getElement := func(name string) (dbtype.Element, bool) {
		for _, node := range in.nodes {
			if node.Name == name {
				return node, true
			}
		}

		for _, edge := range in.edges {
			if edge.Name == name {
				return edge, true
			}
		}

		for _, obj := range in.objects {
			if obj.Name == name {
				return obj, true
			}
		}

		// for _, enum := range in.enums {
		// 	if enum.Name == name {
		// 		return enum, true
		// 	}
		// }

		return nil, false
	}

	for _, node := range source.Nodes {
		dbNode := &dbtype.Node{
			Name: node.Name,
		}

		for _, f := range node.Fields {
			dbField, ok := field.Convert(f, getElement, strcase.ToSnake) // TODO
			if !ok {
				return nil, fmt.Errorf("could not convert field: %v", f)
			}
			dbNode.Fields = append(dbNode.Fields, dbField)
		}

		in.nodes = append(in.nodes, dbNode)
	}

	for _, edge := range source.Edges {
		dbEdge := &dbtype.Edge{
			Name: edge.Name,
		}

		inField, ok := field.Convert(edge.In, getElement, strcase.ToSnake) // TODO
		if !ok {
			return nil, fmt.Errorf("could not convert in field: %v", edge.In)
		}
		dbEdge.In = inField

		outField, ok := field.Convert(edge.Out, getElement, strcase.ToSnake) // TODO
		if !ok {
			return nil, fmt.Errorf("could not convert out field: %v", edge.Out)
		}
		dbEdge.Out = outField

		for _, f := range edge.Fields {
			dbField, ok := field.Convert(f, getElement, strcase.ToSnake) // TODO
			if !ok {
				return nil, fmt.Errorf("could not convert field: %v", f)
			}
			dbEdge.Fields = append(dbEdge.Fields, dbField)
		}

		in.edges = append(in.edges, dbEdge)
	}

	for _, str := range source.Structs {
		dbObject := &dbtype.Object{
			Name: str.Name,
		}

		for _, f := range str.Fields {
			dbField, ok := field.Convert(f, getElement, strcase.ToSnake) // TODO
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
