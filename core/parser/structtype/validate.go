package structtype

import (
	"fmt"

	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
)

func validateFields(typeName string, fields []parser.Field, output *parser.Output) error {
	seen := make(map[string]string, len(fields)) // dbName -> Go field name
	for _, f := range fields {
		dbName := f.DBName()
		if dbName == "" {
			dbName = strcase.ToSnake(f.Name())
		}
		if !f.IsInternal() && parser.ReservedDBNames[dbName] {
			return fmt.Errorf("%s: field %q maps to reserved database name %q", typeName, f.Name(), dbName)
		}
		if prev, ok := seen[dbName]; ok {
			return fmt.Errorf("%s: fields %q and %q both map to database name %q", typeName, prev, f.Name(), dbName)
		}
		seen[dbName] = f.Name()

		if err := f.Verify(output); err != nil {
			return fmt.Errorf("%s: %w", typeName, err)
		}
	}
	return nil
}

func hasDuplicates(names []string) (string, bool) {
	seen := make(map[string]struct{}, len(names))
	for _, name := range names {
		if _, ok := seen[name]; ok {
			return name, true
		}
		seen[name] = struct{}{}
	}
	return "", false
}
