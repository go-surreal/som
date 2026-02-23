package parser

import "github.com/wzshiming/gotype"

type stringFieldHandler struct{}

func init() { RegisterFieldHandler(&stringFieldHandler{}) }

func (h *stringFieldHandler) Priority() int { return 90 }

func (h *stringFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.String
}

func (h *stringFieldHandler) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldString{&fieldAtomic{name: t.Name()}}, nil
}
