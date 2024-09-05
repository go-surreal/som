package parser

import (
	"fmt"
	"go/ast"
)

func (p *Package) parseDecl(decl ast.Decl) error {
	switch decl := decl.(type) {

	case *ast.GenDecl: // import, type, const or var declaration
		{
			for _, spec := range decl.Specs {
				if err := parseSpec(p, spec); err != nil {
					return fmt.Errorf("could not parse spec: %w", err)
				}
			}

			return nil
		}

	case *ast.FuncDecl: // function declaration
		{
			return nil // we do not need those
		}

	case *ast.BadDecl: // invalid declaration
		{
			return nil // TODO: provide useful error messages
		}

	default:
		return fmt.Errorf("unsupported declaration: %T", decl)
	}
}
