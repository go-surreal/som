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

func (f *ComplexID) SchemaStatements(table, prefix string) []string {
	var statements []string

	for i, sub := range f.source.Fields {
		var dbType string
		switch sub.Field.(type) {
		case *parser.FieldString:
			dbType = "string"
		case *parser.FieldNumeric:
			dbType = "int|float"
		case *parser.FieldBool:
			dbType = "bool"
		case *parser.FieldTime:
			dbType = "datetime"
		case *parser.FieldDuration:
			dbType = "duration"
		case *parser.FieldUUID:
			dbType = "string"
		default:
			dbType = "any"
		}

		var fieldName string
		switch f.source.Kind {
		case parser.IDTypeArray:
			fieldName = fmt.Sprintf("id[%d]", i)
		case parser.IDTypeObject:
			fieldName = "id." + sub.DBName
		}

		statements = append(statements, fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+fieldName, table, dbType,
		))
	}

	return statements
}

func (f *ComplexID) CodeGen() *CodeGen {
	return &CodeGen{}
}
