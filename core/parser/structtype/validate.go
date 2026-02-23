package structtype

import (
	"fmt"

	"github.com/go-surreal/som/core/parser"
)

func validateFields(typeName string, fields []parser.Field, output *parser.Output) error {
	for _, f := range fields {
		if err := validateField(typeName, f, output); err != nil {
			return err
		}
	}
	return nil
}

func validateField(typeName string, f parser.Field, output *parser.Output) error {
	if s := f.Search(); s != nil {
		if !searchExists(s.ConfigName, output) {
			return fmt.Errorf("%s: field %q references unknown search config %q", typeName, f.Name(), s.ConfigName)
		}
	}

	switch field := f.(type) {
	case *parser.FieldNode:
		if !nodeExists(field.Node, output) {
			return fmt.Errorf("%s: field %q references unknown node %q", typeName, f.Name(), field.Node)
		}
	case *parser.FieldEdge:
		if !edgeExists(field.Edge, output) {
			return fmt.Errorf("%s: field %q references unknown edge %q", typeName, f.Name(), field.Edge)
		}
	case *parser.FieldEnum:
		if !enumExists(field.Typ, output) {
			return fmt.Errorf("%s: field %q references unknown enum %q", typeName, f.Name(), field.Typ)
		}
	case *parser.FieldStruct:
		if !structExists(field.Struct, output) {
			return fmt.Errorf("%s: field %q references unknown struct %q", typeName, f.Name(), field.Struct)
		}
	case *parser.FieldSlice:
		return validateField(typeName, field.Field, output)
	case *parser.FieldComplexID:
		for _, cf := range field.Fields {
			if err := validateField(typeName, cf.Field, output); err != nil {
				return err
			}
		}
	}

	return nil
}

func nodeExists(name string, output *parser.Output) bool {
	for _, n := range output.Nodes {
		if n.Name == name {
			return true
		}
	}
	return false
}

func edgeExists(name string, output *parser.Output) bool {
	for _, e := range output.Edges {
		if e.Name == name {
			return true
		}
	}
	return false
}

func enumExists(name string, output *parser.Output) bool {
	for _, e := range output.Enums {
		if e.Name == name {
			return true
		}
	}
	return false
}

func structExists(name string, output *parser.Output) bool {
	for _, s := range output.Structs {
		if s.Name == name {
			return true
		}
	}
	return false
}

func searchExists(name string, output *parser.Output) bool {
	if output.Define == nil {
		return false
	}
	for _, s := range output.Define.Searches {
		if s.Name == name {
			return true
		}
	}
	return false
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
