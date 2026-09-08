package fieldtype

import (
	"fmt"

	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/parser/structtype"
	"github.com/wzshiming/gotype"
)

// SinkRefHandler rejects fields that reference a sink type. Sink records are
// discarded immediately after write, so they have no addressable id and
// cannot be linked to. Only nodes and edges may be referenced.
type SinkRefHandler struct{}

func (h *SinkRefHandler) Match(elem gotype.Type, ctx *parser.FieldContext) bool {
	return elem.Kind() == gotype.Struct && structtype.IsSink(elem, ctx.OutPkg)
}

func (h *SinkRefHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return nil, fmt.Errorf(
		"field %q references sink %q: sink records are discarded after write and cannot be linked (only nodes and edges may be referenced)",
		t.Name(), elem.Name(),
	)
}
