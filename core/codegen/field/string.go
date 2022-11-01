package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/def"
	"github.com/marcbinz/sdb/core/parser"
)

type String struct {
	source          *parser.FieldString
	dbNameConverter NameConverter
}

func (f *String) NameGo() string {
	return f.source.Name
}

func (f *String) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *String) FilterDefine(sourcePkg string) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "String").Types(jen.Id("T"))
}

func (f *String) FilterInit(sourcePkg string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewString").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *String) FilterFunc(sourcePkg, elemName string) jen.Code {
	// String does not need a filter function.
	return nil
}

func (f *String) SortDefine(types jen.Code) jen.Code {
	return jen.Id(f.source.Name).Op("*").Qual(def.PkgLibSort, "String").Types(types)
}

func (f *String) SortInit(types jen.Code) jen.Code {
	return jen.Qual(def.PkgLibSort, "NewString").Types(types).Params(jen.Id("key"))
}

func (f *String) ConvFrom() jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *String) ConvTo(elem string) jen.Code {
	return jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.source.Name))).Op(".").Parens(jen.String())
}
