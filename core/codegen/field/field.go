package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/dbtype"
	"github.com/marcbinz/som/core/parser"
)

const (
	funcBuildDatabaseID = "buildDatabaseID"
	funcParseDatabaseID = "parseDatabaseID"
)

type Edge struct {
	Name   string
	In     Field
	Out    Field
	Fields []Field
}

type ElemGetter func(name string) (dbtype.Element, bool)

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
	FilterFunc(sourcePkg string, elem dbtype.Element) jen.Code

	SortDefine(types jen.Code) jen.Code
	SortInit(types jen.Code) jen.Code
	SortFunc(sourcePkg, elemName string) jen.Code

	ConvFrom(sourcePkg, elem string) jen.Code
	ConvTo(sourcePkg, elem string) jen.Code

	FieldDef() jen.Code
}

func Convert(field parser.Field, getElement ElemGetter, nameConv NameConverter) (Field, bool) {
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
			getElement:      getElement,
		}, true
	}

	return nil, false
}
