package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldSemVer is a som.SemVer.
type FieldSemVer struct {
	fieldBase
}

func (f *FieldSemVer) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.String && elem.PkgPath() == ctx.OutPkg && elem.Name() == "SemVer"
}

func (f *FieldSemVer) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldSemVer{fieldBase: newBase(t.Name())}, nil
}
