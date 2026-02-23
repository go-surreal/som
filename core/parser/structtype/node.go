package structtype

import (
	"fmt"
	"go/ast"
	"path"
	"strings"

	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
	"github.com/wzshiming/gotype"
)

type NodeHandler struct{}

func (h *NodeHandler) Match(t gotype.Type, ctx *parser.TypeContext) bool {
	return IsNode(t, ctx.OutPkg)
}

func (h *NodeHandler) Handle(t gotype.Type, ctx *parser.TypeContext) error {
	node, err := ParseNode(t, ctx.OutPkg, ctx.PkgScope)
	if err != nil {
		return err
	}
	ctx.Output.Nodes = append(ctx.Output.Nodes, node)
	return nil
}

func (h *NodeHandler) Validate(_ *parser.TypeContext) error { return nil }

func IsNode(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)

		if !f.IsAnonymous() {
			continue
		}

		if f.Name() == "Node" && IsGenericNodeFromSom(f.Elem(), outPkg, "Node") {
			return true
		}
	}

	return false
}

func IsGenericNodeFromSom(t gotype.Type, outPkg string, name string) bool {
	if pkgPath := t.PkgPath(); pkgPath != "" {
		return pkgPath == outPkg
	}

	origin := t.Origin()
	if origin == nil {
		return false
	}
	indexExpr, ok := origin.(*ast.IndexExpr)
	if !ok {
		return false
	}
	selExpr, ok := indexExpr.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if selExpr.Sel.Name != name {
		return false
	}
	ident, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == path.Base(outPkg)
}

func IsKnownStringIDType(t gotype.Type) bool {
	origin := t.Origin()
	if origin == nil {
		return false
	}
	indexExpr, ok := origin.(*ast.IndexExpr)
	if !ok {
		return true
	}
	_, ok = indexExpr.Index.(*ast.SelectorExpr)
	return ok
}

func ParseIDType(t gotype.Type) parser.IDType {
	origin := t.Origin()
	if origin == nil {
		return parser.IDTypeULID
	}

	if indexExpr, ok := origin.(*ast.IndexExpr); ok {
		if selExpr, ok := indexExpr.Index.(*ast.SelectorExpr); ok {
			switch selExpr.Sel.Name {
			case "UUID":
				return parser.IDTypeUUID
			case "Rand":
				return parser.IDTypeRand
			case "ULID":
				return parser.IDTypeULID
			}
		}
	}

	return parser.IDTypeULID
}

func ParseComplexIDFields(t gotype.Type, outPkg string, pkgScope gotype.Type) (*parser.FieldComplexID, error) {
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

	keyType, ok := pkgScope.ChildByName(typeArgName)
	if !ok {
		return nil, fmt.Errorf("complex ID type argument %s not found in package scope", typeArgName)
	}

	if keyType.Kind() != gotype.Struct {
		return nil, fmt.Errorf("complex ID type parameter must be a struct, got %s", keyType.Kind())
	}

	structName := keyType.Name()
	nf := keyType.NumField()
	fields := make([]parser.ComplexIDField, 0, nf)

	var kind parser.IDType
	var kindSet bool
	for i := 0; i < nf; i++ {
		sf := keyType.Field(i)

		if sf.IsAnonymous() {
			if sf.Elem().PkgPath() == outPkg {
				switch sf.Name() {
				case "ArrayID":
					if kindSet {
						return nil, fmt.Errorf("complex ID struct %s embeds both ArrayID and ObjectID", structName)
					}
					kind = parser.IDTypeArray
					kindSet = true
				case "ObjectID":
					if kindSet {
						return nil, fmt.Errorf("complex ID struct %s embeds both ArrayID and ObjectID", structName)
					}
					kind = parser.IDTypeObject
					kindSet = true
				}
			}
			continue
		}

		if !ast.IsExported(sf.Name()) {
			continue
		}

		parsed, err := parser.ParseFieldInternal(sf, outPkg, false)
		if err != nil {
			return nil, fmt.Errorf("complex ID field %s: %w", sf.Name(), err)
		}

		switch parsed.(type) {
		case *parser.FieldString, *parser.FieldNumeric, *parser.FieldBool, *parser.FieldTime, *parser.FieldDuration, *parser.FieldUUID, *parser.FieldNode:
		default:
			return nil, fmt.Errorf("complex ID field %s: unsupported type %T (only string, numeric, bool, time.Time, time.Duration, UUID, and node references are allowed)", sf.Name(), parsed)
		}

		fields = append(fields, parser.ComplexIDField{
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

	return parser.NewFieldComplexID("ID", kind, structName, fields), nil
}

func ParseNode(v gotype.Type, outPkg string, pkgScope gotype.Type) (*parser.Node, error) {
	internalPkg := path.Join(outPkg, "internal")

	node := &parser.Node{Name: v.Name()}

	var features parser.FeatureSet

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Name() == "Node" && IsGenericNodeFromSom(f.Elem(), outPkg, "Node") {
				if IsKnownStringIDType(f.Elem()) {
					gen := ParseIDType(f.Elem())
					node.IDType = gen
					node.IDEmbed = f.Name()
					node.Fields = append(node.Fields,
						parser.NewFieldID("ID", gen),
					)
				} else {
					node.IDEmbed = f.Name()

					complexID, err := ParseComplexIDFields(f.Elem(), outPkg, pkgScope)
					if err != nil {
						return nil, fmt.Errorf("model %s: %w", v.Name(), err)
					}
					node.IDType = complexID.Kind
					node.ComplexID = complexID
					node.Fields = append(node.Fields, complexID)
				}
				continue
			}

			if parser.ParseFeature(f, internalPkg, &features, &node.Fields) {
				continue
			}

			return nil, fmt.Errorf("model %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("model %s: field ID not allowed, already provided by som.%s", v.Name(), node.IDEmbed)
		}

		field, err := parser.ParseField(f, outPkg)
		if err != nil {
			return nil, err
		}

		node.Fields = append(node.Fields, field)
	}

	parser.ApplyFeatures(features, &node.Timestamps, &node.OptimisticLock, &node.SoftDelete, &node.Fields)

	return node, nil
}
