package structtype

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type ViewHandler struct{}

func (h *ViewHandler) Match(t gotype.Type, ctx *parser.TypeContext) bool {
	return IsView(t, ctx.OutPkg)
}

func (h *ViewHandler) Handle(t gotype.Type, ctx *parser.TypeContext) error {
	view, err := ParseView(t, ctx.OutPkg)
	if err != nil {
		return err
	}
	ctx.Output.Views = append(ctx.Output.Views, view)
	return nil
}

func (h *ViewHandler) Validate(ctx *parser.TypeContext) error {
	names := make([]string, len(ctx.Output.Views))
	for i, v := range ctx.Output.Views {
		names[i] = v.Name
		if err := validateFields("view "+v.Name, v.Fields, ctx.Output); err != nil {
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

	for i := 0; i < nf; i++ {
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

func ParseView(v gotype.Type, outPkg string) (*parser.View, error) {
	view := &parser.View{Name: v.Name()}

	nf := v.NumField()

	for i := 0; i < nf; i++ {
		f := v.Field(i)

		if !ast.IsExported(f.Name()) {
			continue
		}

		if f.IsAnonymous() {
			if f.Elem().PkgPath() == outPkg && f.Name() == "View" {
				// The view marker provides the read-only id column.
				view.Fields = append(view.Fields,
					parser.NewFieldID("ID", parser.IDTypeULID),
				)
				continue
			}

			return nil, fmt.Errorf("view %s: anonymous field %s not allowed", v.Name(), f.Name())
		}

		if strings.ToLower(f.Name()) == "id" {
			return nil, fmt.Errorf("view %s: field ID not allowed, already provided by som.View", v.Name())
		}

		field, err := parser.ParseField(f, outPkg)
		if err != nil {
			return nil, err
		}

		view.Fields = append(view.Fields, field)
	}

	return view, nil
}
