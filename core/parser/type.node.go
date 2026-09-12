package parser

import (
	"fmt"
	"go/ast"
	"path"
	"strings"

	"github.com/wzshiming/gotype"
)

// Node is a model backed by its own table. It is also its own TypeHandler:
// the instance passed to Parse acts as the prototype that Match is called on.
type Node struct {
	Name           string
	Fields         []Field
	IDType         IDType
	IDEmbed        string
	ComplexID      *FieldComplexID
	Timestamps     bool
	OptimisticLock bool
	Changefeed     string
	SoftDelete     bool
	Expiry         bool
	ExpiryDuration string
}

func (n *Node) Match(t gotype.Type, ctx *TypeContext) bool {
	return IsNode(t, ctx.OutPkg)
}

func (n *Node) Handle(t gotype.Type, ctx *TypeContext) error {
	node, err := ParseNode(t, ctx)
	if err != nil {
		return err
	}
	ctx.Output.Nodes = append(ctx.Output.Nodes, node)
	return nil
}

func (n *Node) Validate(ctx *TypeContext) error {
	names := make([]string, len(ctx.Output.Nodes))
	for i, node := range ctx.Output.Nodes {
		names[i] = node.Name
		if err := validateFields("node "+node.Name, node.Fields, ctx.Output); err != nil {
			return err
		}
	}
	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate node name %q", dup)
	}
	return nil
}

func IsNode(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := range nf {
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
		return true
	}
	indexExpr, ok := origin.(*ast.IndexExpr)
	if !ok {
		return true
	}
	_, ok = indexExpr.Index.(*ast.SelectorExpr)
	return ok
}

func ParseIDType(t gotype.Type) IDType {
	origin := t.Origin()
	if origin == nil {
		return IDTypeULID
	}

	if indexExpr, ok := origin.(*ast.IndexExpr); ok {
		if selExpr, ok := indexExpr.Index.(*ast.SelectorExpr); ok {
			switch selExpr.Sel.Name {
			case "UUID":
				return IDTypeUUID
			case "Rand":
				return IDTypeRand
			case "String":
				return IDTypeString
			case "ULID":
				return IDTypeULID
			}
		}
	}

	return IDTypeULID
}

func ParseNode(v gotype.Type, ctx *TypeContext) (*Node, error) {
	outPkg := ctx.OutPkg
	internalPkg := path.Join(outPkg, "internal")

	node := &Node{Name: v.Name()}

	var features FeatureSet

	nf := v.NumField()

	for i := range nf {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Name() == "Node" && IsGenericNodeFromSom(f.Elem(), outPkg, "Node") {
				node.Changefeed = ParseChangefeedTag(f.Tag().Get("som"))

				if IsKnownStringIDType(f.Elem()) {
					gen := ParseIDType(f.Elem())
					node.IDType = gen
					node.IDEmbed = f.Name()
					node.Fields = append(node.Fields,
						NewFieldID("ID", gen),
					)
				} else {
					node.IDEmbed = f.Name()

					complexID, err := ParseComplexIDFields(f.Elem(), ctx)
					if err != nil {
						return nil, fmt.Errorf("model %s: %w", v.Name(), err)
					}
					node.IDType = complexID.Kind
					node.ComplexID = complexID
					node.Fields = append(node.Fields, complexID)
				}
				continue
			}

			matched, err := ParseFeature(f, internalPkg, &features, &node.Fields)
			if err != nil {
				return nil, fmt.Errorf("model %s: %w", v.Name(), err)
			}
			if matched {
				continue
			}

			return nil, fmt.Errorf("model %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("model %s: field ID not allowed, already provided by som.%s", v.Name(), node.IDEmbed)
		}

		field, err := ctx.ParseField(f)
		if err != nil {
			return nil, err
		}

		node.Fields = append(node.Fields, field)
	}

	ApplyFeatures(features, &node.Timestamps, &node.OptimisticLock, &node.SoftDelete, &node.Fields)

	node.Expiry = features.Expiry
	node.ExpiryDuration = features.ExpiryDuration

	return node, nil
}
