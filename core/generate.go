package core

import (
	"github.com/marcbinz/som/core/codegen"
	"github.com/marcbinz/som/core/parser"
	"os"
	"path"
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

	outPkg := path.Join(strings.TrimSuffix(source.PkgPath, inPath), outPath) // TODO: is this really safe?

	err = codegen.Build(source, outPath, outPkg)
	if err != nil {
		return err
	}

	return nil
}
