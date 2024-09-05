package parser

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"golang.org/x/exp/maps"
)

//type Parser struct {
//	structs map[string]*ast.StructType
//}
//
//func (p *Parser) findNodes() {
//	for name, str := range p.structs {
//		for _, field := range str.Fields.List {
//			if len(field.Names) > 0 {
//				continue
//			}
//			fmt.Println("anonymous field:", name, field.Type)
//
//			if _, ok := field.Type.(*ast.SelectorExpr); ok {
//				fmt.Println("ok")
//			}
//		}
//	}
//}

type Parser struct {
	pkgs map[string]*Package
}

func NewParser() *Parser {
	return &Parser{
		pkgs: make(map[string]*Package),
	}
}

type x[T any] = string

func (p *Parser) Parse(path string) error {
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

	rawPkg := maps.Values(pkgs)[0]

	pkg := NewPackage(rawPkg.Name)

	p.pkgs[pkg.Name] = pkg

	if err := pkg.parse(rawPkg.Files); err != nil {
		return fmt.Errorf("could not parse package files: %w", err)
	}

	//aJSON, _ := json.MarshalIndent(p, "", "  ")
	//fmt.Println(string(aJSON))

	return nil
}
