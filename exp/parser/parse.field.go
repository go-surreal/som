package parser

import (
	"errors"
	"fmt"
	"github.com/go-surreal/som/exp/def/field"
	"go/ast"
	"go/types"
)

func parseField(raw *ast.Field, info *types.Info) ([]field.Field, error) {
	var fields []field.Field

	for _, name := range raw.Names {
		if !name.IsExported() {
			continue
		}

		fieldType, ok := info.Types[raw.Type]
		if !ok {
			return nil, errors.New("field type not found")
		}

		parsedField, ok, err := parseFieldType(name.Name, fieldType.Type)
		if err != nil {
			return nil, fmt.Errorf("could not parse field type: %w", err)
		}

		if !ok {
			continue // skip unsupported field types
		}

		fields = append(fields, parsedField)
	}

	return fields, nil
}

func parseFieldType(name string, fieldType types.Type) (field.Field, bool, error) {
	switch mappedField := fieldType.(type) {

	case *types.Alias:
		return parseFieldType(name, mappedField.Underlying()) // TODO: okay?

	case *types.Array:

	case *types.Basic:
		return parseBasicType(name, mappedField)

	case *types.Chan:
		return nil, false, errors.New("chan type is not supported")

	case *types.Interface:
		return nil, false, errors.New("any type is not yet supported") // TODO: support :)

	case *types.Named:
		{
			// TODO: check if underlying type is supported!

			return &field.Named{
				BaseField: &field.BaseField{Name: name},
				Pkg:       mappedField.Obj().Pkg().Path(),
				TypeName:  mappedField.Obj().Name(),
			}, true, nil
		}

	case *types.Map:
		return nil, false, errors.New("map type is not yet supported") // TODO: support :)

	case *types.Pointer:
		{
			ptrField, ok, err := parseFieldType("", mappedField.Elem())
			if err != nil {
				return nil, false, fmt.Errorf("could not parse pointer type: %w", err)
			}

			if !ok {
				return nil, false, nil // TODO: okay?
			}

			return &field.Pointer{
				BaseField: &field.BaseField{
					Name: name,
				},
				Field: ptrField,
			}, true, nil
		}

	case *types.Signature:
		return nil, false, errors.New("function type is not supported")

	case *types.Slice:
		{
			elemField, ok, err := parseFieldType("", mappedField.Elem())
			if err != nil {
				return nil, false, fmt.Errorf("could not parse slice elem type: %w", err)
			}

			if !ok {
				return nil, false, nil // TODO: okay?
			}

			return &field.Slice{
				BaseField: &field.BaseField{
					Name: name,
				},
				Elem: elemField,
			}, true, nil
		}

	case *types.Struct:
		return nil, false, errors.New("anonymous struct type is not supported")

	case *types.Tuple:
		// TODO: test!

	case *types.TypeParam:
		return nil, false, errors.New("generic type is not yet supported") // TODO: support?

	case *types.Union:
		return nil, false, errors.New("union type is not supported")

	default:
		return nil, false, fmt.Errorf("unknown field type: %T", fieldType)
	}

	return nil, false, nil // TODO
}

func parseBasicType(name string, basic *types.Basic) (field.Field, bool, error) {
	switch basic.Kind() {

	case types.Invalid:
		return nil, false, errors.New("invalid type")

	case types.Bool:
		{
			return &field.Bool{
				BaseField: &field.BaseField{Name: name},
			}, true, nil
		}

	case types.Int, types.Int8, types.Int16 /* types.Int32, */, types.Int64,
		types.Uint /* types.Uint8, */, types.Uint16, types.Uint32, types.Uint64,
		types.Uintptr,
		types.Float32, types.Float64,
		types.Complex64, types.Complex128:
		{
			return &field.Numeric{
				BaseField: &field.BaseField{Name: name},
				Kind:      basic.Kind(),
			}, true, nil
		}

	case types.String:
		{
			return &field.String{
				BaseField: &field.BaseField{Name: name},
			}, true, nil
		}

	case types.UnsafePointer:
		return nil, false, errors.New("unsafe pointer type is not supported")

	// types for untyped values
	case types.UntypedBool, types.UntypedInt, types.UntypedRune, types.UntypedFloat,
		types.UntypedComplex, types.UntypedString, types.UntypedNil:
		{
			return nil, false, errors.New("untyped type is not supported")
		}

	case types.Uint8: // either numeric or byte
		{
			if basic.Name() == "byte" {
				return &field.Byte{
					BaseField: &field.BaseField{Name: name},
				}, true, nil
			}

			return &field.Numeric{
				BaseField: &field.BaseField{Name: name},
				Kind:      basic.Kind(),
			}, true, nil
		}

	case types.Int32: // either numeric or rune
		{
			if basic.Name() == "rune" {
				return &field.Rune{
					BaseField: &field.BaseField{Name: name},
				}, true, nil
			}

			return &field.Numeric{
				BaseField: &field.BaseField{Name: name},
				Kind:      basic.Kind(),
			}, true, nil
		}

	default:
		return nil, false, fmt.Errorf("unknown basic type: %s", basic)
	}
}
