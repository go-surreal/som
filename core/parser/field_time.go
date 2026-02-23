package parser

import "github.com/wzshiming/gotype"

type timeFieldHandler struct{}

func init() { RegisterFieldHandler(&timeFieldHandler{}) }

func (h *timeFieldHandler) Priority() int { return 30 }

func (h *timeFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Struct && elem.PkgPath() == "time" && elem.Name() == "Time"
}

func (h *timeFieldHandler) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldTime{&fieldAtomic{name: t.Name()}, false, false, false}, nil
}
