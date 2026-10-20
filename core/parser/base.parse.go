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

// Parse reads the model package at dir and turns it into an Output.
//
// typeHandlers and fieldHandlers are the prototypes that every declared type
// and every struct field is matched against, in the given order. The first
// match wins, so more specific handlers must come first.
func Parse(dir string, outPkg string, typeHandlers []TypeHandler, fieldHandlers []FieldHandler) (*Output, error) {
	res := &Output{}

	tReg := newTypeRegistry(typeHandlers)
	fReg := newFieldRegistry(fieldHandlers)

	imp := gotype.NewImporter()

	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not determine working directory: %w", err)
	}

	if !strings.HasPrefix(dir, "./") {
		dir = "./" + dir
	}

	n, err := imp.Import(dir, workDir)
	if err != nil {
		return nil, fmt.Errorf("could not import model package %s: %w", dir, err)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("could not find absolute path: %w", err)
	}

	mod, err := gomod.FindGoMod(absDir)
	if err != nil {
		return nil, fmt.Errorf("could not find go.mod for %s: %w", absDir, err)
	}

	diff := strings.TrimPrefix(absDir, mod.Dir())
	res.PkgPath = path.Join(mod.Module(), diff)

	ctx := &TypeContext{OutPkg: outPkg, PkgScope: n, Output: res, fields: fReg}

	nc := n.NumChild()
	for i := range nc {
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
