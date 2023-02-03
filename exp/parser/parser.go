package parser

import (
	"fmt"
	"go/ast"
)

type Parser struct {
	structs map[string]*ast.StructType
}

func (p *Parser) findNodes() {
	for name, str := range p.structs {
		for _, field := range str.Fields.List {
			if len(field.Names) > 0 {
				continue
			}
			fmt.Println("anonymous field:", name, field.Type)

			if _, ok := field.Type.(*ast.SelectorExpr)
		}
	}
}
