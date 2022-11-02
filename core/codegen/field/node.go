package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/parser"
	"strings"
)

type Node struct {
	source          *parser.FieldNode
	dbNameConverter NameConverter
}

func (f *Node) NameGo() string {
	return f.source.Name
}

func (f *Node) NameDatabase() string {
	return f.dbNameConverter(f.source.Name)
}

func (f *Node) FilterDefine(sourcePkg string) jen.Code {
	// Node uses a filter function instead.
	return nil
}

func (f *Node) FilterInit(sourcePkg string) jen.Code {
	// Node uses a filter function instead.
	return nil
}

func (f *Node) FilterFunc(sourcePkg, elemName string) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(strings.ToLower(elemName)).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(strings.ToLower(f.source.Node)).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Node).Types(jen.Id("T")).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))))
}

func (f *Node) SortDefine(types jen.Code) jen.Code {
	// Node uses a sort function instead.
	return nil
}

func (f *Node) SortInit(types jen.Code) jen.Code {
	// Node uses a sort function instead.
	return nil
}

func (f *Node) SortFunc(sourcePkg, elemName string) jen.Code {
	return jen.Func().
		Params(jen.Id("n").Id(strings.ToLower(elemName)).Types(jen.Id("T"))).
		Id(f.NameGo()).Params().
		Id(strings.ToLower(f.source.Node)).Types(jen.Id("T")).
		Block(
			jen.Return(jen.Id("new" + f.source.Node).Types(jen.Id("T")).
				Params(jen.Id("keyed").Call(jen.Id("n").Dot("key"), jen.Lit(strcase.ToSnake(f.NameGo()))))))
}

func (f *Node) ConvFrom() jen.Code {
	return jen.Id("to" + f.source.Node + "Record").Call(jen.Id("data").Dot(f.source.Name))
}

func (f *Node) ConvTo(elem string) jen.Code {
	return jen.Id("from" + f.source.Node + "Record").Call(jen.Id("data").
		Index(jen.Lit(strcase.ToSnake(f.source.Name))))
}
