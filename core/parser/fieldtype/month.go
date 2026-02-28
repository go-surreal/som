package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type MonthHandler struct{}

func (h *MonthHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Int && elem.PkgPath() == "time" && elem.Name() == "Month"
}

func (h *MonthHandler) Parse(t gotype.Type, _ gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldMonth(t.Name()), nil
}
