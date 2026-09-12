package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/parser/structtype"
	"github.com/wzshiming/gotype"
)

type UnionRefHandler struct{}

func (h *UnionRefHandler) Match(elem gotype.Type, ctx *parser.FieldContext) bool {
	return structtype.IsUnion(elem, ctx.OutPkg)
}

func (h *UnionRefHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldUnion(t.Name(), elem.Name()), nil
}
