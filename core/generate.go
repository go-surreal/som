package core

import (
	"github.com/marcbinz/sdb/core/codegen"
	"github.com/marcbinz/sdb/core/parser"
	"os"
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

	outPkg := strings.TrimSuffix(source.PkgPath, inPath) + outPath // TODO: is this really safe?

	err = codegen.Build(source, outPath, outPkg)
	if err != nil {
		return err
	}

	return nil
}
