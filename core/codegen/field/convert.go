package field

import (
	"fmt"

	"github.com/go-surreal/som/core/parser"
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
			Name:   node.Name,
			Source: node,
		}

		for _, f := range node.Fields {
			dbField, ok := Convert(source, buildConf, f)
			if !ok {
				return nil, fmt.Errorf("could not convert field a: %v", f)
			}
			dbNode.Fields = append(dbNode.Fields, dbField)
		}

		def.Nodes = append(def.Nodes, dbNode)
	}

	for _, edge := range source.Edges {
		dbEdge := &EdgeTable{
			Name:   edge.Name,
			Source: edge,
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
				return nil, fmt.Errorf("could not convert field b: %v", f)
			}
			dbEdge.Fields = append(dbEdge.Fields, dbField)
		}

		def.Edges = append(def.Edges, dbEdge)
	}

	for _, dbNode := range def.Nodes {
		for _, f := range dbNode.Fields {
			if cid, ok := f.(*ComplexID); ok && cid.element != nil {
				def.Objects = append(def.Objects, &DatabaseObject{
					Name:           cid.element.NameGo(),
					Fields:         cid.element.GetFields(),
					IsArrayIndexed: cid.source.Kind == parser.IDTypeArray,
				})
			}
		}
	}

	for _, str := range source.Structs {
		dbObject := &DatabaseObject{
			Name: str.Name,
		}

		for _, f := range str.Fields {
			dbField, ok := Convert(source, buildConf, f)
			if !ok {
				return nil, fmt.Errorf("could not convert field c: %v", f)
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

	case *parser.FieldByte:
		{
			return &Byte{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldTime:
		{
			return &Time{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldDuration:
		{
			return &Duration{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldMonth:
		{
			return &Month{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldWeekday:
		{
			return &Weekday{
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

	case *parser.FieldRegex:
		{
			return &Regex{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldPassword:
		{
			return &Password{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldEmail:
		{
			return &Email{
				baseField: base,
				source:    f,
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
				return nil, false // TODO: anonymous struct type not supported // return error msg!
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
				element: &NodeTable{
					Name:   f.Struct,
					Fields: fields,
				}, // TODO: struct not a NodeTable?!
			}, true
		}

	case *parser.FieldNode:
		{
			// Find the source node to get its properties (like SoftDelete)
			var sourceNode *parser.Node
			for _, node := range source.Nodes {
				if node.Name == f.Node {
					sourceNode = node
					break
				}
			}

			return &Node{
				baseField: base,
				source:    f,
				table: &NodeTable{
					Name:   f.Node,
					Fields: nil, // TODO: needed? -> node.Fields
					Source: sourceNode,
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
					Name:   f.Edge,
					In:     in.(*Node),
					Out:    out.(*Node),
					Fields: fields,
					Source: edge,
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

	case *parser.FieldVersion:
		{
			return &Version{
				baseField: base,
				source:    f,
			}, true
		}

	case *parser.FieldComplexID:
		{
			var fields []Field
			for _, sf := range f.Fields {
				if _, ok := sf.Field.(*parser.FieldNode); ok {
					continue
				}
				fld, ok := Convert(source, conf, sf.Field)
				if !ok {
					continue
				}
				fields = append(fields, fld)
			}

			cid := &ComplexID{baseField: base, source: f}
			if len(fields) > 0 {
				cid.element = &NodeTable{Name: f.StructName, Fields: fields}
			}
			return cid, true
		}
	}

	return nil, false
}
