package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/parser/structtype"
	"github.com/wzshiming/gotype"
)

type EdgeRefHandler struct{}

func (h *EdgeRefHandler) Match(elem gotype.Type, ctx *parser.FieldContext) bool {
	return elem.Kind() == gotype.Struct && structtype.IsEdge(elem, ctx.OutPkg)
}

func (h *EdgeRefHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldEdge(t.Name(), elem.Name()), nil
}
