package core

import (
	"github.com/marcbinz/sdb/builder"
	"github.com/marcbinz/sdb/parser"
)

func Generate(inPath, outPath string) error {
	res, err := parser.Parse(inPath)
	if err != nil {
		return err
	}

	err = builder.Build(res, outPath)
	if err != nil {
		return err
	}

	return nil
}
