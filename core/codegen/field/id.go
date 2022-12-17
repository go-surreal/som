package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/dbtype"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
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
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "ID").Types(jen.Id("T"))
}

func (f *ID) FilterInit(sourcePkg string, elemName string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewID").Types(jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))), jen.Lit(strcase.ToSnake(elemName)))
}

func (f *ID) FilterFunc(sourcePkg string, elem dbtype.Element) jen.Code {
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

func (f *ID) ConvFrom(sourcePkg, elem string) jen.Code {
	return jen.Id(funcBuildDatabaseID).Call(
		jen.Lit(strcase.ToSnake(elem)),
		jen.Id("data").Dot(f.source.Name),
	)
}

func (f *ID) ConvTo(sourcePkg, elem string) jen.Code {
	return jen.Id(funcParseDatabaseID).Call(
		jen.Lit(strcase.ToSnake(elem)),
		jen.Id("data").Dot(f.source.Name),
	)
}

func (f *ID) FieldDef() jen.Code {
	return jen.Id(f.source.Name).String().
		Tag(map[string]string{"json": strcase.ToSnake(f.source.Name) + ",omitempty"})
}
