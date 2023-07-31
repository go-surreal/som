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
}

func NewDef(source *parser.Output, buildConf *BuildConfig) (*Def, error) {
	var def Def

	for _, node := range source.Nodes {
		dbNode := &NodeTable{
			Name:       node.Name,
			Timestamps: node.Timestamps,
		}

		for _, f := range node.Fields {
			dbField, ok := Convert(source, buildConf, f)
			if !ok {
				return nil, fmt.Errorf("could not convert field: %v", f)
			}
			dbNode.Fields = append(dbNode.Fields, dbField)
		}

		def.Nodes = append(def.Nodes, dbNode)
	}

	for _, edge := range source.Edges {
		dbEdge := &EdgeTable{
			Name:       edge.Name,
			Timestamps: edge.Timestamps,
		}

		inField, ok := Convert(source, buildConf, edge.In)
		if !ok {
			return nil, fmt.Errorf("could not convert in field: %v", edge.In)
		}
		dbEdge.In = inField.(*Node)

		outField, ok := Convert(source, buildConf, edge.Out)
		if !ok {
			return nil, fmt.Errorf("could not convert out field: %v", edge.Out)
		}
		dbEdge.Out = outField.(*Node)

		for _, f := range edge.Fields {
			dbField, ok := Convert(source, buildConf, f)
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
			dbField, ok := Convert(source, buildConf, f)
			if !ok {
				return nil, errors.New("could not convert field")
			}
			dbObject.Fields = append(dbObject.Fields, dbField)
		}

		def.Objects = append(def.Objects, dbObject)
	}

	return &def, nil
}

func Convert(source *parser.Output, conf *BuildConfig, field parser.Field) (Field, bool) {
	base := &baseField{BuildConfig: conf, source: field}

	switch f := field.(type) {

	case *parser.FieldID:
		{
			return &ID{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldString:
		{
			return &String{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldNumeric:
		{
			return &Numeric{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldBool:
		{
			return &Bool{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldStruct:
		{
			var object *parser.Struct
			for _, elem := range source.Structs {
				if elem.Name == f.Struct {
					object = elem
					break
				}
			}

			if object == nil {
				return nil, false
			}

			var fields []Field
			for _, field := range object.Fields {
				fld, ok := Convert(source, conf, field)
				if !ok {
					return nil, false
				}
				fields = append(fields, fld)
			}

			return &Struct{
				baseField: base,
				source:    f,
				table: &NodeTable{
					Name:   f.Struct,
					Fields: fields,
				}, // TODO: struct not a NodeTable?!
			}, true
		}

	case *parser.FieldNode:
		{
			var node *parser.Node
			for _, elem := range source.Nodes {
				if elem.Name == f.Node {
					node = elem
					break
				}
			}

			return &Node{
				baseField: base,
				source:    f,
				table: &NodeTable{
					Name:       f.Node,
					Fields:     nil, // TODO: needed? -> node.Fields
					Timestamps: node.Timestamps,
				},
			}, true
		}

	case *parser.FieldEdge:
		{
			var edge *parser.Edge
			for _, elem := range source.Edges {
				if elem.Name == f.Edge {
					edge = elem
					break
				}
			}

			in, ok := Convert(source, conf, edge.In)
			if !ok {
				return nil, false
			}

			out, ok := Convert(source, conf, edge.Out)
			if !ok {
				return nil, false
			}

			var fields []Field
			for _, field := range edge.Fields {
				fld, ok := Convert(source, conf, field)
				if !ok {
					return nil, false
				}
				fields = append(fields, fld)
			}

			return &Edge{
				baseField: base,
				source:    f,
				table: &EdgeTable{
					Name:       f.Edge,
					In:         in.(*Node),
					Out:        out.(*Node),
					Fields:     fields,
					Timestamps: edge.Timestamps,
				},
			}, true
		}

	case *parser.FieldSlice:
		{
			element, ok := Convert(source, conf, f.Field)
			if !ok {
				return nil, false
			}

			return &Slice{
				baseField: base,
				source:    f,
				element:   element,
			}, true
		}

	case *parser.FieldEnum:
		{
			var values []string
			for _, val := range source.EnumValues {
				if val.Enum == f.Typ {
					values = append(values, val.Value)
				}
			}

			return &Enum{
				baseField: base,
				source:    f,
				model:     EnumModel(f.Typ),
				values:    values,
			}, true
		}

	case *parser.FieldTime:
		{
			return &Time{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldUUID:
		{
			return &UUID{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldURL:
		{
			return &URL{
				baseField: base,
				source:    f,
			}, true
		}
	}

	return nil, false
}
