package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/sdb/core/codegen/dbtype"
	"github.com/marcbinz/sdb/core/codegen/def"
	"github.com/marcbinz/sdb/core/parser"
	"os"
	"path"
	"strings"
)

type build struct {
	input  *input
	outDir string
}

func Build(source *parser.Output, outDir string) error {
	in, err := newInput(source)
	if err != nil {
		return fmt.Errorf("error creating input: %v", err)
	}

	builder := &build{
		input:  in,
		outDir: outDir,
	}

	return builder.build()
}

func (b *build) build() error {
	if err := os.MkdirAll(b.basePath(), os.ModePerm); err != nil {
		return err
	}

	if err := b.buildClientFile(); err != nil {
		return err
	}

	for _, node := range b.input.nodes {
		if err := b.buildBaseFile(node); err != nil {
			return err
		}
	}

	builders := []builder{
		b.newQueryBuilder(),
		b.newFilterBuilder(),
		b.newSortBuilder(),
		b.newConvBuilder(),
	}

	for _, builder := range builders {
		if err := builder.build(); err != nil {
			return err
		}
	}

	return nil
}

func (b *build) buildClientFile() error {
	content := `package sdb

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

type Database interface {
	Close()
	Create(thing string, data map[string]any) (any, error)
	Select(what string) (any, error)
	Query(statement string, vars map[string]any) (any, error)
	Update(thing string, data map[string]any) (any, error)
	Delete(what string) (any, error)
}

type Client struct {
	db Database
}

func NewClient(addr, namespace, database, username, password string) (*Client, error) {
	db, err := surrealdb.New(addr + "/rpc")
	if err != nil {
		return nil, fmt.Errorf("new failed: %v", err)
	}

	_, err = db.Signin(map[string]any{
		"user": username,
		"pass": password,
	})
	if err != nil {
		return nil, err
	}

	_, err = db.Use(namespace, database)
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() {
	c.db.Close()
}
`

	err := os.WriteFile(path.Join(b.basePath(), "client.go"), []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *build) buildBaseFile(node *dbtype.Node) error {
	pkgQuery := b.subPkg(def.PkgQuery)
	pkgConv := b.subPkg(def.PkgConv)

	f := jen.NewFile(def.PkgBase)

	f.Func().
		Params(jen.Id("c").Op("*").Id("Client")).
		Id(node.NameGo()).Params().
		Op("*").Id(strings.ToLower(node.NameGo())).
		Block(
			jen.Return(
				jen.Op("&").Id(strings.ToLower(node.NameGo())).
					Values(jen.Id("client").Op(":").Id("c")),
			),
		)

	f.Type().Id(strings.ToLower(node.NameGo())).Struct(
		jen.Id("client").Op("*").Id("Client"),
	)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strings.ToLower(node.NameGo()))).
		Id("Query").Params().
		Op("*").Qual(pkgQuery, node.NameGo()).
		Block(
			jen.Return(jen.Qual(pkgQuery, "New"+node.NameGo()).Call(jen.Id("n").Dot("client").Dot("db"))),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strings.ToLower(node.NameGo()))).
		Id("Create").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower(node.NameGo())).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(strings.ToLower(node.NameGo())).Dot("ID").Op("!=").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("ID must not be set for a node to be created"))),
				),
			jen.Id("data").Op(":=").Qual(pkgConv, "From"+node.NameGo()).Call(jen.Op("*").Id(strings.ToLower(node.NameGo()))),
			jen.Id("raw").Op(",").Err().Op(":=").Id("n").Dot("client").Dot("db").Dot("Create").Call(jen.Lit(strcase.ToSnake(node.NameGo())), jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.Id("res").Op(":=").Qual(pkgConv, "To"+node.NameGo()).
				Call(jen.Id("raw").Op(".").Parens(jen.Index().Any()).Index(jen.Lit(0)).Op(".").Parens(jen.Map(jen.String()).Any())),
			jen.Op("*").Id(strings.ToLower(node.NameGo())).Op("=").Id("res"),
			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strings.ToLower(node.NameGo()))).
		Id("Read").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
		).
		Params(jen.Op("*").Add(b.input.SourceQual(node.NameGo())), jen.Error()).
		Block( // TODO: what if id already contains "user:" ?!
			jen.Id("raw").Op(",").Err().Op(":=").Id("n").Dot("client").Dot("db").Dot("Select").Call(jen.Lit(strcase.ToSnake(node.Name)).Op("+").Id("id")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.Err()),
			),
			jen.Id("res").Op(":=").Qual(pkgConv, "To"+node.Name).Call(jen.Id("raw").Op(".").Parens(jen.Map(jen.String()).Any())),
			jen.Return(jen.Op("&").Id("res"), jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strings.ToLower(node.Name))).
		Id("Update").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower(node.Name)).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strings.ToLower(node.Name))).
		Id("Delete").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strings.ToLower(node.Name)).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.Return(jen.Nil()),
		)

	if err := f.Save(path.Join(b.basePath(), node.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *build) newQueryBuilder() builder {
	return newQueryBuilder(b.input, b.basePath(), b.basePkg(), def.PkgQuery)
}

func (b *build) newFilterBuilder() builder {
	return newFilterBuilder(b.input, b.basePath(), b.basePkg(), def.PkgFilter)
}

func (b *build) newSortBuilder() builder {
	return newSortBuilder(b.input, b.basePath(), b.basePkg(), def.PkgSort)
}

func (b *build) newConvBuilder() builder {
	return newConvBuilder(b.input, b.basePath(), b.basePkg(), def.PkgConv)
}

func (b *build) basePath() string {
	return path.Join(b.outDir, def.PkgBase)
}

func (b *build) basePkg() string {
	return path.Join("github.com/marcbinz/sdb/example/gen", def.PkgBase) // TODO!!!
}

func (b *build) subPkg(pkg string) string {
	return path.Join(b.basePkg(), pkg)
}
