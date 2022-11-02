package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/def"
	"github.com/marcbinz/sdb/core/parser"
	"strings"
)

type Slice struct {
	source          *parser.FieldSlice
	dbNameConverter NameConverter
}

func (f *Slice) NameGo() string {
	return f.source.Name
}

func (f *Slice) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *Slice) FilterDefine(sourcePkg string) jen.Code {
	// Slice uses a filter function instead.
	return nil
}

func (f *Slice) FilterInit(sourcePkg string) jen.Code {
	// Slice uses a filter function instead.
	return nil
}

func (f *Slice) FilterFunc(sourcePkg, elemName string) jen.Code {
	if f.source.IsNode {
		return jen.Func().
			Params(jen.Id("n").Id(strings.ToLower(elemName)).Types(jen.Id("T"))).
			Id(f.NameGo()).Params().
			Id(strings.ToLower(f.source.Value)+"Slice").Types(jen.Id("T")).
			Block(
				jen.Id("key").Op(":=").Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(strcase.ToSnake(f.NameGo()))),
				jen.Return(
					jen.Id(strings.ToLower(f.source.Value)+"Slice").Types(jen.Id("T")).
						Values(
							jen.Id("new"+f.source.Value).Types(jen.Id("T")).
								Call(jen.Id("key")),
							jen.Qual(def.PkgLibFilter, "NewSlice").Types(jen.Qual(sourcePkg, f.source.Value), jen.Id("T")).
								Call(jen.Id("key")),
						),
				),
			)
	} else if f.source.IsEnum {
		return jen.Func().
			Params(jen.Id("n").Id(strings.ToLower(elemName)).Types(jen.Id("T"))).
			Id(f.NameGo()).Params().
			Op("*").Qual(def.PkgLibFilter, "Slice").Types(jen.Qual(sourcePkg, f.source.Value), jen.Id("T")).
			Block(
				jen.Return(
					jen.Qual(def.PkgLibFilter, "NewSlice").Types(jen.Qual(sourcePkg, f.source.Value), jen.Id("T")).
						Call(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(strcase.ToSnake(f.NameGo())))),
				),
			)
	} else {
		return jen.Func().
			Params(jen.Id("n").Id(strings.ToLower(elemName)).Types(jen.Id("T"))).
			Id(f.NameGo()).Params().
			Op("*").Qual(def.PkgLibFilter, "Slice").Types(jen.Id(f.source.Value), jen.Id("T")).
			Block(
				jen.Return(
					jen.Qual(def.PkgLibFilter, "NewSlice").Types(jen.Id(f.source.Value), jen.Id("T")).
						Call(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(strcase.ToSnake(f.NameGo())))),
				),
			)
	}
}

func (f *Slice) SortDefine(types jen.Code) jen.Code {
	return nil
}

func (f *Slice) SortInit(types jen.Code) jen.Code {
	return nil
}

func (f *Slice) SortFunc(sourcePkg, elemName string) jen.Code {
	return nil // TODO
}

func (f *Slice) ConvFrom() jen.Code {
	return nil
}

func (f *Slice) ConvTo(elem string) jen.Code {
	return nil
}
