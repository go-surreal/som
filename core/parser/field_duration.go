package parser

import "github.com/wzshiming/gotype"

type durationFieldHandler struct{}



func (h *durationFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Int64 && elem.PkgPath() == "time" && elem.Name() == "Duration"
}

func (h *durationFieldHandler) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldDuration{&fieldAtomic{name: t.Name()}}, nil
}
