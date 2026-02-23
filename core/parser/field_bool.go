package parser

import "github.com/wzshiming/gotype"

type boolFieldHandler struct{}



func (h *boolFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Bool
}

func (h *boolFieldHandler) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldBool{&fieldAtomic{name: t.Name()}}, nil
}
