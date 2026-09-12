package parser

import (
	"fmt"
	"go/ast"
	"path"
	"slices"

	"github.com/wzshiming/gotype"
)

// Union is a set of node models sharing a common interface in the model
// package. It has no table of its own: a field typed as a union maps to a
// multi-table record link (record<a|b>), and an edge endpoint typed as a
// union maps to TYPE RELATION IN a|b.
//
// Like the other model types it is also its own TypeHandler.
type Union struct {
	Name    string
	Members []string

	// Open marks a union declared with som.OpenUnion instead of som.Union. It
	// requires no marker method, so a node can be a member of any number of
	// open unions, at the cost of the compiler no longer checking membership.
	Open bool
}

// Membership is a node's declaration that it belongs to a union. A sealed union
// is joined by embedding som.Part, which lets the compiler check the
// membership. An open one is joined by a blank som.Part field, which declares
// it without adding a method.
type Membership struct {
	Union  string
	Sealed bool
}

func (u *Union) Match(t gotype.Type, ctx *TypeContext) bool {
	return IsUnion(t, ctx.OutPkg)
}

func (u *Union) Handle(t gotype.Type, ctx *TypeContext) error {
	args := unionTypeArgs(t, ctx.OutPkg)
	open := isOpenUnion(t, ctx.OutPkg)

	if len(args) > 0 && open {
		return fmt.Errorf("union %s: embeds both som.Union and som.OpenUnion", t.Name())
	}

	if open {
		if unionMethods(t) < 1 {
			return fmt.Errorf(
				"union %s: an open union declares no marker method, so it must declare at least "+
					"one method of its own for a member to be checked against",
				t.Name(),
			)
		}

		ctx.Output.Unions = append(ctx.Output.Unions, &Union{Name: t.Name(), Open: true})

		return nil
	}

	if len(args) > 1 {
		return fmt.Errorf("union %s: embeds som.Union more than once", t.Name())
	}

	if args[0] != t.Name() {
		return fmt.Errorf(
			"union %s: the type parameter of som.Union must be the union interface itself, got %s",
			t.Name(), args[0],
		)
	}

	ctx.Output.Unions = append(ctx.Output.Unions, &Union{Name: t.Name()})

	return nil
}

func (u *Union) Validate(ctx *TypeContext) error {
	names := make([]string, len(ctx.Output.Unions))

	for i, union := range ctx.Output.Unions {
		names[i] = union.Name

		for _, node := range ctx.Output.Nodes {
			for _, member := range node.Unions {
				if member.Union != union.Name {
					continue
				}

				if member.Sealed != !union.Open {
					return fmt.Errorf(
						"node %s: %s is %s, so declare the membership as %s",
						node.Name, union.Name,
						map[bool]string{true: "an open union", false: "a sealed union"}[union.Open],
						map[bool]string{
							true:  fmt.Sprintf("the blank field `_ som.Part[%s]`", union.Name),
							false: fmt.Sprintf("the embed `som.Part[%s]`", union.Name),
						}[union.Open],
					)
				}

				union.Members = append(union.Members, node.Name)
			}
		}

		if len(union.Members) == 0 {
			return fmt.Errorf(
				"union %s: no model is part of it, add som.Part[%s] to at least one node",
				union.Name, union.Name,
			)
		}
	}

	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate union name %q", dup)
	}

	// A union shares the generated namespace with the tables, so its name must
	// not collide with any of them.
	for _, union := range ctx.Output.Unions {
		for kind, taken := range map[string]bool{
			"node":     nodeExists(union.Name, ctx.Output),
			"edge":     edgeExists(union.Name, ctx.Output),
			"enum":     enumExists(union.Name, ctx.Output),
			"struct":   structExists(union.Name, ctx.Output),
			"view":     viewExists(union.Name, ctx.Output),
			"sink":     sinkExists(union.Name, ctx.Output),
			"fragment": fragmentExists(union.Name, ctx.Output),
		} {
			if taken {
				return fmt.Errorf("union %s: name is already taken by a %s", union.Name, kind)
			}
		}
	}

	for _, node := range ctx.Output.Nodes {
		memberOf := make([]string, len(node.Unions))

		for i, member := range node.Unions {
			memberOf[i] = member.Union

			if !slices.Contains(names, member.Union) {
				return fmt.Errorf("node %s: som.Part references unknown union %q", node.Name, member.Union)
			}
		}

		if dup, ok := hasDuplicates(memberOf); ok {
			return fmt.Errorf("node %s: som.Part[%s] declared more than once", node.Name, dup)
		}
	}

	return nil
}

