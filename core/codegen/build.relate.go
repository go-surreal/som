package codegen

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	if err := b.embedStaticFiles(); err != nil {
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

func (b *relateBuilder) embedStaticFiles() error {
	tmpl := &embed.Template{
		GenerateOutPath: b.subPkg(""),
	}

	files, err := embed.Relate(tmpl)
	if err != nil {
		return err
	}

	for _, file := range files {
		content := string(file.Content)
		content = strings.Replace(content, embedComment, codegenComment, 1)

		err := os.WriteFile(filepath.Join(b.path(), file.Path), []byte(content), os.ModePerm)
		if err != nil {
			return err
		}
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
		jen.Id("unmarshal").Func().Params(jen.Id("buf").Index().Byte(), jen.Id("val").Any()).Error(),
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
		jen.Id("unmarshal").Func().Params(jen.Id("buf").Index().Byte(), jen.Id("val").Any()).Error(),
	)

	file.Line()
	file.Func().Params(jen.Id("e").Id(edge.NameGoLower())).Id("Create").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("edge").Op("*").Add(b.SourceQual(edge.Name)),
		).
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

			jen.Id("query").Op(":=").Lit("RELATE ").Op("+").
				Lit(edge.In.NameDatabase()+":").Op("+").Id("edge").Dot(edge.In.NameGo()).Dot("ID").Call().Op("+").
				Lit("->"+edge.NameDatabase()+"->").Op("+").
				Lit(edge.Out.NameDatabase()+":").Op("+").Id("edge").Dot(edge.Out.NameGo()).Dot("ID").Call().Op("+").
				Lit(" CONTENT $data"),

			jen.Id("data").Op(":=").Qual(b.subPkg(def.PkgConv), "From"+edge.NameGo()).Call(jen.Id("edge")),

			jen.List(jen.Id("res"), jen.Err()).Op(":=").Id("e").Dot("db").Dot("Query").Call(
				jen.Id("ctx"),
				jen.Id("query"),
				jen.Map(jen.String()).Any().Values(jen.Lit("data").Op(":").Id("data")),
			),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not create relation: %w"), jen.Err())),
			),

			jen.Var().Id("convEdge").Op("*").Qual(b.subPkg(def.PkgConv), edge.NameGo()),
			jen.Err().Op("=").Id("e").Dot("unmarshal").Call(jen.Id("res"), jen.Op("&").Id("convEdge")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not unmarshal relation: %w"), jen.Err())),
			),

			jen.Op("*").Id("edge").Op("=").
				Op("*").Qual(b.subPkg(def.PkgConv), "To"+edge.NameGo()).Call(jen.Id("convEdge")),

			jen.Return(jen.Nil()),
		)

	file.Line()
	file.Func().Params(jen.Id(edge.NameGoLower())).
		Id("Update").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.NameGo()))).
		Error().
		Block(
			jen.Return(jen.Qual("errors", "New").Call(jen.Lit("not yet implemented"))),
		)

	file.Line()
	file.Func().Params(jen.Id(edge.NameGoLower())).
		Id("Delete").Params(jen.Id("edge").Op("*").Add(b.SourceQual(edge.NameGo()))).
		Error().
		Block(
			jen.Return(jen.Qual("errors", "New").Call(jen.Lit("not yet implemented"))),
		)

	if err := file.Save(path.Join(b.path(), edge.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *relateBuilder) byNew(node field.Element) jen.Code {
	return jen.Func().Id("New"+node.NameGo()).
		Params(
			jen.Id("db").Id("Database"),
			jen.Id("unmarshal").Func().Params(jen.Id("buf").Index().Byte(), jen.Id("val").Any()).Error(),
		).
		Id("*").Id(node.NameGo()).
		Block(
			jen.Return(
				jen.Id("&").Id(node.NameGo()).Values(
					jen.Id("db").Op(":").Id("db"),
					jen.Id("unmarshal").Op(":").Id("unmarshal"),
				),
			),
		)
}
