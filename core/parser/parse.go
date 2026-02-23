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
	var indexInfo *IndexInfo
	var searchInfo *SearchInfo
	if isStructField {
		somTag := t.Tag().Get("som")
		indexInfo, searchInfo = parseSomTag(somTag)
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

	if indexInfo != nil {
		field.setIndex(indexInfo)
	}
	if searchInfo != nil {
		field.setSearch(searchInfo)
	}

	return field, nil
}

// parseSomTag parses the "som" struct tag and extracts index/search info.
func parseSomTag(tag string) (index *IndexInfo, search *SearchInfo) {
	if tag == "" {
		return nil, nil
	}

	if tag == "in" || tag == "out" {
		return nil, nil
	}

	parts := strings.Split(tag, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.HasPrefix(part, "fulltext(") && strings.HasSuffix(part, ")") {
			configName := part[9 : len(part)-1]
			search = &SearchInfo{ConfigName: configName}
			continue
		}

		if part == "index" {
			if index == nil {
				index = &IndexInfo{}
			}
			continue
		}

		if part == "unique" || strings.HasPrefix(part, "unique(") {
			if index == nil {
				index = &IndexInfo{}
			}
			index.Unique = true

			if strings.HasPrefix(part, "unique(") && strings.HasSuffix(part, ")") {
				inner := part[7 : len(part)-1]
				index.UniqueName = inner
			}
			continue
		}

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
