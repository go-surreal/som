package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type RegexHandler struct{}

func (h *RegexHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Struct && elem.PkgPath() == "regexp" && elem.Name() == "Regexp"
}

func (h *RegexHandler) Parse(t gotype.Type, _ gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldRegex(t.Name()), nil
}
