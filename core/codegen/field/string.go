package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/dbtype"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
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

func (f *String) FilterInit(sourcePkg string, elemName string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewString").Types(jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *String) FilterFunc(sourcePkg string, elem dbtype.Element) jen.Code {
	// String does not need a filter function.
	return nil
}

func (f *String) SortDefine(types jen.Code) jen.Code {
	return jen.Id(f.source.Name).Op("*").Qual(def.PkgLibSort, "String").Types(jen.Id("T"))
}

func (f *String) SortInit(types jen.Code) jen.Code {
	return jen.Qual(def.PkgLibSort, "NewString").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *String) SortFunc(sourcePkg, elemName string) jen.Code {
	// String does not need a sort function.
	return nil
}

func (f *String) ConvFrom(sourcePkg, elem string) jen.Code {
	return jen.Id("data").Dot(f.source.Name) // TODO: vulnerability -> record link could be injected
}

func (f *String) ConvTo(sourcePkg, elem string) jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *String) FieldDef() jen.Code {
	return jen.Id(f.source.Name).String().
		Tag(map[string]string{"json": strcase.ToSnake(f.source.Name) + ",omitempty"})
}
