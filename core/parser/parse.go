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

// activeFieldRegistry is set at the start of Parse() so that
// ParseField / ParseFieldInternal can delegate to it.
var activeFieldRegistry *fieldRegistry

func Parse(dir string, outPkg string, typeHandlers []TypeHandler, fieldHandlers []FieldHandler) (*Output, error) {
	res := &Output{}

	tReg := newTypeRegistry(typeHandlers)
	fReg := newFieldRegistry(fieldHandlers)

	activeFieldRegistry = fReg

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

	ctx := &TypeContext{OutPkg: outPkg, PkgScope: n, Output: res}

	nc := n.NumChild()
	for i := 0; i < nc; i++ {
		v := n.Child(i)

		if !ast.IsExported(v.Name()) {
			continue
		}

		matched, err := tReg.handle(v, ctx)
		if err != nil {
			return nil, err
		}
		if !matched {
			fmt.Println("ignoring:", v)
		}
	}

	define, err := parseDefine(absDir)
	if err != nil {
		return nil, fmt.Errorf("could not parse define files: %w", err)
	}
	res.Define = define

	if err := tReg.validate(ctx); err != nil {
		return nil, err
	}

	res.UsedFeatures = collectUsedFeatures(res)

	return res, nil
}

func ParseField(t gotype.Type, outPkg string) (Field, error) {
	field, err := ParseFieldInternal(t, outPkg, true)
	if err != nil {
		return nil, err
	}

	if err := field.Validate(); err != nil {
		return nil, err
	}

	return field, nil
}

func ParseFieldInternal(t gotype.Type, outPkg string, isStructField bool) (Field, error) {
	var tagInfo *TagInfo
	if isStructField {
		somTag := t.Tag().Get("som")
		var err error
		tagInfo, err = parseSomTag(somTag)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", t.Name(), err)
		}
	}

	ctx := &FieldContext{
		OutPkg: outPkg,
	}
	ctx.ParseElem = func(t gotype.Type, elem gotype.Type) (Field, error) {
		return activeFieldRegistry.parse(t, elem, ctx)
	}

	field, err := activeFieldRegistry.parse(t, t.Elem(), ctx)
	if err != nil {
		return nil, err
	}

	if tagInfo != nil {
		if tagInfo.DBName != "" {
			field.setDBName(tagInfo.DBName)
		}
		if len(tagInfo.Indexes) > 0 {
			field.setIndexes(tagInfo.Indexes)
		}
		if tagInfo.Search != nil {
			field.setSearch(tagInfo.Search)
		}
	}

	return field, nil
}

// TagInfo holds all parsed som struct tag data.
type TagInfo struct {
	DBName  string
	Indexes []IndexInfo
	Search  *SearchInfo
}

// parseSomTag parses the "som" struct tag and extracts field metadata.
// All parameterized options use key=value syntax:
//
//	som:"index"
//	som:"index=my_index"
//	som:"unique"
//	som:"unique=composite_name"
//	som:"name=db_field_name"
//	som:"fulltext=english_search"
//	som:"index,unique=login"
func parseSomTag(tag string) (*TagInfo, error) {
	if tag == "" || tag == "in" || tag == "out" {
		return nil, nil
	}

	info := &TagInfo{}

	parts := strings.Split(tag, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		key, value, hasValue := strings.Cut(part, "=")

		switch key {
		case "index":
			idx := IndexInfo{}
			if hasValue {
				if value == "" {
					return nil, fmt.Errorf("invalid tag %q: index name must not be empty", part)
				}
				idx.Name = value
			}
			info.Indexes = append(info.Indexes, idx)

		case "unique":
			idx := IndexInfo{Unique: true}
			if hasValue {
				if value == "" {
					return nil, fmt.Errorf("invalid tag %q: unique name must not be empty", part)
				}
				idx.Name = value
			}
			info.Indexes = append(info.Indexes, idx)

		case "name":
			if !hasValue || value == "" {
				return nil, fmt.Errorf("invalid tag %q: name requires a value (name=db_field_name)", part)
			}
			if info.DBName != "" {
				return nil, fmt.Errorf("invalid tag: name specified multiple times")
			}
			info.DBName = value

		case "fulltext":
			if !hasValue || value == "" {
				return nil, fmt.Errorf("invalid tag %q: fulltext requires a config name (fulltext=english_search)", part)
			}
			info.Search = &SearchInfo{ConfigName: value}

		default:
			return nil, fmt.Errorf("unknown som tag %q", part)
		}
	}

	return info, nil
}

type Output struct {
	PkgPath    string
	Nodes      []*Node
	Edges      []*Edge
	Structs    []*Struct
	Enums      []*Enum
	EnumValues []*EnumValue

	Define *DefineOutput

	UsedFeatures *UsedFeatures
}

type UsedFeatures struct {
	UsesGoogleUUID bool
	UsesGofrsUUID  bool
}

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
