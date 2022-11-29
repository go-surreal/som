package parser

import (
	"fmt"
	"github.com/wzshiming/gotype"
	"go/ast"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const fileGoMod = "go.mod"

const packagePath = "github.com/marcbinz/som"

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

	pkgPath, modPath, err := parseMod(absDir)
	if err != nil {
		return nil, err
	}

	diff := strings.TrimPrefix(absDir, modPath)
	res.PkgPath = path.Join(pkgPath, diff)

	nc := n.NumChild()
	for i := 0; i < nc; i++ {
		v := n.Child(i)

		if !ast.IsExported(v.Name()) {
			continue
		}

		if isNode(v) {
			node, err := parseNode(v)
			if err != nil {
				return nil, err
			}
			res.Nodes = append(res.Nodes, node)
			continue
		}

		if isEdge(v) {
			edge, err := parseEdge(v)
			if err != nil {
				return nil, err
			}
			res.Edges = append(res.Edges, edge)
			continue
		}

		if v.Kind() == gotype.Struct {
			// TODO: prevent external structs!

			str, err := parseStruct(v)
			if err != nil {
				return nil, err
			}
			res.Structs = append(res.Structs, str)
			continue
		}

		if v.Kind() == gotype.String && v.PkgPath() == packagePath {
			res.Enums = append(res.Enums, &Enum{
				Name: v.Name(),
			})
			continue
		}

		if v.Kind() == gotype.Declaration {
			res.EnumValues = append(res.EnumValues, &EnumValue{
				Enum:     v.Declaration().Name(),
				Variable: v.Name(),
				Value:    v.Value(),
			})
			continue
		}

		fmt.Println("ignoring:", v)
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

	return t.String() != "string" && t.PkgPath() == "github.com/marcbinz/som" // TODO: might not be an enum..?!
}

func parseNode(v gotype.Type) (*Node, error) {
	node := &Node{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if i == 0 {
			// TODO: better detect the "Node" field...
			continue
		}

		// TODO: ignore unexported fields?!

		// prevent ID from not being a string type
		if f.Name() == "ID" && f.Elem().Kind() != gotype.String {
			return nil, fmt.Errorf("field ID of model %s must be a string", v.Name())
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

		if i == 0 {
			// TODO: better detect the "Edge" field...
			continue
		}

		// TODO: ignore unexported fields?!

		// prevent ID from not being a string type
		if f.Name() == "ID" && f.Elem().Kind() != gotype.String {
			return nil, fmt.Errorf("field ID of model %s must be a string", v.Name())
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
	var field Field

	atomic := &fieldAtomic{Name: t.Name()}

	switch t.Elem().Kind() {
	case gotype.String:
		if t.Name() == "ID" {
			field = &FieldID{atomic}
		} else if isEnum(t.Elem()) {
			field = &FieldEnum{atomic, t.Elem().Name()}
		} else {
			field = &FieldString{atomic}
		}
	case gotype.Int:
		field = &FieldNumeric{atomic, NumberInt}
	case gotype.Int32:
		field = &FieldNumeric{atomic, NumberInt32}
	case gotype.Int64:
		field = &FieldNumeric{atomic, NumberInt64}
	case gotype.Float32:
		field = &FieldNumeric{atomic, NumberFloat32}
	case gotype.Float64:
		field = &FieldNumeric{atomic, NumberFloat64}
	case gotype.Bool:
		field = &FieldBool{atomic}
	case gotype.Struct:
		// TODO: prevent structs (or general types) from another package (except time and uuid)!
		if t.Elem().PkgPath() == "time" {
			field = &FieldTime{atomic}
		} else if isNode(t.Elem()) {
			field = &FieldNode{atomic, t.Elem().Name(), false} // TODO: handle pointers
		} else if isEdge(t.Elem()) {
			field = &FieldEdge{atomic, t.Elem().Name(), false} // TODO: handle pointers
		} else {
			field = &FieldStruct{atomic, t.Elem().Name(), false} // TODO: handle pointers
		}
	case gotype.Slice:
		subField, err := parseField(t.Elem())
		if err != nil {
			return nil, err
		}
		field = &FieldSlice{
			&fieldAtomic{Name: t.Name()},
			t.Elem().Elem().Name(),
			subField,
			isNode(t.Elem().Elem()),
			isEdge(t.Elem().Elem()),
			isEnum(t.Elem().Elem()),
		}
	// case gotype.Map:
	// 	field = FieldMap{fieldAtomic{Name: t.Name()}, t.Elem().Key().Name(), t.Elem().Elem().Name()}
	case gotype.Array:
		if t.Elem().PkgPath() == "github.com/google/uuid" {
			field = &FieldUUID{&fieldAtomic{Name: t.Name()}}
		}
	default:
		return nil, fmt.Errorf("field %s has unsupported type %s", t.Name(), t.Elem().Kind())
	}

	return field, nil
}

type Output struct {
	PkgPath    string
	Nodes      []*Node
	Edges      []*Edge
	Structs    []*Struct
	Enums      []*Enum
	EnumValues []*EnumValue
}
