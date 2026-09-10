package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type StringHandler struct{}

func (h *StringHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.String
}

func (h *StringHandler) Parse(t gotype.Type, _ gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldString(t.Name()), nil
}
