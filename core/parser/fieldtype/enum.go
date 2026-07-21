package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/parser/structtype"
	"github.com/wzshiming/gotype"
)

type EnumHandler struct{}

func (h *EnumHandler) Match(elem gotype.Type, ctx *parser.FieldContext) bool {
	return elem.Kind() == gotype.String && structtype.IsEnum(elem, ctx.OutPkg)
}

func (h *EnumHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldEnum(t.Name(), elem.Name()), nil
}
