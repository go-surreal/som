package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldEmail is a som.Email.
type FieldEmail struct {
	fieldBase
}

func (f *FieldEmail) Match(elem gotype.Type, ctx *FieldContext) bool {
	return elem.Kind() == gotype.String && elem.PkgPath() == ctx.OutPkg && elem.Name() == "Email"
}

func (f *FieldEmail) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldEmail{fieldBase: newBase(t.Name())}, nil
}
