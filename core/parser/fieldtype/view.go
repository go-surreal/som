package fieldtype

import (
	"fmt"

	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/parser/structtype"
	"github.com/wzshiming/gotype"
)

// ViewRefHandler rejects fields that reference a view type. Views are
// read-only, computed tables; their rows cannot be linked to, so a record
// link (or any field) pointing at a view is not allowed. Only nodes and
// edges may be referenced.
type ViewRefHandler struct{}

func (h *ViewRefHandler) Match(elem gotype.Type, ctx *parser.FieldContext) bool {
	return elem.Kind() == gotype.Struct && structtype.IsView(elem, ctx.OutPkg)
}

func (h *ViewRefHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return nil, fmt.Errorf(
		"field %q references view %q: views are read-only and cannot be linked (only nodes and edges may be referenced)",
		t.Name(), elem.Name(),
	)
}
