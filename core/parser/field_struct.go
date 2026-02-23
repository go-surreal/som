package parser

import "github.com/wzshiming/gotype"

type structFieldHandler struct{}

func init() { RegisterFieldHandler(&structFieldHandler{}) }

func (h *structFieldHandler) Priority() int { return 90 }

func (h *structFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Struct
}

func (h *structFieldHandler) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldStruct{&fieldAtomic{name: t.Name()}, elem.Name()}, nil
}
