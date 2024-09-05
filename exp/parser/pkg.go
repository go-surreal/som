package parser

import (
	"fmt"
	"go/ast"
)

type Package struct {
	Name string

	Structs []*Struct
}

func NewPackage(name string) *Package {
	return &Package{
		Name: name,
	}
}

func (p *Package) parse(files map[string]*ast.File) error {
	for _, file := range files {
		if err := p.parseFile(file); err != nil {
			return fmt.Errorf("could not parse file: %w", err)
		}
	}

	return nil
}

func (p *Package) parseFile(file *ast.File) error {
	//for _, imp := range file.Imports {
	//	fmt.Println("import:", imp.Name, imp.Path.Value)
	//}

	for _, decl := range file.Decls {
		if err := p.parseDecl(decl); err != nil {
			return fmt.Errorf("could not parse declaration: %w", err)
		}
	}

	return nil
}
