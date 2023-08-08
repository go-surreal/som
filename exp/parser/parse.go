package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/marcbinz/som/exp/field"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
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

	p := &Package{
		Name: pkg.Name,
	}

	for _, file := range pkg.Files {

		//for _, imp := range file.Imports {
		//	fmt.Println("import:", imp.Name, imp.Path.Value)
		//}

		for _, decl := range file.Decls {
			if err := parseDecl(p, decl); err != nil {
				return fmt.Errorf("could not parse declaration: %w", err)
			}
		}

	}

	aJSON, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(aJSON))

	return nil
}

func parseDecl(p *Package, decl ast.Decl) error {
	switch matchedDecl := decl.(type) {

	case *ast.GenDecl:
		{
			for _, spec := range matchedDecl.Specs {
				if err := parseSpec(p, spec); err != nil {
					return fmt.Errorf("could not parse spec: %w", err)
				}
			}

			return nil
		}

	default:
		return fmt.Errorf("unsupported declaration: %T", matchedDecl)
	}
}

func parseSpec(p *Package, spec ast.Spec) error {
	switch matchedSpec := spec.(type) {

	case *ast.TypeSpec: // type declaration
		{
			err := parseTypeSpec(p, matchedSpec)
			if err != nil {
				return fmt.Errorf("could not parse type spec: %w", err)
			}

			return nil
		}

	case *ast.ValueSpec: // const or var declaration
		{
			parseValueSpec(p, matchedSpec)
			return nil
		}

	default:
		return nil // fmt.Errorf("unsupported spec: %T", matchedSpec)
	}
}

func parseTypeSpec(p *Package, spec *ast.TypeSpec) error {
	switch matchedSpec := spec.Type.(type) {

	case *ast.StructType:
		{
			s := &Struct{
				Name: spec.Name.Name,
			}

			//p.structs[spec.Name.Name] = matchedSpec
			err := parseStructType(s, matchedSpec)
			if err != nil {
				return fmt.Errorf("could not parse struct type: %w", err)
			}

			p.Structs = append(p.Structs, s)

			return nil
		}

	default:
		return fmt.Errorf("unsupported type spec: %T", matchedSpec)
	}
}

func parseValueSpec(p *Package, spec *ast.ValueSpec) error {
	switch matchedSpec := spec.Type.(type) {

	default:
		return fmt.Errorf("unsupported value spec: %T", matchedSpec)
	}
}

func parseStructType(s *Struct, t *ast.StructType) error {
	for _, rawField := range t.Fields.List {
		//fmt.Println("field:", field.Names, field.Type)

		if len(rawField.Names) < 1 {
			fld, err := parseFieldType("", rawField.Type)
			if err != nil {
				return fmt.Errorf("could not parse field type: %w", err)
			}

			s.Fields = append(s.Fields, fld)
			continue
		}

		for _, ident := range rawField.Names {
			if !ident.IsExported() {
				continue // skip unexported fields
			}

			fld, err := parseFieldType(ident.Name, rawField.Type)
			if err != nil {
				return fmt.Errorf("could not parse field type: %w", err)
			}

			s.Fields = append(s.Fields, fld)
		}
	}

	return nil
}

func parseFieldType(name string, expr ast.Expr) (field.Field, error) {
	base := &field.BaseField{
		Name: name,
	}

	switch fieldType := expr.(type) {

	case *ast.Ident:
		{
			fmt.Println("IDENT:", base.Name, fieldType.Name)
			return nil, nil
		}

	case *ast.SelectorExpr:
		{
			fmt.Println("SELECTOR EXPR:", base.Name, fieldType, fieldType.X, fieldType.Sel.Name, fieldType)
			//parseExpr(name, fieldType.X)
			return nil, nil
		}

	//case *ast.IndexListExpr:
	//	{
	//		fmt.Println("INDEX LIST EXPR:", fieldType.X, fieldType.Indices[0], fieldType.Indices[1])
	//		//parseExpr(name, fieldType.X)
	//		return nil, nil
	//	}

	case *ast.ArrayType:
		{
			fmt.Println("ARRAY TYPE:", base.Name, fieldType)
			fmt.Println(parseFieldType("", fieldType.Elt))
			return nil, nil
		}

	case *ast.StarExpr:
		{
			fld, err := parseFieldType(name, fieldType.X)
			if err != nil {
				return nil, err
			}

			return fld, nil
		}

	default:
		return nil, fmt.Errorf("unsupported field type: %T", fieldType)
	}
}

type Package struct {
	Name string

	Structs []*Struct
}

type Struct struct {
	Name string

	Fields []field.Field
}
