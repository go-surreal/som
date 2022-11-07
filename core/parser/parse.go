package parser

import (
	"fmt"
	"go/ast"
	"path"
	"path/filepath"
	"strings"

	"github.com/wzshiming/gotype"
)

const fileGoMod = "go.mod"

const packagePath = "github.com/marcbinz/sdb"

func Parse(dir string) (*Output, error) {
	res := &Output{}

	imp := gotype.NewImporter()

	n, err := imp.Import(dir, "")
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
			f.Elem().String() == "sdb.Node" && f.Elem().PkgPath() == packagePath {
			return true
		}
	}

	return false
}

func isEnum(t gotype.Type) bool {
	if t.Kind() != gotype.String {
		return false
	}

	return t.String() != "string" && t.PkgPath() == "github.com/marcbinz/sdb" // TODO: might not be an enum..?!
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

	switch t.Elem().Kind() {
	case gotype.String:
		if t.Name() == "ID" {
			field = &FieldID{&fieldAtomic{Name: t.Name()}}
		} else if isEnum(t.Elem()) {
			field = &FieldEnum{&fieldAtomic{Name: t.Name()}, t.Elem().String()}
		} else {
			field = &FieldString{&fieldAtomic{Name: t.Name()}}
		}
	case gotype.Int:
		field = &FieldNumeric{&fieldAtomic{Name: t.Name()}, NumberInt}
	case gotype.Int32:
		field = &FieldNumeric{&fieldAtomic{Name: t.Name()}, NumberInt32}
	case gotype.Int64:
		field = &FieldNumeric{&fieldAtomic{Name: t.Name()}, NumberInt64}
	case gotype.Float32:
		field = &FieldNumeric{&fieldAtomic{Name: t.Name()}, NumberFloat32}
	case gotype.Float64:
		field = &FieldNumeric{&fieldAtomic{Name: t.Name()}, NumberFloat64}
	case gotype.Bool:
		field = &FieldBool{&fieldAtomic{Name: t.Name()}}
	case gotype.Struct:
		// TODO: prevent structs (or general types) from another package (except time and uuid)!
		if t.Elem().PkgPath() == "time" {
			field = &FieldTime{&fieldAtomic{Name: t.Name()}}
		} else if isNode(t.Elem()) {
			field = &FieldNode{&fieldAtomic{Name: t.Name()}, t.Elem().Name(), false} // TODO: handle pointers
		} else {
			field = &FieldStruct{&fieldAtomic{Name: t.Name()}, t.Elem().Name(), false} // TODO: handle pointers
		}
	case gotype.Slice:
		field = &FieldSlice{&fieldAtomic{Name: t.Name()}, t.Elem().Elem().Name(), isNode(t.Elem().Elem()), isEnum(t.Elem().Elem())}
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
	Structs    []*Struct
	Enums      []*Enum
	EnumValues []*EnumValue
}
