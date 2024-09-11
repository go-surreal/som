package parser

import (
	"github.com/go-surreal/som/exp/def"
	"go/ast"
)

func (p *Parser) parseImportSpec(spec *ast.ImportSpec) {
	importDef := &def.Import{
		Path: spec.Path.Value,
	}

	if spec.Name != nil {
		importDef.Name = spec.Name.Name
	}

	p.imports = append(p.imports, importDef)
}
