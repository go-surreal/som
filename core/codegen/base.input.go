package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
)

type input struct {
	sourcePkgPath string
	nodes         []*field.NodeTable
	edges         []*field.EdgeTable
	objects       []*field.DatabaseObject
}

func newInput(source *parser.Output, outPkg string) (*input, error) {
	buildConf := &field.BuildConfig{
		SourcePkg:      source.PkgPath,
		TargetPkg:      outPkg,
		ToDatabaseName: strcase.ToSnake, // TODO
	}

	var in input

	in.sourcePkgPath = source.PkgPath

	def, err := field.NewDef(source, buildConf)
	if err != nil {
		return nil, fmt.Errorf("could not build def: %w", err)
	}

	in.nodes = def.Nodes
	in.edges = def.Edges
	in.objects = def.Objects

	return &in, nil
}

func (in *input) SourceQual(name string) jen.Code {
	return jen.Qual(in.sourcePkgPath, name)
}
