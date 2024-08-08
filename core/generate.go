package core

import (
	"fmt"
	"github.com/go-surreal/som/core/codegen"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util"
	"github.com/go-surreal/som/core/util/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Generate(inPath, outPath string, dry bool) error {
	source, err := parser.Parse(inPath)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(outPath); err != nil {
		return err
	}

	absDir, err := filepath.Abs(outPath)
	if err != nil {
		return fmt.Errorf("could not find absolute path: %w", err)
	}

	pkgPath, modPath, err := util.ParseMod(absDir)
	if err != nil {
		return err
	}

	diff := strings.TrimPrefix(absDir, modPath)
	outPkg := path.Join(pkgPath, diff)

	out := fs.New()

	err = codegen.Build(source, out, outPkg)
	if err != nil {
		return err
	}

	if dry {
		return out.Dry(outPath)
	}

	return out.Flush(outPath)
}
