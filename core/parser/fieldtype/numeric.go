package fieldtype

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/wzshiming/gotype"
)

var numericKinds = map[gotype.Kind]parser.NumberType{
	gotype.Int:     parser.NumberInt,
	gotype.Int8:    parser.NumberInt8,
	gotype.Int16:   parser.NumberInt16,
	gotype.Int32:   parser.NumberInt32,
	gotype.Int64:   parser.NumberInt64,
	gotype.Uint8:   parser.NumberUint8,
	gotype.Uint16:  parser.NumberUint16,
	gotype.Uint32:  parser.NumberUint32,
	gotype.Float32: parser.NumberFloat32,
	gotype.Float64: parser.NumberFloat64,
	gotype.Rune:    parser.NumberRune,
}

type NumericHandler struct{}

func (h *NumericHandler) Match(elem gotype.Type, _ *parser.FieldContext) bool {
	_, ok := numericKinds[elem.Kind()]
	return ok
}

func (h *NumericHandler) Parse(t gotype.Type, elem gotype.Type, _ *parser.FieldContext) (parser.Field, error) {
	return parser.NewFieldNumeric(t.Name(), numericKinds[elem.Kind()]), nil
}
