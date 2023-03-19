package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/codegen/field"
	"os"
	"path"
)

type relateBuilder struct {
	*baseBuilder
}

func newRelateBuilder(input *input, basePath, basePkg, pkgName string) *relateBuilder {
	return &relateBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *relateBuilder) build() error {
	if err := b.createDir(); err != nil {
		return err
	}

	if err := b.buildBaseFile(); err != nil {
		return err
	}

	for _, node := range b.nodes {
		if err := b.buildNodeFile(node); err != nil {
			return err
		}
	}

	for _, edge := range b.edges {
		if err := b.buildEdgeFile(edge); err != nil {
			return err
		}
	}

	return nil
}

func (b *relateBuilder) buildBaseFile() error {
	content := `

package relate

type Database interface {
	Query(statement string, vars any) (any, error)
}
`

	data := []byte(codegenComment + content)

	err := os.WriteFile(path.Join(b.path(), "relate.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *relateBuilder) buildNodeFile(node *field.NodeTable) error {
	file := jen.NewFile(b.pkgName)

	file.PackageComment(codegenComment)

	file.Line()
	file.Add(b.byNew(node))

	file.Line()
	file.Type().Id(node.Name).Struct(
		jen.Id("db").Id("Database"),
	)

	for _, fld := range node.GetFields() {
		slice, ok := fld.(*field.Slice)
		if !ok {
			continue
		}

		edgeElement, ok := slice.Element().(*field.Edge)
		if !ok {
			continue
		}

		file.Line()
		file.Func().Params(jen.Id("n").Id(node.NameGo())).
			Id(fld.NameGo()).Params().
			Id(edgeElement.Table().NameGoLower()).
			Block(
				jen.Return(jen.Id(edgeElement.Table().NameGoLower()).Call(jen.Id("n"))),
			)
	}

	if err := file.Save(path.Join(b.path(), node.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *relateBuilder) buildEdgeFile(edge *field.EdgeTable) error {
	file := jen.NewFile(b.pkgName)

	file.PackageComment(codegenComment)

	file.Line()
	file.Type().Id(edge.NameGoLower()).Struct(
		jen.Id("db").Id("Database"),
	)

	file.Line()
	file.Func().Params(jen.Id("e").Id(edge.NameGoLower())).
		Id("Create").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.Name))).
		Error().
		Block(
			jen.If(jen.Id("edge").Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the given edge must not be nil"))),
				),

			jen.If(jen.Id("edge").Dot("ID").Call().Op("!=").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID must not be set for an edge to be created"))),
				),

			jen.If(jen.Id("edge").Dot(edge.In.NameGo()).Dot("ID").Call().Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID of the incoming node '"+edge.In.NameGo()+"' must not be empty"))),
				),

			jen.If(jen.Id("edge").Dot(edge.Out.NameGo()).Dot("ID").Call().Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID of the outgoing node '"+edge.Out.NameGo()+"' must not be empty"))),
				),

			jen.Id("query").Op(":=").Lit("RELATE "),
			jen.Id("query").Op("+=").Lit(edge.In.NameDatabase()+":").Op("+").Id("edge").Dot(edge.In.NameGo()).Dot("ID").Call(),
			jen.Id("query").Op("+=").Lit("->"+edge.NameDatabase()+"->"),
			jen.Id("query").Op("+=").Lit(edge.Out.NameDatabase()+":").Op("+").Id("edge").Dot(edge.Out.NameGo()).Dot("ID").Call(),
			jen.Id("query").Op("+=").Lit(" CONTENT $data"),

			jen.Id("data").Op(":=").Qual(b.subPkg(def.PkgConv), "From"+edge.NameGo()).Call(jen.Op("*").Id("edge")),
			jen.Id("raw").Op(",").Err().Op(":=").Id("e").Dot("db").Dot("Query").
				Call(jen.Id("query"), jen.Map(jen.String()).Any().Values(jen.Lit("data").Op(":").Id("data"))),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.Var().Id("convEdge").Qual(b.subPkg(def.PkgConv), edge.NameGo()),
			jen.List(jen.Id("ok"), jen.Err()).Op(":=").Qual(def.PkgSurrealDB, "UnmarshalRaw").
				Call(jen.Id("raw"), jen.Op("&").Id("convEdge")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Return(jen.Qual("errors", "New").Call(jen.Lit("result is empty"))),
			),

			jen.Op("*").Id("edge").Op("=").Qual(b.subPkg(def.PkgConv), "To"+edge.NameGo()).Call(jen.Id("convEdge")),
			jen.Return(jen.Nil()),
		)

	file.Line()
	file.Func().Params(jen.Id(edge.NameGoLower())).
		Id("Update").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.NameGo()))).
		Error().
		Block(
			jen.Return(jen.Nil()),
		)

	file.Line()
	file.Func().Params(jen.Id(edge.NameGoLower())).
		Id("Delete").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.NameGo()))).
		Error().
		Block(
			jen.Return(jen.Nil()),
		)

	if err := file.Save(path.Join(b.path(), edge.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *relateBuilder) byNew(node field.Element) jen.Code {
	return jen.Func().Id("New" + node.NameGo()).
		Params(jen.Id("db").Id("Database")).
		Id("*").Id(node.NameGo()).
		Block(
			jen.Return(
				jen.Id("&").Id(node.NameGo()).Values(
					jen.Id("db").Op(":").Id("db"),
				),
			),
		)
}
