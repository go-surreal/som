package parser

import (
	"fmt"
	"go/ast"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-surreal/som/core/util/gomod"
	"github.com/iancoleman/strcase"
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
				node, err := parseNode(v, outPkg, n)
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

				if isComplexIDStruct(v, outPkg) {
					continue
				}

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

	// Parse analyzer and search definitions.
	define, err := parseDefine(absDir)
	if err != nil {
		return nil, fmt.Errorf("could not parse define files: %w", err)
	}
	res.Define = define

	// Collect which optional features are used.
	res.UsedFeatures = collectUsedFeatures(res)

	return res, nil
}

func isNode(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)

		if !f.IsAnonymous() {
			continue
		}

		if f.Name() == "Node" && f.Elem().Name() == "Node" &&
			f.Elem().PkgPath() == outPkg {
			return true
		}

		if f.Name() == "CustomNode" {
			if isCustomNodeFromSom(f.Elem(), outPkg) {
				return true
			}
		}

	}

	return false
}

func isGenericNodeFromSom(t gotype.Type, outPkg string, name string) bool {
	if pkgPath := t.PkgPath(); pkgPath != "" {
		return pkgPath == outPkg
	}

	origin := t.Origin()
	if origin == nil {
		return false
	}
	indexExpr, ok := origin.(*ast.IndexExpr)
	if !ok {
		return false
	}
	selExpr, ok := indexExpr.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if selExpr.Sel.Name != name {
		return false
	}
	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == path.Base(outPkg)
}

func isCustomNodeFromSom(t gotype.Type, outPkg string) bool {
	return isGenericNodeFromSom(t, outPkg, "CustomNode")
}

func isKnownStringIDType(t gotype.Type) bool {
	origin := t.Origin()
	if origin == nil {
		return true // Node alias (no type arg) is ULID
	}
	indexExpr, ok := origin.(*ast.IndexExpr)
	if !ok {
		return true
	}
	_, ok = indexExpr.Index.(*ast.SelectorExpr)
	return ok // SelectorExpr means som.ULID/som.UUID/som.Rand
}

