package parser

import "github.com/wzshiming/gotype"

type byteFieldHandler struct{}

func init() { RegisterFieldHandler(&byteFieldHandler{}) }

func (h *byteFieldHandler) Priority() int { return 90 }

func (h *byteFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Byte
}

func (h *byteFieldHandler) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldByte{&fieldAtomic{name: t.Name()}}, nil
}
