package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/dbtype"
	"github.com/marcbinz/sdb/core/parser"
)

type Struct struct {
	source          *parser.FieldStruct
	dbNameConverter NameConverter
}

func (f *Struct) NameGo() string {
	return f.source.Name
}

func (f *Struct) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *Struct) FilterDefine(sourcePkg string) jen.Code {
	// Struct uses a filter function instead.
	return nil
}

func (f *Struct) FilterInit(sourcePkg string, elemName string) jen.Code {
	// Struct uses a filter function instead.
	return nil
}

func (f *Struct) FilterFunc(sourcePkg string, elem dbtype.Element) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(strcase.ToLowerCamel(elem.NameGo())).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(strcase.ToLowerCamel(f.source.Struct)).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Struct).Types(jen.Id("T")).
				Params(jen.Id("n").Dot("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))))
}

func (f *Struct) SortDefine(types jen.Code) jen.Code {
	return nil // TODO
}

func (f *Struct) SortInit(types jen.Code) jen.Code {
	return nil // TODO
}

func (f *Struct) SortFunc(sourcePkg, elemName string) jen.Code {
	return nil // TODO
}

func (f *Struct) ConvFrom(sourcePkg, elem string) jen.Code {
	return jen.Op("*").Id("From" + f.source.Struct).Call(jen.Op("&").Id("data").Dot(f.source.Name))
}

func (f *Struct) ConvTo(sourcePkg, elem string) jen.Code {
	return jen.Op("*").Id("To" + f.source.Struct).Call(jen.Op("&").Id("data").Dot(f.source.Name))
}

func (f *Struct) FieldDef() jen.Code {
	return jen.Id(f.source.Name).Id(f.source.Struct).
		Tag(map[string]string{"json": strcase.ToSnake(f.source.Name) + ",omitempty"})
}
