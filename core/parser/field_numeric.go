package parser

import "github.com/wzshiming/gotype"

var numericKinds = map[gotype.Kind]NumberType{
	gotype.Int:     NumberInt,
	gotype.Int8:    NumberInt8,
	gotype.Int16:   NumberInt16,
	gotype.Int32:   NumberInt32,
	gotype.Int64:   NumberInt64,
	gotype.Uint8:   NumberUint8,
	gotype.Uint16:  NumberUint16,
	gotype.Uint32:  NumberUint32,
	gotype.Float32: NumberFloat32,
	gotype.Float64: NumberFloat64,
	gotype.Rune:    NumberRune,
}

type numericFieldHandler struct{}



func (h *numericFieldHandler) Match(elem gotype.Type, _ *FieldContext) bool {
	_, ok := numericKinds[elem.Kind()]
	return ok
}

func (h *numericFieldHandler) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldNumeric{&fieldAtomic{name: t.Name()}, numericKinds[elem.Kind()]}, nil
}
