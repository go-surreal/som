package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type WeekdayHandler struct{}

func (h *WeekdayHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Int && elem.PkgPath() == "time" && elem.Name() == "Weekday"
}

func (h *WeekdayHandler) Parse(t gotype.Type, _ gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldWeekday(t.Name()), nil
}
