package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/def"
	"github.com/marcbinz/sdb/core/parser"
)

type ID struct {
	source          *parser.FieldID
	dbNameConverter NameConverter
}

func (f *ID) NameGo() string {
	return f.source.Name
}

func (f *ID) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *ID) FilterDefine(sourcePkg string) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Base").Types(jen.String(), jen.Id("T"))
}

func (f *ID) FilterInit(sourcePkg string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewBase").Types(jen.String(), jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *ID) FilterFunc(sourcePkg, elemName string) jen.Code {
	// ID does not need a filter function.
	return nil
}

func (f *ID) SortDefine(types jen.Code) jen.Code {
	return jen.Id(f.source.Name).Op("*").Qual(def.PkgLibSort, "Sort").Types(jen.Id("T"))
}

func (f *ID) SortInit(types jen.Code) jen.Code {
	return jen.Qual(def.PkgLibSort, "NewSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *ID) SortFunc(sourcePkg, elemName string) jen.Code {
	// ID does not need a sort function.
	return nil
}

func (f *ID) ConvFrom() jen.Code {
	// ID is never set, because the database is providing them, not the application.
	return nil
}

func (f *ID) ConvTo(elem string) jen.Code {
	return jen.Id(funcPrepareID).Call(
		jen.Lit(strcase.ToSnake(elem)),
		jen.Id("data").Index(jen.Lit(f.NameDatabase())),
	)
}
