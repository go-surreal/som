package parser

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

// validateFields runs the post-parse checks shared by every model type: the
// database names must be unique and must not collide with a reserved name, and
// each field verifies itself against the fully parsed output.
func validateFields(typeName string, fields []Field, output *Output) error {
	seen := make(map[string]string, len(fields)) // dbName -> Go field name
	for _, f := range fields {
		dbName := f.DBName()
		if dbName == "" {
			dbName = strcase.ToSnake(f.Name())
		}
		if !f.IsInternal() && ReservedDBNames[dbName] {
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

func nodeExists(name string, out *Output) bool {
	for _, n := range out.Nodes {
		if n.Name == name {
			return true
		}
	}
	return false
}

func edgeExists(name string, out *Output) bool {
	for _, e := range out.Edges {
		if e.Name == name {
			return true
		}
	}
	return false
}

func enumExists(name string, out *Output) bool {
	for _, e := range out.Enums {
		if e.Name == name {
			return true
		}
	}
	return false
}

func structExists(name string, out *Output) bool {
	for _, s := range out.Structs {
		if s.Name == name {
			return true
		}
	}
	return false
}

func searchExists(name string, out *Output) bool {
	if out.Define == nil {
		return false
	}
	for _, s := range out.Define.Searches {
		if s.Name == name {
			return true
		}
	}
	return false
}
