package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/parser"
)

// FieldDescriptorInfo holds the information needed to generate a field descriptor
// for the distinct values API.
type FieldDescriptorInfo struct {
	// TypeCode is the jen.Code representing the Go type of the field (e.g., jen.String()).
	TypeCode jen.Code
	// FactoryName is the name of the factory function in the query package
	// (e.g., "NewField", "NewTimeField", "NewUUIDGoogleField").
	FactoryName string
}

// FieldDescriptorFor returns the descriptor info for a field, or nil if the field type
// is not supported for distinct queries (e.g., slices, structs, nodes, edges).
func FieldDescriptorFor(f Field) *FieldDescriptorInfo {
	switch ft := f.(type) {
	case *String:
		return &FieldDescriptorInfo{
			TypeCode:    ft.typeGo(),
			FactoryName: "NewField",
		}

	case *Numeric:
		return &FieldDescriptorInfo{
			TypeCode:    ft.typeGo(),
			FactoryName: "NewField",
		}

	case *Bool:
		return &FieldDescriptorInfo{
			TypeCode:    ft.typeGo(),
			FactoryName: "NewField",
		}

	case *Enum:
		return &FieldDescriptorInfo{
			TypeCode:    ft.typeGo(),
			FactoryName: "NewField",
		}

	case *Time:
		if ft.baseField.source.Pointer() {
			return &FieldDescriptorInfo{
				TypeCode:    ft.typeGo(),
				FactoryName: "NewTimePtrField",
			}
		}
		return &FieldDescriptorInfo{
			TypeCode:    ft.typeGo(),
			FactoryName: "NewTimeField",
		}

	case *UUID:
		if ft.source.Package == parser.UUIDPackageGofrs {
			if ft.baseField.source.Pointer() {
				return &FieldDescriptorInfo{TypeCode: ft.typeGo(), FactoryName: "NewUUIDGofrsPtrField"}
			}
			return &FieldDescriptorInfo{TypeCode: ft.typeGo(), FactoryName: "NewUUIDGofrsField"}
		}
		if ft.baseField.source.Pointer() {
			return &FieldDescriptorInfo{TypeCode: ft.typeGo(), FactoryName: "NewUUIDGooglePtrField"}
		}
		return &FieldDescriptorInfo{TypeCode: ft.typeGo(), FactoryName: "NewUUIDGoogleField"}

	default:
		return nil
	}
}
