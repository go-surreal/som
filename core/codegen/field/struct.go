package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/parser"
	"strings"
)

type Struct struct {
	source          *parser.FieldStruct
	dbNameConverter NameConverter
}

func (f *Struct) NameGo() string {
	return f.source.Name
}

func (f *Struct) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *Struct) FilterDefine(sourcePkg string) jen.Code {
	// Struct uses a filter function instead.
	return nil
}

func (f *Struct) FilterInit(sourcePkg string) jen.Code {
	// Struct uses a filter function instead.
	return nil
}

func (f *Struct) FilterFunc(sourcePkg, elemName string) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(strings.ToLower(elemName)).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(strings.ToLower(f.source.Struct)).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Struct).Types(jen.Id("T")).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))))
}

func (f *Struct) SortDefine(types jen.Code) jen.Code {
	return nil
}

func (f *Struct) SortInit(types jen.Code) jen.Code {
	return nil
}

func (f *Struct) ConvFrom() jen.Code {
	return jen.Id("From" + f.source.Struct).Call(jen.Id("data").Dot(f.source.Name))
}

func (f *Struct) ConvTo(elem string) jen.Code {
	return jen.Id("To" + f.source.Struct).Call(jen.Id("data").
		Index(jen.Lit(strcase.ToSnake(f.source.Name))).Op(".").Parens(jen.Map(jen.String()).Any()))
}
