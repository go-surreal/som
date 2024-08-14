package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"golang.org/x/exp/maps"
)

func Parse(path string) error {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		return fmt.Errorf("could not parse code in source path: %w", err)
	}

	if len(pkgs) < 1 {
		return errors.New("no packages found in source path")
	}

	if len(pkgs) > 1 {
		return errors.New("more than one package found in source path")
	}

	pkg := maps.Values(pkgs)[0]

	p := NewPackage(pkg.Name)

	if err := p.parse(pkg.Files); err != nil {
		return err
	}

	aJSON, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(aJSON))

	return nil
}
