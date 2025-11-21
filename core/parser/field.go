package parser

import (
	"fmt"
)

type Field interface {
	fmt.Stringer
	field()
	Name() string
	Pointer() bool
	setName(string)
	setPointer(bool)
}

type fieldAtomic struct {
	name    string
	pointer bool
}

func (*fieldAtomic) field() {}

func (f *fieldAtomic) String() string {
	return f.Name()
}

func (f *fieldAtomic) Name() string {
	return f.name
}

func (f *fieldAtomic) setName(name string) {
	f.name = name
}

func (f *fieldAtomic) Pointer() bool {
	return f.pointer
}

func (f *fieldAtomic) setPointer(val bool) {
	f.pointer = val
}

type FieldID struct {
	*fieldAtomic
}

type FieldString struct {
	*fieldAtomic
}

type FieldNumeric struct {
	*fieldAtomic
	Type NumberType
}

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

type FieldBool struct {
	*fieldAtomic
}

type FieldByte struct {
	*fieldAtomic
}

type FieldDuration struct {
	*fieldAtomic
}

type FieldTime struct {
	*fieldAtomic
	IsCreatedAt bool
	IsUpdatedAt bool
}

type FieldUUID struct {
	*fieldAtomic
}

type FieldURL struct {
	*fieldAtomic
}

type FieldNode struct {
	*fieldAtomic
	Node string
}

type FieldEdge struct {
	*fieldAtomic
	Edge string
}

type FieldEnum struct {
	*fieldAtomic
	Typ      string
	BaseType string
}

type FieldStruct struct {
	*fieldAtomic
	Struct string
}

type FieldSlice struct {
	*fieldAtomic
	// Value  string
	Field Field
	// IsNode bool
	// IsEdge bool
	// IsEnum bool
}

// type FieldMap struct {
// 	fieldAtomic
// 	Key   string
// 	Value string
// }
