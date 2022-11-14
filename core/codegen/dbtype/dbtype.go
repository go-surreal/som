package dbtype

import (
	"github.com/dave/jennifer/jen"
)

type Element interface {
	FileName() string
	GetFields() []Field
	NameGo() string
	NameDatabase() string
}

type Field interface {
	NameGo() string
	NameDatabase() string

	FilterDefine(sourcePkg string) jen.Code
	FilterInit(sourcePkg string, elemName string) jen.Code
	FilterFunc(sourcePkg string, elem Element) jen.Code

	SortDefine(types jen.Code) jen.Code
	SortInit(types jen.Code) jen.Code
	SortFunc(sourcePkg, elemName string) jen.Code

	ConvFrom(sourcePkg, elem string) jen.Code
	ConvTo(sourcePkg, elem string) jen.Code

	FieldDef() jen.Code
}
