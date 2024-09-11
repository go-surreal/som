package parser

import (
	"errors"
	"fmt"
	"github.com/go-surreal/som/exp/def"
	"go/ast"
	"log/slog"
)

//go:autowire
func parseField(raw *ast.Field) ([]def.Field, error) {
	var fields []def.Field

	for _, name := range raw.Names {
		if !name.IsExported() {
			continue
		}

		field, ok, err := parseFieldType(name.Name, raw.Type)
		if err != nil {
			return nil, fmt.Errorf("could not parse field type: %w", err)
		}

		if !ok {
			continue // skip unsupported field types
		}

		fields = append(fields, field)
	}

	return fields, nil
}

func parseFieldType(name string, expr ast.Expr) (def.Field, bool, error) {
	switch field := expr.(type) {
	case *ast.ArrayType:
	case *ast.BadExpr:
	case *ast.BasicLit:
	case *ast.BinaryExpr:
	case *ast.CallExpr:
	case *ast.ChanType:
	case *ast.CompositeLit:

	case *ast.FuncLit:

	case *ast.Ident:
		return parseIdentType(name, field)

	case *ast.IndexExpr:
	case *ast.IndexListExpr:
	case *ast.InterfaceType:
	case *ast.KeyValueExpr:

	case *ast.MapType:
		return nil, false, errors.New("map type is not yet supported") // TODO: support :)

	case *ast.ParenExpr:
	case *ast.SelectorExpr:
	case *ast.SliceExpr:
	case *ast.StarExpr:
		{
			ptrField, ok, err := parseFieldType("", field.X)
			if err != nil {
				return nil, false, fmt.Errorf("could not parse pointer type: %w", err)
			}

			if !ok {
				return nil, false, nil // TODO: okay?
			}

			return &def.Pointer{
				BaseField: &def.BaseField{
					Name: name,
				},
				Field: ptrField,
			}, true, nil
		}

	case *ast.StructType:
		return nil, false, errors.New("anonymous struct type is not supported")

	case *ast.TypeAssertExpr:
		return nil, false, errors.New("type assertion is not supported")

	case *ast.UnaryExpr:

	case *ast.Ellipsis, *ast.FuncType:
		slog.Error("field type: %T", expr)

	default:
		return nil, false, fmt.Errorf("unknown field type: %T", expr)
	}

	return nil, false, nil // TODO
}

func parseIdentType(name string, ident *ast.Ident) (def.Field, bool, error) {
	fmt.Println(ident.Name)

	//return &def.String{
	//	BaseField: &def.BaseField{Name: name},
	//}, true, nil

	return nil, false, nil
}
