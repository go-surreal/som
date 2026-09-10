package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type BoolHandler struct{}

func (h *BoolHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Bool
}

func (h *BoolHandler) Parse(t gotype.Type, _ gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldBool(t.Name()), nil
}
