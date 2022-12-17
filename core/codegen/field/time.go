package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/dbtype"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
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

func (f *Time) FilterInit(sourcePkg string, elemName string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewTime").Types(jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Time) FilterFunc(sourcePkg string, elem dbtype.Element) jen.Code {
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

func (f *Time) ConvFrom(sourcePkg, elem string) jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *Time) ConvTo(sourcePkg, elem string) jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *Time) FieldDef() jen.Code {
	return jen.Id(f.source.Name).Qual("time", "Time").
		Tag(map[string]string{"json": strcase.ToSnake(f.source.Name) + ",omitempty"})
}
