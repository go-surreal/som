package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/def"
	"github.com/marcbinz/sdb/core/parser"
)

type Enum struct {
	source          *parser.FieldEnum
	dbNameConverter NameConverter
}

func (f *Enum) NameGo() string {
	return f.source.Name
}

func (f *Enum) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *Enum) FilterDefine(sourcePkg string) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Base").Types(jen.Qual(sourcePkg, f.source.Typ), jen.Id("T"))
}

func (f *Enum) FilterInit(sourcePkg string, elemName string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewBase").Types(jen.Qual(sourcePkg, f.source.Typ), jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Enum) FilterFunc(sourcePkg, elemName string) jen.Code {
	// Enum does not need a filter function.
	return nil
}

func (f *Enum) SortDefine(types jen.Code) jen.Code {
	return nil
}

func (f *Enum) SortInit(types jen.Code) jen.Code {
	return nil
}

func (f *Enum) SortFunc(sourcePkg, elemName string) jen.Code {
	// Enum does not need a sort function.
	return nil
}

func (f *Enum) ConvFrom() jen.Code {
	return nil // TODO
}

func (f *Enum) ConvTo(elem string) jen.Code {
	return nil // TODO
}

func (f *Enum) FieldDef() jen.Code {
	return jen.Id(f.source.Name).String(). // TODO: support other enum base types (atomic)
						Tag(map[string]string{"json": strcase.ToSnake(f.source.Name)})
}
