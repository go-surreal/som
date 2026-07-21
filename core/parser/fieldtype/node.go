package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/parser/structtype"
	"github.com/wzshiming/gotype"
)

type NodeRefHandler struct{}

func (h *NodeRefHandler) Match(elem gotype.Type, ctx *parser.FieldContext) bool {
	return elem.Kind() == gotype.Struct && structtype.IsNode(elem, ctx.OutPkg)
}

func (h *NodeRefHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldNode(t.Name(), elem.Name()), nil
}
