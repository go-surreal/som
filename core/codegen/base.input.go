package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/field"
	"github.com/marcbinz/som/core/parser"
)

type input struct {
	sourcePkgPath string
	nodes         []*field.NodeTable
	edges         []*field.EdgeTable
	objects       []*field.DatabaseObject
	enums         []*field.DatabaseEnum
}

func newInput(source *parser.Output) (*input, error) {
	buildConf := &field.BuildConfig{
		SourcePkg:      source.PkgPath,
		ToDatabaseName: strcase.ToSnake, // TODO
	}

	var in input

	in.sourcePkgPath = source.PkgPath

	// getElement := func(name string) (field.Element, bool) {
	// 	for _, node := range in.nodes {
	// 		if node.Name == name {
	// 			return node, true
	// 		}
	// 	}
	//
	// 	for _, edge := range in.edges {
	// 		if edge.Name == name {
	// 			return edge, true
	// 		}
	// 	}
	//
	// 	for _, obj := range in.objects {
	// 		if obj.Name == name {
	// 			return obj, true
	// 		}
	// 	}
	//
	// 	// for _, enum := range in.enums {
	// 	// 	if enum.Name == name {
	// 	// 		return enum, true
	// 	// 	}
	// 	// }
	//
	// 	return nil, false
	// }

	def, err := field.NewDef(source, buildConf)
	if err != nil {
		return nil, fmt.Errorf("could not build def: %v", err)
	}

	in.nodes = def.Nodes
	in.edges = def.Edges
	in.objects = def.Objects
	in.enums = def.Enums

	return &in, nil
}

func (in *input) SourceQual(name string) jen.Code {
	return jen.Qual(in.sourcePkgPath, name)
}