// IsUnion reports whether the given type is an interface declaring itself a
// union of node models by embedding som.Union or som.OpenUnion.
func IsUnion(t gotype.Type, outPkg string) bool {
	return len(unionTypeArgs(t, outPkg)) > 0 || isOpenUnion(t, outPkg)
}

// isOpenUnion reports whether the given interface embeds som.OpenUnion.
func isOpenUnion(t gotype.Type, outPkg string) bool {
	return interfaceEmbeds(t, func(expr ast.Expr) bool {
		sel, ok := expr.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "OpenUnion" {
			return false
		}

		ident, ok := sel.X.(*ast.Ident)

		return ok && ident.Name == path.Base(outPkg)
	})
}

// unionTypeArgs returns the names of the type arguments the given interface
// passes to the embedded som.Union, e.g. "Shape" for som.Union[Shape]. A union
// is free to declare methods of its own next to that embed, which every member
// then has to implement, so only the embeds are looked at.
func unionTypeArgs(t gotype.Type, outPkg string) []string {
	var args []string

	interfaceEmbeds(t, func(expr ast.Expr) bool {
		if arg, ok := genericTypeArg(expr, outPkg, "Union"); ok {
			args = append(args, arg)
		}

		return false
	})

	return args
}

// unionMethods counts the methods the given interface declares itself, which
// for an open union are the only thing a member is checked against.
func unionMethods(t gotype.Type) int {
	if t.Kind() != gotype.Interface {
		return 0
	}

	var count int

	nf := t.NumField()

	for i := range nf {
		elem, ok := t.Field(i).Origin().(*ast.Field)
		if ok && len(elem.Names) > 0 {
			count++
		}
	}

	return count
}

// interfaceEmbeds reports whether any element the given interface embeds
// matches the given predicate.
func interfaceEmbeds(t gotype.Type, match func(ast.Expr) bool) bool {
	if t.Kind() != gotype.Interface {
		return false
	}

	nf := t.NumField()

	for i := range nf {
		elem, ok := t.Field(i).Origin().(*ast.Field)
		if !ok || len(elem.Names) > 0 {
			continue
		}

		if match(elem.Type) {
			return true
		}
	}

	return false
}

// PartUnionName returns the name of the union the given som.Part embed makes
// its model a member of.
func PartUnionName(t gotype.Type, outPkg string) (string, error) {
	origin := t.Origin()
	if origin == nil {
		return "", fmt.Errorf("som.Part embed has no AST origin")
	}

	expr, ok := origin.(ast.Expr)
	if !ok {
		return "", fmt.Errorf("som.Part embed: unsupported AST node %T", origin)
	}

	arg, ok := genericTypeArg(expr, outPkg, "Part")
	if !ok {
		return "", fmt.Errorf("som.Part embed requires a union as type parameter, e.g. som.Part[Shape]")
	}

	return arg, nil
}

// BlankPartUnionName returns the union a blank som.Part field declares a
// membership in. A blank field of any other type is not a membership marker.
func BlankPartUnionName(t gotype.Type, outPkg string) (string, bool) {
	expr, ok := t.Origin().(ast.Expr)
	if !ok {
		return "", false
	}

	return genericTypeArg(expr, outPkg, "Part")
}

// genericTypeArg returns the type argument of an instantiation of the generic
// type som.<name>, e.g. "Shape" for the expression som.Union[Shape].
func genericTypeArg(expr ast.Expr, outPkg string, name string) (string, bool) {
	index, ok := expr.(*ast.IndexExpr)
	if !ok {
		return "", false
	}

	sel, ok := index.X.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != name {
		return "", false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok || ident.Name != path.Base(outPkg) {
		return "", false
	}

	arg, ok := index.Index.(*ast.Ident)
	if !ok {
		return "", false
	}

	return arg.Name, true
}
