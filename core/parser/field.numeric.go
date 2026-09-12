package parser

import (
	"github.com/wzshiming/gotype"
)

type NumberType int32

const (
	NumberInt NumberType = iota
	NumberInt8
	NumberInt16
	NumberInt32
	NumberInt64
	//NumberUint
	NumberUint8
	NumberUint16
	NumberUint32
	//NumberUint64
	//NumberUintptr
	NumberFloat32
	NumberFloat64
	NumberRune
)

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

// FieldNumeric is any of the supported Go number types.
type FieldNumeric struct {
	fieldBase
	Type NumberType
}

func (f *FieldNumeric) Match(elem gotype.Type, _ *FieldContext) bool {
	_, ok := numericKinds[elem.Kind()]
	return ok
}

func (f *FieldNumeric) Parse(t gotype.Type, elem gotype.Type, _ *FieldContext) (Field, error) {
	return &FieldNumeric{fieldBase: newBase(t.Name()), Type: numericKinds[elem.Kind()]}, nil
}
