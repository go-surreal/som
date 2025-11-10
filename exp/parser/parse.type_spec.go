package parser

import (
	"fmt"
	"go/ast"
	"go/types"
)

func (p *Parser) parseType(spec *ast.TypeSpec, info *types.Info) error {
	switch matchedSpec := spec.Type.(type) {

	// *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes

	case *ast.Ident:
		if err := parseIdent(spec, matchedSpec); err != nil {
			return fmt.Errorf("could not parse ident: %w", err)
		}

	case *ast.ParenExpr:
		fmt.Println("PAREN EXPR:", spec.Name.Name, matchedSpec)
	case *ast.SelectorExpr:
		fmt.Println("SELECTOR EXPR:", spec.Name.Name, matchedSpec.X, matchedSpec.Sel.Name)
	case *ast.StarExpr:
		fmt.Println("STAR EXPR:", spec.Name.Name, matchedSpec.X)
	case *ast.ArrayType:
		fmt.Println("ARRAY TYPE:", spec.Name.Name, matchedSpec)
	case *ast.ChanType:
		fmt.Println("CHAN TYPE:", spec.Name.Name, matchedSpec)
	case *ast.FuncType:
		fmt.Println("FUNC TYPE:", spec.Name.Name, matchedSpec)
	case *ast.InterfaceType:
		fmt.Println("INTERFACE TYPE:", spec.Name.Name, matchedSpec)
	case *ast.MapType:
		fmt.Println("MAP TYPE:", spec.Name.Name, matchedSpec)

	case *ast.StructType:
		if err := p.parseStructType(spec, matchedSpec, info); err != nil {
			return fmt.Errorf("could not parse struct type: %w", err)
		}

	default:
		//return fmt.Errorf("unsupported type spec: %T", matchedSpec)
	}

	return nil
}

func parseIdent(spec *ast.TypeSpec, ident *ast.Ident) error {
	if !spec.Name.IsExported() {
		return nil
	}

	fmt.Println("IDENT:", spec.Name.Obj.Decl, spec.Name, spec.Name.IsExported(), ident.Name)

	return nil
}

func parseValueSpec(spec *ast.ValueSpec) error {
	switch matchedSpec := spec.Type.(type) {

	default:
		return fmt.Errorf("unsupported value spec: %T", matchedSpec)
	}
}

//func parseFieldType(name string, expr ast.Expr) (types.Field, error) {
//	base := &types.BaseField{
//		Name: name,
//	}
//
//	switch fieldType := expr.(type) {
//
//	case *ast.Ident:
//		{
//			fmt.Println("IDENT:", base.Name, fieldType.Name)
//			return nil, nil
//		}
//
//	case *ast.SelectorExpr:
//		{
//			fmt.Println("SELECTOR EXPR:", base.Name, fieldType, fieldType.X, fieldType.Sel.Name, fieldType)
//			//parseExpr(name, fieldType.X)
//			return nil, nil
//		}
//
//	//case *ast.IndexListExpr:
//	//	{
//	//		fmt.Println("INDEX LIST EXPR:", fieldType.X, fieldType.Indices[0], fieldType.Indices[1])
//	//		//parseExpr(name, fieldType.X)
//	//		return nil, nil
//	//	}
//
//	case *ast.ArrayType:
//		{
//			fmt.Println("ARRAY TYPE:", base.Name, fieldType)
//			fmt.Println(parseFieldType("", fieldType.Elt))
//			return nil, nil
//		}
//
//	case *ast.StarExpr:
//		{
//			fld, err := parseFieldType(name, fieldType.X)
//			if err != nil {
//				return nil, err
//			}
//
//			return fld, nil
//		}
//
//	default:
//		return nil, fmt.Errorf("unsupported field type: %T", fieldType)
//	}
//}

//func (s *Struct) parseField(field *ast.Field) error {
//	//fmt.Println("field:", field.Names, field.Type)
//
//	if len(field.Names) < 1 {
//		fld, err := parseFieldType("", field.Type)
//		if err != nil {
//			return fmt.Errorf("could not parse field type: %w", err)
//		}
//
//		s.Fields = append(s.Fields, fld)
//		return nil
//	}
//
//	for _, ident := range field.Names {
//		if !ident.IsExported() {
//			continue // skip unexported fields
//		}
//
//		fld, err := parseFieldType(ident.Name, field.Type)
//		if err != nil {
//			return fmt.Errorf("could not parse field type: %w", err)
//		}
//
//		s.Fields = append(s.Fields, fld)
//	}
//
//	return nil
//}
