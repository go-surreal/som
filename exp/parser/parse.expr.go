package parser

import (
	"fmt"
	"go/ast"
)

func (p *Package) parseExpr(expr ast.Expr) error {
	switch matchedExpr := expr.(type) {

	case *ast.ArrayType:
		{
			return nil
		}

	case *ast.BasicLit:
		{
			return nil
		}

	case *ast.BinaryExpr:
		{
			return nil
		}

	case *ast.BadExpr, *ast.CallExpr, *ast.ChanType, *ast.CompositeLit, *ast.Ellipsis:
		{
			// map to unsupported type
			return nil
		}

	case *ast.FuncLit:
	case *ast.FuncType:
	case *ast.Ident:
		{
			return nil
		}
	case *ast.IndexExpr:
	case *ast.IndexListExpr:
	case *ast.InterfaceType:
	case *ast.KeyValueExpr:
	case *ast.MapType:
	case *ast.ParenExpr:
	case *ast.SelectorExpr:
	case *ast.SliceExpr:
	case *ast.StarExpr:
	case *ast.StructType:
	case *ast.TypeAssertExpr:
	case *ast.UnaryExpr:

	default:
		return fmt.Errorf("unsupported expression: %T", matchedExpr)
	}
}
