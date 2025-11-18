package parser

import (
	"fmt"
	"github.com/go-surreal/som/core/util/gomod"
	"github.com/wzshiming/gotype"
	"go/ast"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Parse(dir string, outPkg string) (*Output, error) {
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

	mod, err := gomod.FindGoMod(absDir)
	if err != nil {
		return nil, err
	}

	diff := strings.TrimPrefix(absDir, mod.Dir())
	res.PkgPath = path.Join(mod.Module(), diff)

	// First pass: collect enum names
	enumNames := make(map[string]bool)
	nc := n.NumChild()
	for i := 0; i < nc; i++ {
		v := n.Child(i)

		if !ast.IsExported(v.Name()) {
			continue
		}

		if isEnum(v, res.PkgPath) {
			enumNames[v.Name()] = true
			enumType := v.Kind().String()
			res.Enums = append(res.Enums, &Enum{
				Name: v.Name(),
				Type: enumType,
			})
		}
	}

	// Second pass: parse types using enum names
	for i := 0; i < nc; i++ {
		v := n.Child(i)

		switch {

		case !ast.IsExported(v.Name()):
			{
				continue
			}

		case isNode(v, outPkg):
			{
				node, err := parseNode(v, outPkg, res.PkgPath, enumNames)
				if err != nil {
					return nil, err
				}
				res.Nodes = append(res.Nodes, node)
				continue
			}

		case isEdge(v, outPkg):
			{
				edge, err := parseEdge(v, outPkg, res.PkgPath, enumNames)
				if err != nil {
					return nil, err
				}
				res.Edges = append(res.Edges, edge)
				continue
			}

		case v.Kind() == gotype.Struct:
			{
				// TODO: prevent external structs!

				str, err := parseStruct(v, outPkg, res.PkgPath, enumNames)
				if err != nil {
					return nil, err
				}
				res.Structs = append(res.Structs, str)
				continue
			}

		case isEnum(v, res.PkgPath):
			{
				// Already processed in first pass
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

func isNode(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)

		if f.Name() == "Node" && f.Elem().Name() == "Node" &&
			f.Elem().PkgPath() == outPkg {
			return true
		}
	}

	return false
}

func isEdge(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)

		if f.Name() == "Edge" && f.Elem().Name() == "Edge" &&
			f.Elem().PkgPath() == outPkg {
			return true
		}
	}

	return false
}

func isEnum(t gotype.Type, srcPkg string) bool {
	// Check if it's a named type
	if t.Name() == "" {
		return false
	}

	// Type aliases like `type Role som.EnumString` resolve to the underlying type.
	// So when gotype resolves the type, t.Name() will be "EnumString" for the alias "Role".
	// We need to check the origin to see the actual type definition.

	origin := t.Origin()
	if origin == nil {
		return false
	}

	// Handle *ast.Ident (when type is resolved through alias or field access)
	if ident, ok := origin.(*ast.Ident); ok {
		// Check if the identifier itself is an Enum* type
		return strings.HasPrefix(ident.Name, "Enum")
	}

	// Handle *ast.TypeSpec (when accessed directly from package)
	typeSpec, ok := origin.(*ast.TypeSpec)
	if !ok {
		return false
	}

	// Check if the type definition references an Enum* identifier
	switch typeExpr := typeSpec.Type.(type) {
	case *ast.Ident:
		// Direct reference like: type Role EnumString
		return strings.HasPrefix(typeExpr.Name, "Enum")
	case *ast.SelectorExpr:
		// Qualified reference like: type Role som.EnumString
		return strings.HasPrefix(typeExpr.Sel.Name, "Enum")
	}

	return false
}

func parseNode(v gotype.Type, outPkg, srcPkg string, enumNames map[string]bool) (*Node, error) {
	node := &Node{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() != outPkg {
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

		field, err := parseField(f, outPkg, srcPkg, enumNames)
		if err != nil {
			return nil, err
		}

		node.Fields = append(node.Fields, field)
	}

	return node, nil
}

func parseEdge(v gotype.Type, outPkg, srcPkg string, enumNames map[string]bool) (*Edge, error) {
	edge := &Edge{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() != outPkg {
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

		field, err := parseField(f, outPkg, srcPkg, enumNames)
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

func parseStruct(v gotype.Type, outPkg, srcPkg string, enumNames map[string]bool) (*Struct, error) {
	str := &Struct{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		field, err := parseField(f, outPkg, srcPkg, enumNames)
		if err != nil {
			return nil, err
		}

		str.Fields = append(str.Fields, field)
	}

	return str, nil
}

func parseField(t gotype.Type, outPkg, srcPkg string, enumNames map[string]bool) (Field, error) {
	atomic := &fieldAtomic{
		name:    t.Name(),
		pointer: false,
	}

	// Check for enum types first (supports any comparable base type)
	// Use the enumNames map to check if this field's type is a known enum
	if enumNames[t.Elem().Name()] {
		baseType := t.Elem().Kind().String()
		return &FieldEnum{atomic, t.Elem().Name(), baseType}, nil
	}

	switch t.Elem().Kind() {

	case gotype.String:
		{
			return &FieldString{atomic}, nil
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
			switch {
			case t.Elem().PkgPath() == "time" && t.Elem().Name() == "Duration":
				{
					return &FieldDuration{atomic}, nil
				}
			default:
				{
					return &FieldNumeric{atomic, NumberInt64}, nil
				}
			}

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
			case isNode(t.Elem(), outPkg):
				{
					return &FieldNode{atomic, t.Elem().Name()}, nil
				}
			case isEdge(t.Elem(), outPkg):
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
			field, err := parseField(t.Elem(), outPkg, srcPkg, enumNames)
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
			field, err := parseField(t.Elem(), outPkg, srcPkg, enumNames)
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
