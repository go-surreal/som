package def

import "fmt"

type NumericType int

const (
	NumericUnknown NumericType = iota
	NumericInt
	NumericInt8
	NumericInt16
	NumericInt32
	NumericInt64
	NumericUint
	NumericUint8
	NumericUint16
	NumericUint32
	NumericUint64
	NumericUintptr
	NumericFloat32
	NumericFloat64
)

type Numeric struct {
	Base

	Type NumericType
}

func (n *Numeric) String() string {
	return fmt.Sprintf("[%s] %s (%d)", n.Package, n.Name, n.Type)
}
