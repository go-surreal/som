package parser

import (
	"fmt"
	"github.com/go-surreal/som/core/util"
	"github.com/wzshiming/gotype"
	"go/ast"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const packagePath = "github.com/go-surreal/som"

func Parse(dir string) (*Output, error) {
	res := &Output{}

	imp := gotype.NewImporter()

	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(dir, "./") {
		dir = "./" + dir
	}

	n, err := imp.Import(dir, workDir)
	if err != nil {
		return nil, err
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("could not find absolute path: %v", err)
	}

	mod, err := util.FindGoMod(absDir)
	if err != nil {
		return nil, err
	}

	diff := strings.TrimPrefix(absDir, mod.Dir())
	res.PkgPath = path.Join(mod.Module(), diff)

	nc := n.NumChild()
	for i := 0; i < nc; i++ {
		v := n.Child(i)

		switch {

		case !ast.IsExported(v.Name()):
			{
				continue
			}

		case isNode(v):
			{
				node, err := parseNode(v)
				if err != nil {
					return nil, err
				}
				res.Nodes = append(res.Nodes, node)
				continue
			}

		case isEdge(v):
			{
				edge, err := parseEdge(v)
				if err != nil {
					return nil, err
				}
				res.Edges = append(res.Edges, edge)
				continue
			}

		case v.Kind() == gotype.Struct:
			{
				// TODO: prevent external structs!

				str, err := parseStruct(v)
				if err != nil {
					return nil, err
				}
				res.Structs = append(res.Structs, str)
				continue
			}

		case v.Kind() == gotype.String && v.PkgPath() == packagePath:
			{
				res.Enums = append(res.Enums, &Enum{
					Name: v.Name(),
				})
				continue
			}

		case v.Kind() == gotype.Declaration:
			{
				res.EnumValues = append(res.EnumValues, &EnumValue{
					Enum:     v.Declaration().Name(),
					Variable: v.Name(),
					Value:    strings.Trim(v.Value(), "\""),
				})
				continue
			}

		default:
			{
				fmt.Println("ignoring:", v)
			}
		}
	}

	return res, nil
}

func isNode(t gotype.Type) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)

		if f.Name() == "Node" && f.Elem().Name() == "Node" &&
			f.Elem().String() == "som.Node" && f.Elem().PkgPath() == packagePath {
			return true
		}
	}

	return false
}

func isEdge(t gotype.Type) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)

		if f.Name() == "Edge" && f.Elem().Name() == "Edge" &&
			f.Elem().String() == "som.Edge" && f.Elem().PkgPath() == packagePath {
			return true
		}
	}

	return false
}

func isEnum(t gotype.Type) bool {
	if t.Kind() != gotype.String {
		return false
	}

	return t.String() != "string" && t.PkgPath() == packagePath // TODO: might not be an enum..?!
}

func parseNode(v gotype.Type) (*Node, error) {
	node := &Node{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() != packagePath {
				return nil, fmt.Errorf("model %s: anonymous field %s not allowed", v.Name(), f.Name())
			}

			if f.Name() == "Node" {
				node.Fields = append(node.Fields,
					&FieldID{&fieldAtomic{"ID", false}},
				)
				continue
			}

			if f.Name() == "Timestamps" {
				node.Timestamps = true
				node.Fields = append(node.Fields,
					&FieldTime{
						&fieldAtomic{"CreatedAt", false},
						true,
						false,
					},
					&FieldTime{
						&fieldAtomic{"UpdatedAt", false},
						false,
						true,
					},
				)
				continue
			}

			return nil, fmt.Errorf("model %s: unexpected anonymous field %s", v.Name(), f.Name())
		}

		// prevent custom ID field
		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("model %s: field ID not allowed, already provided by som.Node", v.Name())
		}

		field, err := parseField(f)
		if err != nil {
			return nil, err
		}

		node.Fields = append(node.Fields, field)
	}

	return node, nil
}

func parseEdge(v gotype.Type) (*Edge, error) {
	edge := &Edge{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() != packagePath {
				return nil, fmt.Errorf("model %s: anonymous field %s not allowed", v.Name(), f.Name())
			}

			if f.Name() == "Edge" {
				edge.Fields = append(edge.Fields,
					&FieldID{&fieldAtomic{"ID", false}},
				)
				continue
			}

			if f.Name() == "Timestamps" {
				edge.Timestamps = true
				edge.Fields = append(edge.Fields,
					&FieldTime{
						&fieldAtomic{"CreatedAt", false},
						true,
						false,
					},
					&FieldTime{
						&fieldAtomic{"UpdatedAt", false},
						false,
						true,
					},
				)
				continue
			}

			return nil, fmt.Errorf("model %s: unexpected anonymous field %s", v.Name(), f.Name())
		}

		// prevent custom ID field
		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("model %s: field ID not allowed, already provided by som.Edge", v.Name())
		}

		field, err := parseField(f)
		if err != nil {
			return nil, err
		}

		if f.Tag().Get("som") == "in" {
			edge.In = field
			continue
		}

		if f.Tag().Get("som") == "out" {
			edge.Out = field
			continue
		}

		edge.Fields = append(edge.Fields, field)
	}

	return edge, nil
}

