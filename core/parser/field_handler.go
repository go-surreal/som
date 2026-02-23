package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

type FieldContext struct {
	OutPkg string
}

type FieldHandler interface {
	Priority() int
	Match(elem gotype.Type, ctx *FieldContext) bool
	Parse(t gotype.Type, elem gotype.Type, ctx *FieldContext) (Field, error)
}

var defaultFieldRegistry = &fieldRegistry{}

type fieldRegistry struct {
	registry[FieldHandler]
}

func RegisterFieldHandler(h FieldHandler) {
	defaultFieldRegistry.register(h)
}

func (r *fieldRegistry) parse(t gotype.Type, elem gotype.Type, ctx *FieldContext) (Field, error) {
	for _, h := range r.all(func(h FieldHandler) int { return h.Priority() }) {
		if h.Match(elem, ctx) {
			return h.Parse(t, elem, ctx)
		}
	}
	return nil, fmt.Errorf("field %s has unsupported type %s", t.Name(), elem.Kind())
}
