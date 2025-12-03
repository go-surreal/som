package field

import (
	"testing"

	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/parser"
	"gotest.tools/v3/assert"
)

func TestUUID_uuidPkg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		pkg      parser.UUIDPackage
		expected string
	}{
		{
			name:     "Google UUID package",
			pkg:      parser.UUIDPackageGoogle,
			expected: def.PkgUUIDGoogle,
		},
		{
			name:     "Gofrs UUID package",
			pkg:      parser.UUIDPackageGofrs,
			expected: def.PkgUUIDGofrs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			field := &UUID{
				source: &parser.FieldUUID{
					Package: tt.pkg,
				},
			}

			assert.Equal(t, tt.expected, field.uuidPkg())
		})
	}
}

func TestUUID_uuidTypeName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		pkg      parser.UUIDPackage
		expected string
	}{
		{
			name:     "Google UUID type name",
			pkg:      parser.UUIDPackageGoogle,
			expected: "UUIDGoogle",
		},
		{
			name:     "Gofrs UUID type name",
			pkg:      parser.UUIDPackageGofrs,
			expected: "UUIDGofrs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			field := &UUID{
				source: &parser.FieldUUID{
					Package: tt.pkg,
				},
			}

			assert.Equal(t, tt.expected, field.uuidTypeName())
		})
	}
}

func TestUUID_TypeDatabase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		pointer  bool
		expected string
	}{
		{
			name:     "non-pointer UUID",
			pointer:  false,
			expected: "uuid",
		},
		{
			name:     "pointer UUID",
			pointer:  true,
			expected: "option<uuid>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			field := &UUID{
				source: parser.NewFieldUUID("test_uuid", tt.pointer, parser.UUIDPackageGoogle),
			}

			assert.Equal(t, tt.expected, field.TypeDatabase())
		})
	}
}
