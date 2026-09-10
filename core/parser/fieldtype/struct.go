package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type StructHandler struct{}

func (h *StructHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Struct
}

func (h *StructHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldStruct(t.Name(), elem.Name()), nil
}
