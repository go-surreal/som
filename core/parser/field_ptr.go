package parser

import (
	"fmt"

	"github.com/wzshiming/gotype"
)

type ptrFieldHandler struct{}

func init() { RegisterFieldHandler(&ptrFieldHandler{}) }

func (h *ptrFieldHandler) Priority() int { return 50 }

func (h *ptrFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Ptr
}

func (h *ptrFieldHandler) Parse(t gotype.Type, elem gotype.Type, ctx *FieldContext) (Field, error) {
	field, err := defaultFieldRegistry.parse(elem, elem.Elem(), ctx)
	if err != nil {
		return nil, fmt.Errorf("could not parse elem for ptr field %s: %v", t.Name(), err)
	}

	if t.Name() != "" {
		field.setName(t.Name())
	}
	field.setPointer(true)

	return field, nil
}
