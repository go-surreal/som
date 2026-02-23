package parser

import "github.com/wzshiming/gotype"

type stringFieldHandler struct{}



func (h *stringFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.String
}

func (h *stringFieldHandler) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldString{&fieldAtomic{name: t.Name()}}, nil
}
