package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/wzshiming/gotype"
)

// View is a read-only, pre-computed table view. It has a struct shape
// (its projected columns) like a Node, but no ID type, features or
// write operations. The SELECT that populates it is supplied separately
// via a //go:build som definition and linked back by database name.
//
// It is also its own TypeHandler.
type View struct {
	Name   string
	Fields []Field
}

func (v *View) Match(t gotype.Type, ctx *TypeContext) bool {
	return IsView(t, ctx.OutPkg)
}

func (v *View) Handle(t gotype.Type, ctx *TypeContext) error {
	view, err := ParseView(t, ctx)
	if err != nil {
		return err
	}
	ctx.Output.Views = append(ctx.Output.Views, view)
	return nil
}

func (v *View) Validate(ctx *TypeContext) error {
	names := make([]string, len(ctx.Output.Views))
	for i, view := range ctx.Output.Views {
		names[i] = view.Name
		if err := validateFields("view "+view.Name, view.Fields, ctx.Output); err != nil {
			return err
		}
	}
	if dup, ok := hasDuplicates(names); ok {
		return fmt.Errorf("duplicate view name %q", dup)
	}
	return nil
}

func IsView(t gotype.Type, outPkg string) bool {
	if t.Kind() != gotype.Struct {
		return false
	}

	nf := t.NumField()

	for i := range nf {
		f := t.Field(i)

		if !f.IsAnonymous() {
			continue
		}

		if f.Name() == "View" && f.Elem().Name() == "View" &&
			f.Elem().PkgPath() == outPkg {
			return true
		}
	}

	return false
}

func ParseView(v gotype.Type, ctx *TypeContext) (*View, error) {
	outPkg := ctx.OutPkg

	view := &View{Name: v.Name()}

	nf := v.NumField()

	for i := range nf {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() == outPkg && f.Name() == "View" {
				// The view marker (som.View) contributes no queryable column.
				// The record id is read into som.View during conversion, but
				// no id filter is generated: a view's id may be a composite
				// (GROUP BY) array that the string-based id filter cannot
				// express, so filtering a view by id is intentionally omitted.
				continue
			}

			return nil, fmt.Errorf("view %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("view %s: field ID not allowed, already provided by som.View", v.Name())
		}

		field, err := ctx.ParseField(f)
		if err != nil {
			return nil, fmt.Errorf("view %s: %w", v.Name(), err)
		}

		view.Fields = append(view.Fields, field)
	}

	return view, nil
}
