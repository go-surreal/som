package field

import (
	"errors"
	"fmt"
	"github.com/marcbinz/som/core/parser"
)

type Def struct {
	Nodes   []*NodeTable
	Edges   []*EdgeTable
	Objects []*DatabaseObject
	Enums   []*DatabaseEnum
}

func NewDef(source *parser.Output, buildConf *BuildConfig) (*Def, error) {
	var def Def

	for _, node := range source.Nodes {
		dbNode := &NodeTable{
			Name: node.Name,
		}

		for _, f := range node.Fields {
			dbField, ok := Convert(buildConf, f)
			if !ok {
				return nil, fmt.Errorf("could not convert field: %v", f)
			}
			dbNode.Fields = append(dbNode.Fields, dbField)
		}

		def.Nodes = append(def.Nodes, dbNode)
	}

	for _, edge := range source.Edges {
		dbEdge := &EdgeTable{
			Name: edge.Name,
		}

		inField, ok := Convert(buildConf, edge.In)
		if !ok {
			return nil, fmt.Errorf("could not convert in field: %v", edge.In)
		}
		in := inField.(*Node)
		dbEdge.In = *in

		outField, ok := Convert(buildConf, edge.Out)
		if !ok {
			return nil, fmt.Errorf("could not convert out field: %v", edge.Out)
		}
		out := outField.(*Node)
		dbEdge.Out = *out

		for _, f := range edge.Fields {
			dbField, ok := Convert(buildConf, f)
			if !ok {
				return nil, fmt.Errorf("could not convert field: %v", f)
			}
			dbEdge.Fields = append(dbEdge.Fields, dbField)
		}

		def.Edges = append(def.Edges, dbEdge)
	}

	for _, str := range source.Structs {
		dbObject := &DatabaseObject{
			Name: str.Name,
		}

		for _, f := range str.Fields {
			dbField, ok := Convert(buildConf, f)
			if !ok {
				return nil, errors.New("could not convert field")
			}
			dbObject.Fields = append(dbObject.Fields, dbField)
		}

		def.Objects = append(def.Objects, dbObject)
	}

	for _, enum := range source.Enums {
		dbEnum := &DatabaseEnum{
			Name: enum.Name,
		}

		def.Enums = append(def.Enums, dbEnum)
	}

	return &def, nil
}

func Convert(conf *BuildConfig, field parser.Field) (Field, bool) {
	base := &baseField{BuildConfig: conf, source: field}

	switch f := field.(type) {

	case *parser.FieldID:
		return &ID{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldString:
		return &String{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldNumeric:
		return &Numeric{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldBool:
		return &Bool{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldTime:
		return &Time{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldUUID:
		return &UUID{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldEnum:
		return &Enum{
			baseField: base,
			source:    f,
			model:     Model(f.Typ),
		}, true

	case *parser.FieldStruct:
		return &Struct{
			baseField: base,
			source:    f,
		}, true

	case *parser.FieldNode:
		return &Node{
			baseField: base,
			source:    f,
			table: NodeTable{
				Name:   "NodeTable",
				Fields: nil,
			},
		}, true

	case *parser.FieldEdge:
		// in,_ := Convert(conf, f.Edge)
		return &Edge{
			baseField: base,
			source:    f,
			table: EdgeTable{
				Name:   "EdgeTable",
				In:     Node{},
				Out:    Node{},
				Fields: nil,
			}, // TODO
		}, true

	case *parser.FieldSlice:
		element, ok := Convert(conf, f.Field)
		if !ok {
			return nil, false
		}

		return &Slice{
			baseField: base,
			source:    f,
			element:   element,
		}, true
	}

	return nil, false
}
