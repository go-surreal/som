package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/def"
	"github.com/marcbinz/sdb/core/parser"
)

type Bool struct {
	source          *parser.FieldBool
	dbNameConverter NameConverter
}

func (f *Bool) NameGo() string {
	return f.source.Name
}

func (f *Bool) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *Bool) FilterDefine(sourcePkg string) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Bool").Types(jen.Id("T"))
}

func (f *Bool) FilterInit(sourcePkg string, elemName string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewBool").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Bool) FilterFunc(sourcePkg, elemName string) jen.Code {
	// Bool does not need a filter function.
	return nil
}

func (f *Bool) SortDefine(types jen.Code) jen.Code {
	// TODO: should bool be sortable?
	return nil
}

func (f *Bool) SortInit(types jen.Code) jen.Code {
	// TODO: should bool be sortable?
	return nil
}

func (f *Bool) SortFunc(sourcePkg, elemName string) jen.Code {
	// Bool does not need a sort function.
	return nil
}

func (f *Bool) ConvFrom() jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *Bool) ConvTo(elem string) jen.Code {
	return jen.Id("data").Index(jen.Lit(strcase.ToSnake(f.source.Name))).Op(".").Parens(jen.Bool())
}
