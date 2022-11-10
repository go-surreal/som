package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/def"
	"github.com/marcbinz/sdb/core/parser"
)

type UUID struct {
	source          *parser.FieldUUID
	dbNameConverter NameConverter
}

func (f *UUID) NameGo() string {
	return f.source.Name
}

func (f *UUID) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *UUID) FilterDefine(sourcePkg string) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Base").Types(jen.Qual(def.PkgUUID, "UUID"), jen.Id("T"))
}

func (f *UUID) FilterInit(sourcePkg string, elemName string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewBase").Types(jen.Qual(def.PkgUUID, "UUID"), jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *UUID) FilterFunc(sourcePkg, elemName string) jen.Code {
	// UUID does not need a filter function.
	return nil
}

func (f *UUID) SortDefine(types jen.Code) jen.Code {
	// UUID is not really sortable.
	return nil
}

func (f *UUID) SortInit(types jen.Code) jen.Code {
	// UUID is not really sortable.
	return nil
}

func (f *UUID) SortFunc(sourcePkg, elemName string) jen.Code {
	// UUID does not need a sort function.
	return nil
}

func (f *UUID) ConvFrom() jen.Code {
	return jen.Id("data").Dot(f.source.Name).Dot("String").Call()
}

func (f *UUID) ConvTo(elem string) jen.Code {
	return jen.Id("parseUUID").Call(jen.Id("data").Dot(f.source.Name))
}

func (f *UUID) FieldDef() jen.Code {
	return jen.Id(f.source.Name).String().
		Tag(map[string]string{"json": strcase.ToSnake(f.source.Name)})
}
