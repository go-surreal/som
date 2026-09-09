package fieldtype

import (
	"fmt"

	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/parser/structtype"
	"github.com/wzshiming/gotype"
)

// FragmentRefHandler rejects fields that reference a fragment type. A fragment
// is a projection of a node, not a table of its own, so it cannot be stored or
// linked to. Only nodes and edges may be referenced.
type FragmentRefHandler struct{}

func (h *FragmentRefHandler) Match(elem gotype.Type, ctx *parser.FieldContext) bool {
	return elem.Kind() == gotype.Struct && structtype.IsFragment(elem, ctx.OutPkg)
}

func (h *FragmentRefHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return nil, fmt.Errorf(
		"field %q references fragment %q: a fragment is a projection of a node and cannot be stored or linked (only nodes and edges may be referenced)",
		t.Name(), elem.Name(),
	)
}
