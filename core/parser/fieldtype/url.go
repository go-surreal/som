package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type URLHandler struct{}

func (h *URLHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	return elem.Kind() == gotype.Struct && elem.PkgPath() == "net/url" && elem.Name() == "URL"
}

func (h *URLHandler) Parse(t gotype.Type, _ gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldURL(t.Name()), nil
}
