package parser

import (
	"github.com/wzshiming/gotype"
)

// FieldURL is a net/url.URL.
type FieldURL struct {
	fieldBase
}

func (f *FieldURL) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Struct && elem.PkgPath() == "net/url" && elem.Name() == "URL"
}

func (f *FieldURL) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldURL{fieldBase: newBase(t.Name())}, nil
}
