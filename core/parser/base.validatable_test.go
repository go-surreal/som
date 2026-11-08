package parser

import (
	"os"
	"testing"

	"github.com/wzshiming/gotype"
)

func loadValidatable(t *testing.T) gotype.Type {
	t.Helper()

	workDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}

	scope, err := gotype.NewImporter().Import("./testdata/validatable", workDir)
	if err != nil {
		t.Fatalf("could not import testdata: %v", err)
	}

	return scope
}

func TestHasValidateMethodOnDeclaredType(t *testing.T) {
	scope := loadValidatable(t)

	tests := map[string]bool{
		"ValueRecv":      true,
		"PtrRecv":        true,
		"NamedResult":    true,
		"Enum":           true,
		"Promoted":       true,
		"Plain":          false,
		"WrongResult":    false,
		"WrongParams":    false,
		"TooManyResults": false,
	}

	for name, want := range tests {
		typ, ok := scope.ChildByName(name)
		if !ok {
			t.Fatalf("type %s not found", name)
		}

		if got := hasValidateMethod(typ); got != want {
			t.Errorf("hasValidateMethod(%s) = %v, want %v", name, got, want)
		}
	}
}

// TestHasValidateMethodOnField checks the detection as it happens while
// parsing, where the type is reached through the field holding it.
func TestHasValidateMethodOnField(t *testing.T) {
	scope := loadValidatable(t)

	fields, ok := scope.ChildByName("Fields")
	if !ok {
		t.Fatal("type Fields not found")
	}

	tests := map[string]bool{
		"Value":    true,
		"Ptr":      true,
		"Enum":     true,
		"Promoted": true,
		"Plain":    false,
		"Imported": true,
		"Foreign":  false,
		"String":   false,
		"Slice":    false,
	}

	for name, want := range tests {
		field, ok := fields.FieldByName(name)
		if !ok {
			t.Fatalf("field %s not found", name)
		}

		elem := field.Elem()
		if elem.Kind() == gotype.Ptr {
			elem = elem.Elem()
		}

		if got := hasValidateMethod(elem); got != want {
			t.Errorf("hasValidateMethod(field %s) = %v, want %v", name, got, want)
		}
	}
}
