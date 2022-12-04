package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/dbtype"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/parser"
)

type Numeric struct {
	source          *parser.FieldNumeric
	dbNameConverter NameConverter
}

func (f *Numeric) NameGo() string {
	return f.source.Name
}

func (f *Numeric) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

// TODO: cleaner expose?
func (f *Numeric) TypeDatabase() string {
	switch f.source.Type {
	case parser.NumberInt:
		return "int"
	case parser.NumberInt32:
		return "int"
	case parser.NumberInt64:
		return "int"
	case parser.NumberFloat32:
		return "float"
	case parser.NumberFloat64:
		return "float"
	default:
		return "number" // TODO: okay?
	}
}

func (f *Numeric) FilterDefine(sourcePkg string) jen.Code {
	return jen.Id(f.NameGo()).Op("*").Qual(def.PkgLibFilter, "Numeric").Types(f.CodeNumberType(), jen.Id("T"))
}

func (f *Numeric) FilterInit(sourcePkg string, elemName string) jen.Code {
	return jen.Qual(def.PkgLibFilter, "NewNumeric").Types(f.CodeNumberType(), jen.Id("T")).
		Params(jen.Id("key").Dot("Dot").Call(jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Numeric) FilterFunc(sourcePkg string, elem dbtype.Element) jen.Code {
	// Numeric does not need a filter function.
	return nil
}

func (f *Numeric) SortDefine(types jen.Code) jen.Code {
	return jen.Id(f.source.Name).Op("*").Qual(def.PkgLibSort, "Sort").Types(jen.Id("T"))
}

func (f *Numeric) SortInit(types jen.Code) jen.Code {
	return jen.Qual(def.PkgLibSort, "NewSort").Types(jen.Id("T")).
		Params(jen.Id("keyed").Call(jen.Id("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))
}

func (f *Numeric) SortFunc(sourcePkg, elemName string) jen.Code {
	// Numeric does not need a sort function.
	return nil
}

func (f *Numeric) ConvFrom(sourcePkg, elem string) jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *Numeric) ConvTo(sourcePkg, elem string) jen.Code {
	return jen.Id("data").Dot(f.source.Name)
}

func (f *Numeric) FieldDef() jen.Code {
	return jen.Id(f.source.Name).Add(f.CodeNumberType()).
		Tag(map[string]string{"json": strcase.ToSnake(f.source.Name) + ",omitempty"})
}

func (f *Numeric) CodeNumberType() jen.Code {
	switch f.source.Type {
	case parser.NumberInt:
		return jen.Int()
	case parser.NumberInt32:
		return jen.Int32()
	case parser.NumberInt64:
		return jen.Int64()
	case parser.NumberFloat32:
		return jen.Float32()
	case parser.NumberFloat64:
		return jen.Float64()
	default:
		return jen.Int() // TODO: okay?
	}
}
