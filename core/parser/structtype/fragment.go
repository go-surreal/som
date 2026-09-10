package structtype

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

// fragmentEmbed is the name of the som.Fragment embed and therefore the name of
// the embedded field on every fragment model.
const fragmentEmbed = "Fragment"

// FragmentHandler parses fragment models, i.e. projections of a node model.
//
// A fragment can only be resolved once its parent node is known, which may be
// parsed after the fragment itself. Handle therefore only collects the declared
// types; the parent lookup, the field validation and the takeover of the parent
// fields all happen in Validate, which runs after all types are parsed.
type FragmentHandler struct {
	declared []declaredFragment
}

type declaredFragment struct {
	typ    gotype.Type
	name   string
	parent string
}

func (h *FragmentHandler) Match(t gotype.Type, ctx *parser.TypeContext) bool {
	return IsFragment(t, ctx.OutPkg)
}

func (h *FragmentHandler) Handle(t gotype.Type, ctx *parser.TypeContext) error {
	parent, err := fragmentParentName(t, ctx.OutPkg)
	if err != nil {
		return fmt.Errorf("fragment %s: %w", t.Name(), err)
	}

	h.declared = append(h.declared, declaredFragment{typ: t, name: t.Name(), parent: parent})

	return nil
}

func (h *FragmentHandler) Validate(ctx *parser.TypeContext) error {
	names := make([]string, 0, len(h.declared))

	for _, decl := range h.declared {
		fragment, err := resolveFragment(decl, ctx)
		if err != nil {
			return err
		}

		names = append(names, fragment.Name)
		ctx.Output.Fragments = append(ctx.Output.Fragments, fragment)
	}

	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate fragment name %q", dup)
	}

	return nil
}

func IsFragment(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)

		if !f.IsAnonymous() {
			continue
		}

		if f.Name() == fragmentEmbed && IsGenericNodeFromSom(f.Elem(), outPkg, fragmentEmbed) {
			return true
		}
	}

	return false
}

// fragmentParentName extracts the node model name from the som.Fragment type
// argument, e.g. "Person" for som.Fragment[Person].
func fragmentParentName(t gotype.Type, outPkg string) (string, error) {
	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)

		if !f.IsAnonymous() || f.Name() != fragmentEmbed ||
			!IsGenericNodeFromSom(f.Elem(), outPkg, fragmentEmbed) {
			continue
		}

		origin := f.Elem().Origin()
		if origin == nil {
			return "", fmt.Errorf("som.%s embed has no AST origin", fragmentEmbed)
		}

		indexExpr, ok := origin.(*ast.IndexExpr)
		if !ok {
			return "", fmt.Errorf("som.%s embed is not a generic instantiation (expected IndexExpr, got %T)", fragmentEmbed, origin)
		}

		switch idx := indexExpr.Index.(type) {
		case *ast.Ident:
			return idx.Name, nil
		case *ast.SelectorExpr:
			return idx.Sel.Name, nil
		default:
			return "", fmt.Errorf("som.%s type argument: unsupported AST node %T", fragmentEmbed, indexExpr.Index)
		}
	}

	return "", fmt.Errorf("no som.%s embed found", fragmentEmbed)
}

// resolveFragment validates the declared fragment against its parent node and
// takes over the matching parent fields.
func resolveFragment(decl declaredFragment, ctx *parser.TypeContext) (*parser.Fragment, error) {
	parentNode := findNode(decl.parent, ctx.Output)
	if parentNode == nil {
		return nil, fmt.Errorf("fragment %s: %q is not a node model", decl.name, decl.parent)
	}

	parentType, ok := ctx.PkgScope.ChildByName(decl.parent)
	if !ok {
		return nil, fmt.Errorf("fragment %s: node %s not found in package scope", decl.name, decl.parent)
	}

	fragment := &parser.Fragment{Name: decl.name, Parent: decl.parent}

	nf := decl.typ.NumField()

	for i := 0; i < nf; i++ {
		f := decl.typ.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Name() == fragmentEmbed {
				// The fragment marker (som.Fragment) projects no field of its
				// own. The record id is always projected and read into it.
				continue
			}

			return nil, fmt.Errorf("fragment %s: anonymous field %s not allowed", decl.name, f.Name())
		}

		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("fragment %s: field ID not allowed, the record id is always projected", decl.name)
		}

		if f.Tag().Get("som") != "" {
			return nil, fmt.Errorf(
				"fragment %s: field %s must not carry a som tag, it is taken over from node %s",
				decl.name, f.Name(), decl.parent,
			)
		}

		parentField, ok := findStructField(parentType, f.Name())
		if !ok {
			return nil, fmt.Errorf("fragment %s: field %s does not exist on node %s", decl.name, f.Name(), decl.parent)
		}

		if got, want := typeString(f), typeString(parentField); got != want {
			return nil, fmt.Errorf(
				"fragment %s: field %s has type %s, but %s on node %s",
				decl.name, f.Name(), got, want, decl.parent,
			)
		}

		field := findNodeField(parentNode, f.Name())
		if field == nil {
			return nil, fmt.Errorf("fragment %s: field %s is not projectable from node %s", decl.name, f.Name(), decl.parent)
		}

		fragment.Fields = append(fragment.Fields, field)
	}

	if len(fragment.Fields) == 0 {
		return nil, fmt.Errorf("fragment %s: no projected fields, at least one is required", decl.name)
	}

	return fragment, nil
}

func findNode(name string, output *parser.Output) *parser.Node {
	for _, n := range output.Nodes {
		if n.Name == name {
			return n
		}
	}
	return nil
}

func findNodeField(node *parser.Node, name string) parser.Field {
	for _, f := range node.Fields {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

func findStructField(t gotype.Type, name string) (gotype.Type, bool) {
	nf := t.NumField()

	for i := 0; i < nf; i++ {
		f := t.Field(i)
		if f.Name() == name {
			return f, true
		}
	}

	return nil, false
}

// typeString renders the declared type of a struct field, used to compare a
// fragment field against the node field it projects.
func typeString(f gotype.Type) string {
	return f.Elem().String()
}
