package parser

import (
	"fmt"
	"go/ast"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-surreal/som/core/util/gomod"
	"github.com/wzshiming/gotype"
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

	nc := n.NumChild()
	for i := 0; i < nc; i++ {
		v := n.Child(i)

		switch {

		case !ast.IsExported(v.Name()):
			{
				// TODO: imagine an unexported enum that is set with exported method.. possible, no?
				continue
			}

		case isNode(v, outPkg):
			{
				node, err := parseNode(v, outPkg)
				if err != nil {
					return nil, err
				}
				res.Nodes = append(res.Nodes, node)
				continue
			}

		case isEdge(v, outPkg):
			{
				edge, err := parseEdge(v, outPkg)
				if err != nil {
					return nil, err
				}
				res.Edges = append(res.Edges, edge)
				continue
			}

		case v.Kind() == gotype.Struct:
			{
				// TODO: prevent external structs!

				str, err := parseStruct(v, outPkg)
				if err != nil {
					return nil, err
				}
				res.Structs = append(res.Structs, str)
				continue
			}

		case isEnum(v, outPkg):
			{
				res.Enums = append(res.Enums, &Enum{
					Name: v.Name(),
				})
				continue
			}

		case v.Kind() == gotype.Declaration: // TODO: what about other decls? :/ -> new parser
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

	// Parse //go:build som config files for analyzer and search definitions
	config, err := ParseConfig(absDir)
	if err != nil {
		return nil, fmt.Errorf("could not parse config files: %w", err)
	}
	res.Config = config

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

func isEnum(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.String {
		return false
	}

	return t.String() != "string" && t.PkgPath() == outPkg // TODO: might not be an enum..?!
}

func isPassword(t gotype.Type, outPkg string) bool {
	if t.PkgPath() != outPkg {
		return false
	}
	return t.Name() == "Password"
}

func parsePasswordAlgorithm(t gotype.Type) PasswordAlgorithm {
	// Extract the type parameter from the AST
	origin := t.Origin()
	if origin == nil {
		return PasswordBcrypt
	}

	// Handle *ast.IndexExpr for Password[Algo]
	if indexExpr, ok := origin.(*ast.IndexExpr); ok {
		if selExpr, ok := indexExpr.Index.(*ast.SelectorExpr); ok {
			switch selExpr.Sel.Name {
			case "Bcrypt":
				return PasswordBcrypt
			case "Argon2":
				return PasswordArgon2
			case "Pbkdf2":
				return PasswordPbkdf2
			case "Scrypt":
				return PasswordScrypt
			}
		}
	}

	// Default to Bcrypt
	return PasswordBcrypt
}

func parseNode(v gotype.Type, outPkg string) (*Node, error) {
	internalPkg := path.Join(outPkg, "internal")

	node := &Node{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() == outPkg && f.Name() == "Node" {
				node.Fields = append(node.Fields,
					&FieldID{&fieldAtomic{name: "ID"}},
				)
				continue
			}

			if f.Elem().PkgPath() == internalPkg && f.Name() == "Timestamps" {
				node.Timestamps = true
				node.Fields = append(node.Fields,
					&FieldTime{
						fieldAtomic: &fieldAtomic{name: "CreatedAt"},
						IsCreatedAt: true,
						IsUpdatedAt: false,
					},
					&FieldTime{
						fieldAtomic: &fieldAtomic{name: "UpdatedAt"},
						IsCreatedAt: false,
						IsUpdatedAt: true,
					},
				)
				continue
			}

			return nil, fmt.Errorf("model %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		// prevent custom ID field
		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("model %s: field ID not allowed, already provided by som.Node", v.Name())
		}

		field, err := parseField(f, outPkg)
		if err != nil {
			return nil, err
		}

		node.Fields = append(node.Fields, field)
	}

	return node, nil
}

func parseEdge(v gotype.Type, outPkg string) (*Edge, error) {
	internalPkg := path.Join(outPkg, "internal")

	edge := &Edge{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() == outPkg && f.Name() == "Edge" {
				edge.Fields = append(edge.Fields,
					&FieldID{&fieldAtomic{name: "ID"}},
				)
				continue
			}

			if f.Elem().PkgPath() == internalPkg && f.Name() == "Timestamps" {
				edge.Timestamps = true
				edge.Fields = append(edge.Fields,
					&FieldTime{
						fieldAtomic: &fieldAtomic{name: "CreatedAt"},
						IsCreatedAt: true,
						IsUpdatedAt: false,
					},
					&FieldTime{
						fieldAtomic: &fieldAtomic{name: "UpdatedAt"},
						IsCreatedAt: false,
						IsUpdatedAt: true,
					},
				)
				continue
			}

			return nil, fmt.Errorf("model %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		// prevent custom ID field
		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("model %s: field ID not allowed, already provided by som.Edge", v.Name())
		}

		field, err := parseField(f, outPkg)
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

func parseStruct(v gotype.Type, outPkg string) (*Struct, error) {
	str := &Struct{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		field, err := parseField(f, outPkg)
		if err != nil {
			return nil, err
		}

		str.Fields = append(str.Fields, field)
	}

	return str, nil
}

func parseField(t gotype.Type, outPkg string) (Field, error) {
	return parseFieldInternal(t, outPkg, true)
}

func parseFieldInternal(t gotype.Type, outPkg string, isStructField bool) (Field, error) {
	// Parse som tag for index/search info
	// Only struct fields have tags
	var indexInfo *IndexInfo
	var searchInfo *SearchInfo
	if isStructField {
		somTag := t.Tag().Get("som")
		indexInfo, searchInfo = parseSomTag(somTag)
	}

	atomic := &fieldAtomic{
		name:    t.Name(),
		pointer: false,
		index:   indexInfo,
		search:  searchInfo,
	}

	switch t.Elem().Kind() {

	case gotype.String:
		{
			switch {
			case isPassword(t.Elem(), outPkg):
				{
					return &FieldPassword{atomic, parsePasswordAlgorithm(t.Elem())}, nil
				}
			case t.Elem().PkgPath() == outPkg && t.Elem().Name() == "Email":
				{
					return &FieldEmail{atomic}, nil
				}
			case isEnum(t.Elem(), outPkg):
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
			field, err := parseFieldInternal(t.Elem(), outPkg, false)
			if err != nil {
				return nil, err
			}

			return &FieldSlice{
				&fieldAtomic{
					name:   t.Name(),
					index:  indexInfo,
					search: searchInfo,
				},
				field,
			}, nil
		}

	case gotype.Array:
		{
			if t.Elem().PkgPath() == "github.com/google/uuid" {
				return &FieldUUID{&fieldAtomic{
					name:   t.Name(),
					index:  indexInfo,
					search: searchInfo,
				}}, nil
			}
		}

	case gotype.Ptr:
		{
			field, err := parseFieldInternal(t.Elem(), outPkg, false)
			if err != nil {
				return nil, fmt.Errorf("could not parse elem for ptr field %s: %v", t.Name(), err)
			}

			if t.Name() != "" {
				field.setName(t.Name())
			}
			field.setPointer(true)

			// Index/search info comes from the outer pointer field, not the inner element
			if indexInfo != nil {
				field.setIndex(indexInfo)
			}
			if searchInfo != nil {
				field.setSearch(searchInfo)
			}

			return field, nil
		}
	}

	return nil, fmt.Errorf("field %s has unsupported type %s", t.Name(), t.Elem().Kind())
}

// parseSomTag parses the "som" struct tag and extracts index/search info.
// Tag format examples:
//   - som:"index"              -> simple index
//   - som:"index,unique"       -> unique index
//   - som:"index,name:idx_foo" -> named index (for composite)
//   - som:"search:config_name" -> fulltext search with named config
//   - som:"in"                 -> edge in field (existing)
//   - som:"out"                -> edge out field (existing)
func parseSomTag(tag string) (index *IndexInfo, search *SearchInfo) {
	if tag == "" {
		return nil, nil
	}

	// Handle existing edge markers
	if tag == "in" || tag == "out" {
		return nil, nil
	}

	parts := strings.Split(tag, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Check for search:config_name
		if strings.HasPrefix(part, "search:") {
			configName := strings.TrimPrefix(part, "search:")
			search = &SearchInfo{ConfigName: configName}
			continue
		}

		// Check for index keyword
		if part == "index" {
			if index == nil {
				index = &IndexInfo{}
			}
			continue
		}

		// Check for unique modifier
		if part == "unique" {
			if index == nil {
				index = &IndexInfo{}
			}
			index.Unique = true
			continue
		}

		// Check for name:xxx modifier
		if strings.HasPrefix(part, "name:") {
			name := strings.TrimPrefix(part, "name:")
			if index == nil {
				index = &IndexInfo{}
			}
			index.Name = name
			continue
		}
	}

	return index, search
}

type Output struct {
	PkgPath    string
	Nodes      []*Node
	Edges      []*Edge
	Structs    []*Struct
	Enums      []*Enum
	EnumValues []*EnumValue

	// Config holds analyzer and search definitions from //go:build som files
	Config *ConfigOutput
}