func parseStruct(v gotype.Type) (*Struct, error) {
	str := &Struct{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		field, err := parseField(f)
		if err != nil {
			return nil, err
		}

		str.Fields = append(str.Fields, field)
	}

	return str, nil
}

func parseField(t gotype.Type) (Field, error) {
	atomic := &fieldAtomic{
		name:    t.Name(),
		pointer: false,
	}

	switch t.Elem().Kind() {

	case gotype.String:
		{
			switch {
			case isEnum(t.Elem()):
				{
					return &FieldEnum{atomic, t.Elem().Name()}, nil
				}
			default:
				{
					return &FieldString{atomic}, nil
				}
			}
		}

	case gotype.Int:
		{
			return &FieldNumeric{atomic, NumberInt}, nil
		}

	case gotype.Int8:
		{
			return &FieldNumeric{atomic, NumberInt8}, nil
		}

	case gotype.Int16:
		{
			return &FieldNumeric{atomic, NumberInt16}, nil
		}

	case gotype.Int32:
		{
			return &FieldNumeric{atomic, NumberInt32}, nil
		}

	case gotype.Int64:
		{
			return &FieldNumeric{atomic, NumberInt64}, nil
		}

	//case gotype.Uint:
	//	{
	//		return &FieldNumeric{atomic, NumberUint}, nil
	//	}

	case gotype.Uint8:
		{
			return &FieldNumeric{atomic, NumberUint8}, nil
		}

	case gotype.Uint16:
		{
			return &FieldNumeric{atomic, NumberUint16}, nil
		}

	case gotype.Uint32:
		{
			return &FieldNumeric{atomic, NumberUint32}, nil
		}

	//case gotype.Uint64:
	//	{
	//		return &FieldNumeric{atomic, NumberUint64}, nil
	//	}

	//case gotype.Uintptr:
	//	{
	//		return &FieldNumeric{atomic, NumberUintptr}, nil
	//	}

	case gotype.Float32:
		{
			return &FieldNumeric{atomic, NumberFloat32}, nil
		}

	case gotype.Float64:
		{
			return &FieldNumeric{atomic, NumberFloat64}, nil
		}

	case gotype.Rune:
		{
			return &FieldNumeric{atomic, NumberRune}, nil
		}

	case gotype.Bool:
		{
			return &FieldBool{atomic}, nil
		}

	case gotype.Byte:
		{
			return &FieldByte{atomic}, nil
		}

	case gotype.Struct:
		{
			// TODO: prevent structs (or general types) from another package (except time and uuid)!
			switch {
			case t.Elem().PkgPath() == "time" && t.Elem().Name() == "Time":
				{
					return &FieldTime{atomic, false, false}, nil
				}
			case t.Elem().PkgPath() == "net/url" && t.Elem().Name() == "URL":
				{
					return &FieldURL{atomic}, nil
				}
			case isNode(t.Elem()):
				{
					return &FieldNode{atomic, t.Elem().Name()}, nil
				}
			case isEdge(t.Elem()):
				{
					return &FieldEdge{atomic, t.Elem().Name()}, nil
				}
			default:
				{
					return &FieldStruct{atomic, t.Elem().Name()}, nil
				}
			}
		}

	case gotype.Slice:
		{
			field, err := parseField(t.Elem())
			if err != nil {
				return nil, err
			}

			return &FieldSlice{
				&fieldAtomic{name: t.Name()},
				field,
			}, nil
		}

	case gotype.Array:
		{
			if t.Elem().PkgPath() == "github.com/google/uuid" {
				return &FieldUUID{&fieldAtomic{name: t.Name()}}, nil
			}
		}

	case gotype.Ptr:
		{
			field, err := parseField(t.Elem())
			if err != nil {
				return nil, fmt.Errorf("could not parse elem for ptr field %s: %v", t.Name(), err)
			}

			if t.Name() != "" {
				field.setName(t.Name())
			}
			field.setPointer(true)

			return field, nil
		}
	}

	return nil, fmt.Errorf("field %s has unsupported type %s", t.Name(), t.Elem().Kind())
}

type Output struct {
	PkgPath    string
	Nodes      []*Node
	Edges      []*Edge
	Structs    []*Struct
	Enums      []*Enum
	EnumValues []*EnumValue
}
