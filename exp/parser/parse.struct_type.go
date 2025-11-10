package parser

import (
	"fmt"
	"github.com/go-surreal/som/exp/def"
	"go/ast"
	"go/types"
)

func (p *Parser) parseStructType(spec *ast.TypeSpec, structType *ast.StructType, info *types.Info) error {
	if !spec.Name.IsExported() {
		return nil // ignore unexported structs
	}

	infoDef, ok := info.Defs[spec.Name]
	if !ok {
		return fmt.Errorf("could not find def for struct %s", spec.Name.Name)
	}

	structDef := &def.Struct{
		Base: &def.Base{
			Package: infoDef.Pkg().Path(),
			Name:    spec.Name.Name,
		},
	}

	if spec.TypeParams != nil {
		for _, typeParam := range spec.TypeParams.List {
			parseTypeParam(typeParam)
		}
	}

	for _, field := range structType.Fields.List {
		if len(field.Names) < 1 {
			continue // TODO: support embedded fields!
		}

		fields, err := parseField(field, info)
		if err != nil {
			return err
		}

		structDef.Fields = append(structDef.Fields, fields...)
	}

	switch {

	case isNode(structType):
		p.nodes = append(p.nodes, &def.Node{
			Struct: structDef,
		})

	default:
		p.structs = append(p.structs, structDef)
	}

	// spec.Name.Name

	//s := &Struct{
	//	Name: spec.Name.Name,
	//}
	//
	//for _, structField := range matchedSpec.Fields.List {
	//	if err := s.parseField(structField); err != nil {
	//		return fmt.Errorf("could not parse struct type: %w", err)
	//	}
	//}
	//
	//p.Structs = append(p.Structs, s)

	//return nil

	//structDef := &def.Struct{}
	//
	//for _, field := range spec.Fields.List {
	//	if err := parseField(field); err != nil {
	//		return fmt.Errorf("could not parse field: %w", err)
	//	}
	//}

	return nil
}

func parseTypeParam(typeParam *ast.Field) []*def.TypeParam {
	var typeParams []*def.TypeParam

	for _, name := range typeParam.Names {
		typeParams = append(typeParams, &def.TypeParam{
			Name: name.Name,
			// TODO: Type
		})
	}

	return typeParams
}

func isNode(structType *ast.StructType) bool {
	for _, field := range structType.Fields.List {
		if len(field.Names) > 0 {
			continue
		}

		if expr, ok := field.Type.(*ast.SelectorExpr); ok {
			x, ok := expr.X.(*ast.Ident)
			if !ok {
				continue
			}

			if x.Name == "som" && expr.Sel.Name == "Node" {
				return true
			}
		}
	}

	return false
}

func isEdge(structType *ast.StructType) bool {
	for _, field := range structType.Fields.List {
		if len(field.Names) > 0 {
			continue
		}

		if expr, ok := field.Type.(*ast.SelectorExpr); ok {
			x, ok := expr.X.(*ast.Ident)
			if !ok {
				continue
			}

			if x.Name == "som" && expr.Sel.Name == "Edge" {
				return true
			}
		}
	}

	return false
}