func isComplexIDStruct(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}
	nf := t.NumField()
	for i := 0; i < nf; i++ {
		f := t.Field(i)
		if !f.IsAnonymous() {
			continue
		}
		if f.Name() == "ArrayID" || f.Name() == "ObjectID" {
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

		if !f.IsAnonymous() {
			continue
		}

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

func parseIDType(t gotype.Type) IDType {
	origin := t.Origin()
	if origin == nil {
		return IDTypeULID
	}

	if indexExpr, ok := origin.(*ast.IndexExpr); ok {
		if selExpr, ok := indexExpr.Index.(*ast.SelectorExpr); ok {
			switch selExpr.Sel.Name {
			case "UUID":
				return IDTypeUUID
			case "Rand":
				return IDTypeRand
			case "ULID":
				return IDTypeULID
			}
		}
	}

	return IDTypeULID
}

func parseComplexIDFields(t gotype.Type, outPkg string, pkgScope gotype.Type) (*FieldComplexID, error) {
	origin := t.Origin()
	if origin == nil {
		return nil, fmt.Errorf("complex ID type has no AST origin")
	}

	indexExpr, ok := origin.(*ast.IndexExpr)
	if !ok {
		return nil, fmt.Errorf("complex ID type is not a generic instantiation (expected IndexExpr, got %T)", origin)
	}

	var typeArgName string
	switch idx := indexExpr.Index.(type) {
	case *ast.Ident:
		typeArgName = idx.Name
	case *ast.SelectorExpr:
		typeArgName = idx.Sel.Name
	default:
		return nil, fmt.Errorf("complex ID type argument: unsupported AST node %T", indexExpr.Index)
	}

	keyType, ok := pkgScope.ChildByName(typeArgName)
	if !ok {
		return nil, fmt.Errorf("complex ID type argument %s not found in package scope", typeArgName)
	}

	if keyType.Kind() != gotype.Struct {
		return nil, fmt.Errorf("complex ID type parameter must be a struct, got %s", keyType.Kind())
	}

	structName := keyType.Name()
	nf := keyType.NumField()
	fields := make([]ComplexIDField, 0, nf)

	var kind IDType
	for i := 0; i < nf; i++ {
		sf := keyType.Field(i)

		if sf.IsAnonymous() {
			switch sf.Name() {
			case "ArrayID":
				kind = IDTypeArray
			case "ObjectID":
				kind = IDTypeObject
			}
			continue
		}

		if !ast.IsExported(sf.Name()) {
			continue
		}

		parsed, err := parseFieldInternal(sf, outPkg, false)
		if err != nil {
			return nil, fmt.Errorf("complex ID field %s: %w", sf.Name(), err)
		}

		switch parsed.(type) {
		case *FieldString, *FieldNumeric, *FieldBool, *FieldTime, *FieldDuration, *FieldUUID, *FieldNode:
		default:
			return nil, fmt.Errorf("complex ID field %s: unsupported type %T (only string, numeric, bool, time.Time, time.Duration, UUID, and node references are allowed)", sf.Name(), parsed)
		}

		fields = append(fields, ComplexIDField{
			Name:   sf.Name(),
			DBName: strcase.ToSnake(sf.Name()),
			Field:  parsed,
		})
	}

	if kind == "" {
		return nil, fmt.Errorf("complex ID struct %s must embed som.ArrayID or som.ObjectID", structName)
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("complex ID struct %s has no exported fields", structName)
	}

	return &FieldComplexID{
		fieldAtomic: &fieldAtomic{name: "ID"},
		Kind:        kind,
		StructName:  structName,
		Fields:      fields,
	}, nil
}

func parseNode(v gotype.Type, outPkg string, pkgScope gotype.Type) (*Node, error) {
	internalPkg := path.Join(outPkg, "internal")

	node := &Node{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if (f.Name() == "Node" && f.Elem().PkgPath() == outPkg) ||
				(f.Name() == "CustomNode" && isCustomNodeFromSom(f.Elem(), outPkg)) {
				if isKnownStringIDType(f.Elem()) {
					gen := parseIDType(f.Elem())
					node.IDType = gen
					node.IDEmbed = f.Name()
					node.Fields = append(node.Fields,
						&FieldID{&fieldAtomic{name: "ID"}, gen},
					)
				} else {
					node.IDEmbed = f.Name()

					complexID, err := parseComplexIDFields(f.Elem(), outPkg, pkgScope)
					if err != nil {
						return nil, fmt.Errorf("model %s: %w", v.Name(), err)
					}
					node.IDType = complexID.Kind
					node.ComplexID = complexID
					node.Fields = append(node.Fields, complexID)
				}
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

			if f.Elem().PkgPath() == internalPkg && f.Name() == "OptimisticLock" {
				node.OptimisticLock = true
				continue
			}

			if f.Elem().PkgPath() == internalPkg && f.Name() == "SoftDelete" {
				node.SoftDelete = true
				node.Fields = append(node.Fields,
					&FieldTime{
						fieldAtomic: &fieldAtomic{name: "DeletedAt", pointer: true},
						IsDeletedAt: true,
					},
				)
				continue
			}

			return nil, fmt.Errorf("model %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		// prevent custom ID field
		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("model %s: field ID not allowed, already provided by som.%s", v.Name(), node.IDEmbed)
		}

		field, err := parseField(f, outPkg)
		if err != nil {
			return nil, err
		}

		node.Fields = append(node.Fields, field)
	}

	// If OptimisticLock is enabled, always add a version field
	if node.OptimisticLock {
		node.Fields = append(node.Fields, &FieldVersion{&fieldAtomic{name: "Version", pointer: false}})
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
					&FieldID{&fieldAtomic{name: "ID"}, IDTypeULID},
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

			if f.Elem().PkgPath() == internalPkg && f.Name() == "OptimisticLock" {
				edge.OptimisticLock = true
				continue
			}

			if f.Elem().PkgPath() == internalPkg && f.Name() == "SoftDelete" {
				edge.SoftDelete = true
				edge.Fields = append(edge.Fields,
					&FieldTime{
						fieldAtomic: &fieldAtomic{name: "DeletedAt", pointer: true},
						IsDeletedAt: true,
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

	// If OptimisticLock is enabled, always add a version field
	if edge.OptimisticLock {
		edge.Fields = append(edge.Fields, &FieldVersion{&fieldAtomic{name: "Version", pointer: false}})
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
	field, err := parseFieldInternal(t, outPkg, true)
	if err != nil {
		return nil, err
	}

	if err := field.Validate(); err != nil {
		return nil, err
	}

	return field, nil
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
					return &FieldTime{atomic, false, false, false}, nil
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
			pkgPath := t.Elem().PkgPath()
			switch pkgPath {
			case string(UUIDPackageGoogle):
				return &FieldUUID{
					fieldAtomic: &fieldAtomic{
						name:   t.Name(),
						index:  indexInfo,
						search: searchInfo,
					},
					Package: UUIDPackageGoogle,
				}, nil
			case string(UUIDPackageGofrs):
				return &FieldUUID{
					fieldAtomic: &fieldAtomic{
						name:   t.Name(),
						index:  indexInfo,
						search: searchInfo,
					},
					Package: UUIDPackageGofrs,
				}, nil
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
//   - som:"index,name:idx_foo" -> named index (for composite)
//   - som:"unique"             -> unique index on single field
//   - som:"unique(name)"       -> composite unique index (grouped by name)
//   - som:"fulltext(config_name)" -> fulltext search with named config
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

		// Check for fulltext(config_name)
		if strings.HasPrefix(part, "fulltext(") && strings.HasSuffix(part, ")") {
			configName := part[9 : len(part)-1] // len("fulltext(") = 9
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

		// Check for unique or unique(name)
		if part == "unique" || strings.HasPrefix(part, "unique(") {
			if index == nil {
				index = &IndexInfo{}
			}
			index.Unique = true

			// Parse unique(name) for composite unique index
			if strings.HasPrefix(part, "unique(") && strings.HasSuffix(part, ")") {
				// Extract the name from unique(name)
				inner := part[7 : len(part)-1] // len("unique(") = 7
				index.UniqueName = inner
			}
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

	// Define holds analyzer and search definitions.
	Define *DefineOutput

	// UsedFeatures tracks which optional features are used in the models.
	UsedFeatures *UsedFeatures
}

// UsedFeatures tracks which optional features are used across all models.
type UsedFeatures struct {
	UsesGoogleUUID bool
	UsesGofrsUUID  bool
}

// collectUsedFeatures walks through all parsed fields to determine which features are used.
func collectUsedFeatures(output *Output) *UsedFeatures {
	features := &UsedFeatures{}

	var checkField func(f Field)
	checkField = func(f Field) {
		switch field := f.(type) {
		case *FieldUUID:
			switch field.Package {
			case UUIDPackageGoogle:
				features.UsesGoogleUUID = true
			case UUIDPackageGofrs:
				features.UsesGofrsUUID = true
			}
		case *FieldSlice:
			checkField(field.Field)
		}
	}

	for _, node := range output.Nodes {
		for _, field := range node.Fields {
			checkField(field)
		}
	}

	for _, edge := range output.Edges {
		for _, field := range edge.Fields {
			checkField(field)
		}
	}

	for _, str := range output.Structs {
		for _, field := range str.Fields {
			checkField(field)
		}
	}

	return features
}
