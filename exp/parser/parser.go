package parser

import (
	"errors"
	"fmt"
	"github.com/go-surreal/som/exp/def"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"strings"
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
	imports []*def.Import
	nodes   []*def.Node
	edges   []*def.Edge
	structs []*def.Struct
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(path string) error {
	fileSet := token.NewFileSet()

	pkgs, err := parser.ParseDir(fileSet, path,
		func(info os.FileInfo) bool {
			return strings.HasSuffix(info.Name(), ".go")
		},
		parser.AllErrors,
	)
	if err != nil {
		return fmt.Errorf("could not parse code in source path: %w", err)
	}

	if len(pkgs) < 1 {
		return errors.New("no packages found in source path")
	}

	if len(pkgs) > 1 {
		return errors.New("more than one package found in source path")
	}

	for pkg := range maps.Values(pkgs) {
		if err := p.parse(pkg); err != nil {
			return fmt.Errorf("could not parse package: %w", err)
		}
	}

	return nil
}

func (p *Parser) parse(pkg ast.Node) error {
	// TODO: use ast.Walk() instead?

	for node := range ast.Preorder(pkg) {
		switch mappedNode := node.(type) {

		case *ast.TypeSpec:
			if err := p.parseTypeSpec(mappedNode); err != nil {
				return fmt.Errorf("could not parse type spec: %w", err)
			}

		case *ast.ImportSpec:
			p.parseImportSpec(mappedNode)
		}
	}

	return nil
}

func (p *Parser) String() string {
	out := "\n"

	out += "---- Imports ----\n\n"

	for _, imp := range p.imports {
		out += fmt.Sprintf("%s\n", imp)
	}

	out += "\n---- Structs ----\n\n"

	for _, str := range p.structs {
		out += fmt.Sprintf("%s\n", str)
	}

	out += "\n---- Nodes ----\n\n"

	for _, node := range p.nodes {
		out += fmt.Sprintf("%s\n", node)
	}

	out += "\n---- Edges ----\n\n"

	for _, edge := range p.edges {
		out += fmt.Sprintf("%s\n", edge)
	}

	return out
}
