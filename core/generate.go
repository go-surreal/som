package core

import (
	"github.com/marcbinz/sdb/core/codegen"
	"github.com/marcbinz/sdb/core/parser"
	"os"
)

func Generate(inPath, outPath string) error {
	res, err := parser.Parse(inPath)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(outPath); err != nil {
		return err
	}

	err = codegen.Build(res, outPath)
	if err != nil {
		return err
	}

	return nil
}
