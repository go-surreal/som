package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type SliceHandler struct{}

func (h *SliceHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Slice
}

func (h *SliceHandler) Parse(t gotype.Type, elem gotype.Type, ctx *parser.FieldContext) (parser.Field, error) {
	inner, err := ctx.RecurseParse(elem, elem.Elem(), ctx)
	if err != nil {
		return nil, err
	}

	return parser.NewFieldSlice(t.Name(), inner), nil
}
