package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/dbtype"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
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
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Bool) FilterFunc(sourcePkg string, elem dbtype.Element) jen.Code {
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

func (f *Bool) ConvFrom(sourcePkg, elem string) jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *Bool) ConvTo(sourcePkg, elem string) jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *Bool) FieldDef() jen.Code {
	return jen.Id(f.source.Name).Bool().
		Tag(map[string]string{"json": strcase.ToSnake(f.source.Name) + ",omitempty"}) // TODO: store "false" (no omitempty)?
}
