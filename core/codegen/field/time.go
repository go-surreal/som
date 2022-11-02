package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/def"
	"github.com/marcbinz/sdb/core/parser"
)

type Time struct {
	source          *parser.FieldTime
	dbNameConverter NameConverter
}

func (f *Time) NameGo() string {
	return f.source.Name
}

func (f *Time) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *Time) FilterDefine(sourcePkg string) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Time").Types(jen.Id("T"))
}

func (f *Time) FilterInit(sourcePkg string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewTime").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Time) FilterFunc(sourcePkg, elemName string) jen.Code {
	// Time does not need a filter function.
	return nil
}

func (f *Time) SortDefine(types jen.Code) jen.Code {
	return jen.Id(f.source.Name).Op("*").Qual(def.PkgLibSort, "Sort").Types(jen.Id("T"))
}

func (f *Time) SortInit(types jen.Code) jen.Code {
	return jen.Qual(def.PkgLibSort, "NewSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Time) SortFunc(sourcePkg, elemName string) jen.Code {
	// Time does not need a sort function.
	return nil
}

func (f *Time) ConvFrom() jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *Time) ConvTo(elem string) jen.Code {
	return jen.Id("parseTime").Call(jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.source.Name))))
}
