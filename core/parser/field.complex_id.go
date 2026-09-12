package parser

import (
	"fmt"
	"go/ast"

	"github.com/iancoleman/strcase"
	"github.com/wzshiming/gotype"
)

type ComplexIDField struct {
	Name   string
	DBName string
	Field  Field
}

// FieldComplexID is a record id made up of several values, declared as a key
// struct embedding som.ArrayID or som.ObjectID. It is never matched against a
// source field, the node type creates it from its som.Node type argument.
type FieldComplexID struct {
	fieldBase
	Kind       IDType
	StructName string
	Fields     []ComplexIDField
}

func (f *FieldComplexID) HasNodeRef() bool {
	for _, sf := range f.Fields {
		if _, ok := sf.Field.(*FieldNode); ok {
			return true
		}
	}
	return false
}

func (f *FieldComplexID) Verify(out *Output) error {
	if err := f.fieldBase.Verify(out); err != nil {
		return err
	}

	for _, sf := range f.Fields {
		if err := sf.Field.Verify(out); err != nil {
			return err
		}
	}

	return nil
}

func (f *FieldComplexID) IsInternal() bool {
	return true
}

// ParseComplexIDFields resolves the key struct given as the type argument of a
// som.Node embed, e.g. PersonID for som.Node[PersonID].
func ParseComplexIDFields(t gotype.Type, ctx *TypeContext) (*FieldComplexID, error) {
	outPkg := ctx.OutPkg

	origin := t.Origin()
	if origin == nil {
		return nil, fmt.Errorf("complex ID type has no AST origin")
	}

	indexExpr, ok := origin.(*ast.IndexExpr)
	if !ok {
		return nil, fmt.Errorf("complex ID type is not a generic instantiation (expected IndexExpr, got %T)", origin)
	}

	var typeArgName string
	switch idx := indexExpr.Index.(type) {
	case *ast.Ident:
		typeArgName = idx.Name
	case *ast.SelectorExpr:
		typeArgName = idx.Sel.Name
	default:
		return nil, fmt.Errorf("complex ID type argument: unsupported AST node %T", indexExpr.Index)
	}

	keyType, ok := ctx.PkgScope.ChildByName(typeArgName)
	if !ok {
		return nil, fmt.Errorf("complex ID type argument %s not found in package scope", typeArgName)
	}

	if keyType.Kind() != gotype.Struct {
		return nil, fmt.Errorf("complex ID type parameter must be a struct, got %s", keyType.Kind())
	}

	structName := keyType.Name()
	nf := keyType.NumField()
	fields := make([]ComplexIDField, 0, nf)

	var kind IDType
	var kindSet bool
	for i := range nf {
		sf := keyType.Field(i)

		if sf.IsAnonymous() {
			if sf.Elem().PkgPath() == outPkg {
				switch sf.Name() {
				case "ArrayID":
					if kindSet {
						return nil, fmt.Errorf("complex ID struct %s embeds both ArrayID and ObjectID", structName)
					}
					kind = IDTypeArray
					kindSet = true
				case "ObjectID":
					if kindSet {
						return nil, fmt.Errorf("complex ID struct %s embeds both ArrayID and ObjectID", structName)
					}
					kind = IDTypeObject
					kindSet = true
				}
			}
			continue
		}

		if !ast.IsExported(sf.Name()) {
			continue
		}

		parsed, err := ctx.ParseUntaggedField(sf)
		if err != nil {
			return nil, fmt.Errorf("complex ID field %s: %w", sf.Name(), err)
		}

		switch parsed.(type) {
		case *FieldString, *FieldNumeric, *FieldBool, *FieldTime, *FieldDuration, *FieldUUID, *FieldNode:
		default:
			return nil, fmt.Errorf("complex ID field %s: unsupported type %T (only string, numeric, bool, time.Time, time.Duration, UUID, and node references are allowed)", sf.Name(), parsed)
		}

		fields = append(fields, ComplexIDField{
			Name:   sf.Name(),
			DBName: strcase.ToSnake(sf.Name()),
			Field:  parsed,
		})
	}

	if !kindSet {
		return nil, fmt.Errorf("complex ID struct %s must embed exactly one of som.ArrayID or som.ObjectID", structName)
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("complex ID struct %s has no exported fields", structName)
	}

	return &FieldComplexID{fieldBase: newBase("ID"), Kind: kind, StructName: structName, Fields: fields}, nil
}
