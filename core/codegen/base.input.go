package codegen

import (
	"fmt"
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
)

type input struct {
	sourcePkgPath string
	nodes         []*field.NodeTable
	edges         []*field.EdgeTable
	objects       []*field.DatabaseObject
	define        *parser.DefineOutput
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
	in.define = source.Define

	return &in, nil
}

func (in *input) SourceQual(name string) jen.Code {
	return jen.Qual(in.sourcePkgPath, name)
}

func (in *input) findNodeByName(name string) *field.NodeTable {
	for _, node := range in.nodes {
		if node.NameGo() == name {
			return node
		}
	}
	return nil
}

func fieldValueFrom(in *input, basePkg string, sf parser.ComplexIDField, accessor jen.Code) jen.Code {
	switch f := sf.Field.(type) {
	case *parser.FieldTime:
		return jen.Op("&").Qual(path.Join(basePkg, def.PkgTypes), "DateTime").Values(
			jen.Id("Time").Op(":").Add(accessor),
		)
	case *parser.FieldDuration:
		return jen.Op("&").Qual(path.Join(basePkg, def.PkgTypes), "Duration").Values(
			jen.Id("Duration").Op(":").Add(accessor),
		)
	case *parser.FieldNode:
		refNode := in.findNodeByName(f.Node)
		if refNode == nil {
			return accessor
		}
		tableName := refNode.NameDatabase()
		idVal := nodeRefValue(in, basePkg, refNode, accessor)
		return jen.Qual(def.PkgModels, "NewRecordID").Call(jen.Lit(tableName), idVal)
	default:
		return accessor
	}
}

func nodeRefValue(in *input, basePkg string, refNode *field.NodeTable, accessor jen.Code) jen.Code {
	if !refNode.HasComplexID() {
		return jen.String().Call(jen.Add(accessor).Dot("ID").Call())
	}
	cid := refNode.Source.ComplexID
	innerAccessor := jen.Add(accessor).Dot("ID").Call()
	if cid.Kind == parser.IDTypeArray {
		var elems []jen.Code
		for _, sf := range cid.Fields {
			elems = append(elems, fieldValueFrom(in, basePkg, sf, jen.Add(innerAccessor).Dot(sf.Name)))
		}
		return jen.Index().Any().Values(elems...)
	}
	dict := jen.Dict{}
	for _, sf := range cid.Fields {
		dict[jen.Lit(sf.DBName)] = fieldValueFrom(in, basePkg, sf, jen.Add(innerAccessor).Dot(sf.Name))
	}
	return jen.Map(jen.String()).Any().Values(dict)
}
