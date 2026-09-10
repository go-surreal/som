package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

type SemVerHandler struct{}

func (h *SemVerHandler) Match(elem gotype.Type, ctx *parser.FieldContext) bool {
	return elem.Kind() == gotype.String && elem.PkgPath() == ctx.OutPkg && elem.Name() == "SemVer"
}

func (h *SemVerHandler) Parse(t gotype.Type, _ gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldSemVer(t.Name()), nil
}
