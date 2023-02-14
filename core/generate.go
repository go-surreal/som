package core

import (
	"fmt"
	"github.com/marcbinz/som/core/codegen"
	"github.com/marcbinz/som/core/parser"
	"github.com/marcbinz/som/core/util"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Generate(inPath, outPath string) error {
	source, err := parser.Parse(inPath)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(outPath); err != nil {
		return err
	}

	absDir, err := filepath.Abs(outPath)
	if err != nil {
		return fmt.Errorf("could not find absolute path: %v", err)
	}

	pkgPath, modPath, err := util.ParseMod(absDir)
	if err != nil {
		return err
	}

	diff := strings.TrimPrefix(absDir, modPath)
	outPkg := path.Join(pkgPath, diff)

	err = codegen.Build(source, outPath, outPkg)
	if err != nil {
		return err
	}

	return nil
}
