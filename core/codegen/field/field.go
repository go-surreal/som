package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/core/parser"
)

const (
	funcPrepareID = "prepareID"
)

type NameConverter func(base string) string

type CodeDict struct {
	Key   jen.Code
	Value jen.Code
}

type Field interface {
	NameGo() string
	NameDatabase() string

	FilterDefine(sourcePkg string) jen.Code
	FilterInit(sourcePkg string, elemName string) jen.Code
	FilterFunc(sourcePkg, elemName string) jen.Code

	SortDefine(types jen.Code) jen.Code
	SortInit(types jen.Code) jen.Code
	SortFunc(sourcePkg, elemName string) jen.Code

	ConvFrom() jen.Code
	ConvTo(elem string) jen.Code
}

func Convert(field parser.Field, nameConv NameConverter) (Field, bool) {
	switch f := field.(type) {

	case *parser.FieldID:
		return &ID{
			source:          f,
			dbNameConverter: nameConv,
		}, true

	case *parser.FieldString:
		return &String{
			source:          f,
			dbNameConverter: nameConv,
		}, true

	case *parser.FieldNumeric:
		return &Numeric{
			source:          f,
			dbNameConverter: nameConv,
		}, true

	case *parser.FieldBool:
		return &Bool{
			source:          f,
			dbNameConverter: nameConv,
		}, true

	case *parser.FieldTime:
		return &Time{
			source:          f,
			dbNameConverter: nameConv,
		}, true

	case *parser.FieldUUID:
		return &UUID{
			source:          f,
			dbNameConverter: nameConv,
		}, true

	case *parser.FieldNode:
		return &Node{
			source:          f,
			dbNameConverter: nameConv,
		}, true

	case *parser.FieldEnum:
		return &Enum{
			source:          f,
			dbNameConverter: nameConv,
		}, true

	case *parser.FieldStruct:
		return &Struct{
			source:          f,
			dbNameConverter: nameConv,
		}, true

	case *parser.FieldSlice:
		return &Slice{
			source:          f,
			dbNameConverter: nameConv,
		}, true
	}

	return nil, false
}
