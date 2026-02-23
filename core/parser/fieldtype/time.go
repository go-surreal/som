package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type TimeHandler struct{}

func (h *TimeHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Struct && elem.PkgPath() == "time" && elem.Name() == "Time"
}

func (h *TimeHandler) Parse(t gotype.Type, _ gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldTime(t.Name()), nil
}
