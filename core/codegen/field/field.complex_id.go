package field

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
)

type ComplexID struct {
	*baseField

	source *parser.FieldComplexID
}

func (f *ComplexID) typeGo() jen.Code {
	return jen.Qual(f.SourcePkg, f.source.StructName)
}

func (f *ComplexID) typeConv(_ Context) jen.Code {
	return f.typeGo()
}

func (f *ComplexID) TypeDatabase() string {
	switch f.source.Kind {
	case parser.IDTypeArray:
		return "array"
	case parser.IDTypeObject:
		return "object"
	default:
		return "any"
	}
}

func (f *ComplexID) SchemaStatements(table, _ string) []string {
	var typeDef string
	switch f.source.Kind {
	case parser.IDTypeArray:
		typeDef = fmt.Sprintf("array<any, %d>", len(f.source.Fields))
	case parser.IDTypeObject:
		typeDef = "object"
	default:
		return nil
	}
	return []string{
		fmt.Sprintf("DEFINE FIELD id ON TABLE %s TYPE %s;", table, typeDef),
	}
}

func (f *ComplexID) CodeGen() *CodeGen {
	return &CodeGen{}
}
