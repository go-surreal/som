package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type DurationHandler struct{}

func (h *DurationHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Int64 && elem.PkgPath() == "time" && elem.Name() == "Duration"
}

func (h *DurationHandler) Parse(t gotype.Type, _ gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldDuration(t.Name()), nil
}
