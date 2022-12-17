package codegen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/dbtype"
	"github.com/marcbinz/som/core/codegen/field"
)

type schemaBuilder struct {
	*baseBuilder
}

func newSchemaBuilder(input *input, basePath, basePkg, pkgName string) *schemaBuilder {
	return &schemaBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *schemaBuilder) build() error {
	if err := b.createDir(); err != nil {
		return err
	}

	var statements []string

	for _, node := range b.nodes {
		statement := fmt.Sprintf("DEFINE TABLE %s SCHEMAFULL", node.NameDatabase())
		statements = append(statements, statement)

		for _, f := range node.GetFields() {
			if f.NameDatabase() == "id" { // TODO: cleaner approach?
				continue
			}

			fieldType, ok := b.mapFieldType(f)
			if !ok {
				continue
			}

			statement := fmt.Sprintf(
				"DEFINE FIELD %s ON TABLE %s TYPE %s",
				f.NameDatabase(), node.NameDatabase(), fieldType,
			)
			statements = append(statements, statement)

			if slice, ok := f.(*field.Slice); ok {
				statement := fmt.Sprintf(
					"DEFINE FIELD %s.* ON TABLE %s TYPE %s",
					f.NameDatabase(), node.NameDatabase(), "record("+strcase.ToLowerCamel(slice.Value())+")",
				)
				statements = append(statements, statement)
			}
		}
	}

	for _, statement := range statements {
		fmt.Println(statement)
	}

	return nil
}

func (b *schemaBuilder) mapFieldType(f dbtype.Field) (string, bool) {
	switch mappedField := f.(type) {
	case *field.Bool:
		return "bool", true
	case *field.String:
		return "string", true
	case *field.Numeric:
		return mappedField.TypeDatabase(), true
	case *field.Time:
		return "datetime", true
	case *field.Enum:
		return "string", true
	case *field.Slice:
		return "array", true
	case *field.Struct:
		// mappedField.StructName()
		return "object", true
	case *field.Node:
		return "record(" + strcase.ToLowerCamel(mappedField.NodeName()) + ")", true
	}

	return "", false
}

type schemaData struct{}

const schemaTmpl = ``
