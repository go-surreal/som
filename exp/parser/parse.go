package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"golang.org/x/exp/maps"
)

func Parse() error {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, "../examples/basic/model", nil, 0)
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

	fmt.Println("package found:", pkg.Name)

	p := &Parser{
		structs: map[string]*ast.StructType{},
	}

	for _, file := range pkg.Files {

		for _, imp := range file.Imports {
			fmt.Println("import:", imp.Name, imp.Path.Value)
		}

		for _, decl := range file.Decls {
			parseDecl(p, decl)
		}

	}

	p.findNodes()

	return nil
}

func parseDecl(p *Parser, decl ast.Decl) {
	switch matchedDecl := decl.(type) {

	case *ast.GenDecl:
		{
			for _, spec := range matchedDecl.Specs {
				parseSpec(p, spec)
			}
		}
	}
}

func parseSpec(p *Parser, spec ast.Spec) {
	switch matchedSpec := spec.(type) {

	case *ast.TypeSpec:
		{
			parseTypeSpec(p, matchedSpec)
		}
	}
}

func parseTypeSpec(p *Parser, spec *ast.TypeSpec) {
	switch matchedSpec := spec.Type.(type) {

	case *ast.StructType:
		{
			p.structs[spec.Name.Name] = matchedSpec
		}

	}
}

func parseStructType(name string, t *ast.StructType) {
	var fields []string

	for _, field := range t.Fields.List {
		fld := ""

		for _, name := range field.Names {
			fld += " " + name.Name
		}

		parseExpr(fld, field.Type)

		fields = append(fields, fld)
	}

	fmt.Println("struct:", name, fields)
}

func parseExpr(name string, expr ast.Expr) {
	switch fieldType := expr.(type) {

	case *ast.Ident:
		{
			fmt.Println("IDENT:", fieldType.Name)
		}

	case *ast.SelectorExpr:
		{
			fmt.Println("SELECTOR EXPR:", fieldType.X, fieldType.Sel.Name, fieldType)
			parseExpr(name, fieldType.X)
		}

	case *ast.IndexListExpr:
		{
			fmt.Println("INDEX LIST EXPR:", fieldType.X, fieldType.Indices[0], fieldType.Indices[1])
			parseExpr(name, fieldType.X)
		}

	default:
		{
			fmt.Println("UNMAPPED:", "'"+name+"'", fmt.Sprintf("type: %#v", expr))
		}

	}
}
