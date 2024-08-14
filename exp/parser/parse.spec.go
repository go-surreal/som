package parser

import (
	"fmt"
	"github.com/go-surreal/som/exp/field"
	"go/ast"
)

func parseSpec(p *Package, spec ast.Spec) error {
	switch spec := spec.(type) {

	case *ast.TypeSpec: // type declaration
		{
			if err := p.parseTypeSpec(spec); err != nil {
				return fmt.Errorf("could not parse type spec: %w", err)
			}

			return nil
		}

	case *ast.ValueSpec: // const or var declaration
		{
			parseValueSpec(p, spec)
			return nil
		}

	case *ast.ImportSpec: // import declaration
		{
			return nil
		}

	default:
		return fmt.Errorf("unsupported spec: %T", spec)
	}
}

func (p *Package) parseTypeSpec(spec *ast.TypeSpec) error {
	switch matchedSpec := spec.Type.(type) {

	case *ast.StructType:
		{
			s := &Struct{
				Name: spec.Name.Name,
			}

			for _, structField := range matchedSpec.Fields.List {
				if err := s.parseField(structField); err != nil {
					return fmt.Errorf("could not parse struct type: %w", err)
				}
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

type Struct struct {
	Name string

	Fields []field.Field
}

func (s *Struct) parseField(field *ast.Field) error {
	//fmt.Println("field:", field.Names, field.Type)

	if len(field.Names) < 1 {
		fld, err := parseFieldType("", field.Type)
		if err != nil {
			return fmt.Errorf("could not parse field type: %w", err)
		}

		s.Fields = append(s.Fields, fld)
		return nil
	}

	for _, ident := range field.Names {
		if !ident.IsExported() {
			continue // skip unexported fields
		}

		fld, err := parseFieldType(ident.Name, field.Type)
		if err != nil {
			return fmt.Errorf("could not parse field type: %w", err)
		}

		s.Fields = append(s.Fields, fld)
	}

	return nil
}
