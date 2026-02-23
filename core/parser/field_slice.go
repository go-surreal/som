package parser

import "github.com/wzshiming/gotype"

type sliceFieldHandler struct{}



func (h *sliceFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	return elem.Kind() == gotype.Slice
}

func (h *sliceFieldHandler) Parse(t gotype.Type, elem gotype.Type, ctx *FieldContext) (Field, error) {
	inner, err := defaultFieldRegistry.parse(elem, elem.Elem(), ctx)
	if err != nil {
		return nil, err
	}

	return &FieldSlice{
		fieldAtomic: &fieldAtomic{name: t.Name()},
		Field:       inner,
	}, nil
}
