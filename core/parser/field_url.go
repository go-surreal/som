package parser

import "github.com/wzshiming/gotype"

type urlFieldHandler struct{}

func init() { RegisterFieldHandler(&urlFieldHandler{}) }

func (h *urlFieldHandler) Priority() int { return 30 }

func (h *urlFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Struct && elem.PkgPath() == "net/url" && elem.Name() == "URL"
}

func (h *urlFieldHandler) Parse(t gotype.Type, _ gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldURL{&fieldAtomic{name: t.Name()}}, nil
}
