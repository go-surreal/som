package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/codegen/field"
	"github.com/marcbinz/som/core/parser"
	"os"
	"path"
	"path/filepath"
)

type build struct {
	input  *input
	outDir string
	outPkg string
}

func Build(source *parser.Output, outDir string, outPkg string) error {
	in, err := newInput(source)
	if err != nil {
		return fmt.Errorf("error creating input: %v", err)
	}

	builder := &build{
		input:  in,
		outDir: outDir,
		outPkg: outPkg,
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

	if err := b.buildDatabaseFile(); err != nil {
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
		b.newFetchBuilder(),
		b.newConvBuilder(),
		b.newRelateBuilder(),
	}

	for _, builder := range builders {
		if err := builder.build(); err != nil {
			return err
		}
	}

	return nil
}

func (b *build) buildClientFile() error {
	content := `

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

type Database interface {
	Close()
	Create(thing string, data any) (any, error)
	Select(what string) (any, error)
	Query(statement string, vars any) (any, error)
	Update(thing string, data any) (any, error)
	Delete(what string) (any, error)
}

type Client struct {
	db Database
}

func NewClient(addr, user, pass, ns, db string) (*Client, error) {
	surreal, err := surrealdb.New(addr + "/rpc")
	if err != nil {
		return nil, fmt.Errorf("new failed: %v", err)
	}

	_, err = surreal.Signin(map[string]any{
		"user": user,
		"pass": pass,
	})
	if err != nil {
		return nil, err
	}

	_, err = surreal.Use(ns, db)
	if err != nil {
		return nil, err
	}

	return &Client{db: &database{DB: surreal}}, nil
}

func (c *Client) Close() {
	c.db.Close()
}
`

	data := []byte(codegenComment + "\n\npackage " + b.basePkgName() + content)

	err := os.WriteFile(path.Join(b.basePath(), "client.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *build) buildDatabaseFile() error {
	content := `

import (
	"fmt"
	"github.com/surrealdb/surrealdb.go"
)

type database struct {
	*surrealdb.DB
}

func (db *database) Create(thing string, data any) (any, error) {
	return db.DB.Create(thing, data)
}

func (db *database) Select(what string) (any, error) {
	return db.DB.Select(what)
}

func (db *database) Query(statement string, vars any) (any, error) {
	fmt.Println(statement)

	raw, err := db.DB.Query(statement, vars)
	if err != nil {
		return nil, err
	}

	fmt.Println(raw)

	return raw, err
}

func (db *database) Update(what string, data any) (any, error) {
	return db.DB.Update(what, data)
}

func (db *database) Delete(what string) (any, error) {
	return db.DB.Delete(what)
}
`

	data := []byte(codegenComment + "\n\npackage " + b.basePkgName() + content)

	err := os.WriteFile(path.Join(b.basePath(), "database.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *build) buildBaseFile(node *field.DatabaseNode) error {
	pkgQuery := b.subPkg(def.PkgQuery)
	pkgConv := b.subPkg(def.PkgConv)

	f := jen.NewFile(b.basePkgName())

	f.PackageComment(codegenComment)

	f.Func().
		Params(jen.Id("c").Op("*").Id("Client")).
		Id(node.NameGo()).Params().
		Op("*").Id(strcase.ToLowerCamel(node.NameGo())).
		Block(
			jen.Return(
				jen.Op("&").Id(strcase.ToLowerCamel(node.NameGo())).
					Values(jen.Id("client").Op(":").Id("c")),
			),
		)

	f.Type().Id(strcase.ToLowerCamel(node.NameGo())).Struct(
		jen.Id("client").Op("*").Id("Client"),
	)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strcase.ToLowerCamel(node.NameGo()))).
		Id("Query").Params().
		Op("*").Qual(pkgQuery, node.NameGo()).
		Block(
			jen.Return(jen.Qual(pkgQuery, "New"+node.NameGo()).Call(jen.Id("n").Dot("client").Dot("db"))),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strcase.ToLowerCamel(node.NameGo()))).
		Id("Create").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			// TODO: maybe add an option to the generator: "strict" vs. "non-strict" mode (regarding custom IDs)?
			jen.Id("key").Op(":=").Lit(node.NameDatabase()),
			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Op("!=").Lit("")).
				Block(
					jen.Id("key").Op("+=").
						Lit(":").Op("+").Lit("⟨").Op("+").Id(node.NameGoLower()).Dot("ID").Op("+").Lit("⟩"),
				),
			jen.Id("data").Op(":=").Qual(pkgConv, "From"+node.NameGo()).Call(jen.Id(strcase.ToLowerCamel(node.NameGo()))),
			jen.Id("raw").Op(",").Err().Op(":=").
				Id("n").Dot("client").Dot("db").Dot("Create").
				Call(jen.Id("key"), jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.Var().Id("convNode").Qual(b.subPkg(def.PkgConv), node.NameGo()),
			jen.Err().Op("=").Qual(def.PkgSurrealDB, "Unmarshal").
				Call(jen.Id("raw"), jen.Op("&").Id("convNode")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.Op("*").Id(strcase.ToLowerCamel(node.NameGo())).Op("=").Op("*").Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).Call(jen.Op("&").Id("convNode")),
			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strcase.ToLowerCamel(node.NameGo()))).
		Id("Read").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
		).
		Params(jen.Op("*").Add(b.input.SourceQual(node.NameGo())), jen.Bool(), jen.Error()).
		Block(
			jen.List(jen.Id("raw"), jen.Err()).Op(":=").
				Id("n").Dot("client").Dot("db").Dot("Select").
				Call(jen.Lit(strcase.ToSnake(node.Name)+":").Op("+").Id("id")),

			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.If(jen.Qual("errors", "As").Call(jen.Err(), jen.Op("&").Qual(def.PkgSurrealDB, "PermissionError").Values())).
					Block(jen.Return(jen.Nil(), jen.False(), jen.Nil())),
				jen.Return(jen.Nil(), jen.False(), jen.Err()),
			),

			jen.Var().Id("convNode").Op("*").Qual(b.subPkg(def.PkgConv), node.NameGo()),
			jen.Err().Op("=").Qual(def.PkgSurrealDB, "Unmarshal").Call(jen.Index().Any().Values(jen.Id("raw")), jen.Op("&").Id("convNode")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.False(), jen.Err()),
			),

			jen.Return(jen.Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).Call(jen.Id("convNode")), jen.True(), jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strcase.ToLowerCamel(node.Name))).
		Id("Update").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strcase.ToLowerCamel(node.Name)).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(strcase.ToLowerCamel(node.NameGo())).Dot("ID").Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot update "+node.NameGo()+" without existing record ID"))),
				),

			jen.Id("data").Op(":=").Qual(pkgConv, "From"+node.NameGo()).Call(jen.Id(strcase.ToLowerCamel(node.NameGo()))),

			jen.Id("raw").Op(",").Err().Op(":=").Id("n").Dot("client").Dot("db").Dot("Update").Call(jen.Lit(node.NameDatabase()+":").Op("+").Id(strcase.ToLowerCamel(node.Name)).Dot("ID"), jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.Var().Id("convNode").Qual(b.subPkg(def.PkgConv), node.NameGo()),
			jen.Err().Op("=").Qual(def.PkgSurrealDB, "Unmarshal").
				Call(jen.Index().Any().Values(jen.Id("raw")), jen.Op("&").Id("convNode")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.Op("*").Id(strcase.ToLowerCamel(node.NameGo())).Op("=").Op("*").Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).Call(jen.Op("&").Id("convNode")),
			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strcase.ToLowerCamel(node.Name))).
		Id("Delete").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(strcase.ToLowerCamel(node.Name)).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.List(jen.Id("_"), jen.Err()).Op(":=").Id("n").Dot("client").Dot("db").Dot("Delete").Call(jen.Lit(node.NameDatabase()+":").Op("+").Id(strcase.ToLowerCamel(node.Name)).Dot("ID")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(strcase.ToLowerCamel(node.NameGo()))).
		Id("Relate").Params().
		Op("*").Qual(b.subPkg(def.PkgRelate), node.NameGo()).
		Block(
			jen.Return(jen.Qual(b.subPkg(def.PkgRelate), "New"+node.NameGo()).Call(jen.Id("n").Dot("client").Dot("db"))),
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

func (b *build) newFetchBuilder() builder {
	return newFetchBuilder(b.input, b.basePath(), b.basePkg(), def.PkgFetch)
}

func (b *build) newConvBuilder() builder {
	return newConvBuilder(b.input, b.basePath(), b.basePkg(), def.PkgConv)
}

func (b *build) newRelateBuilder() builder {
	return newRelateBuilder(b.input, b.basePath(), b.basePkg(), def.PkgRelate)
}

func (b *build) basePath() string {
	return b.outDir
}

func (b *build) basePkg() string {
	return b.outPkg
}

func (b *build) basePkgName() string {
	_, name := filepath.Split(b.outPkg)
	return name
}

func (b *build) subPkg(pkg string) string {
	return path.Join(b.basePkg(), pkg)
}
