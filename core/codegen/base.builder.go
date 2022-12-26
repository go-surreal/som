package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/codegen/field"
	"github.com/marcbinz/som/core/parser"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	if err := b.buildSchemaFile(); err != nil {
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

type Config struct {
	Address string
	Username string
	Password string
	Namespace string
	Database string
}

type Client struct {
	db Database
}

func NewClient(conf Config) (*Client, error) {
	surreal, err := surrealdb.New(conf.Address + "/rpc")
	if err != nil {
		return nil, fmt.Errorf("new failed: %v", err)
	}

	_, err = surreal.Signin(map[string]any{
		"user": conf.Username,
		"pass": conf.Password,
	})
	if err != nil {
		return nil, err
	}

	_, err = surreal.Use(conf.Namespace, conf.Database)
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
	raw, err := db.DB.Query(statement, vars)
	if err != nil {
		return nil, err
	}

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

func (b *build) buildSchemaFile() error {
	statements := []string{"", ""}

	var fieldFn func(table string, f field.Field, prefix string)
	fieldFn = func(table string, f field.Field, prefix string) {
		fieldType := f.TypeDatabase()
		if fieldType == "" {
			return
		}

		statement := fmt.Sprintf(
			"DEFINE FIELD %s ON TABLE %s TYPE %s;",
			prefix+f.NameDatabase(), table, fieldType,
		)
		statements = append(statements, statement)

		if object, ok := f.(*field.Struct); ok {
			for _, fld := range object.Table().GetFields() {
				fieldFn(table, fld, prefix+f.NameDatabase()+".")
			}
		}

		if slice, ok := f.(*field.Slice); ok {
			statement := fmt.Sprintf(
				"DEFINE FIELD %s ON TABLE %s TYPE %s;",
				prefix+f.NameDatabase()+".*", table, slice.Element().TypeDatabase(),
			)
			statements = append(statements, statement)

			if object, ok := slice.Element().(*field.Struct); ok {
				for _, fld := range object.Table().GetFields() {
					fieldFn(table, fld, prefix+f.NameDatabase()+".*.")
				}
			}
		}
	}

	for _, node := range b.input.nodes {
		statement := fmt.Sprintf("DEFINE TABLE %s SCHEMAFULL;", node.NameDatabase())
		statements = append(statements, statement)

		for _, f := range node.GetFields() {
			fieldFn(node.NameDatabase(), f, "")
		}

		statements = append(statements, "")
	}

	for _, edge := range b.input.edges {
		statement := fmt.Sprintf("DEFINE TABLE %s SCHEMAFULL;", edge.NameDatabase())
		statements = append(statements, statement)

		for _, f := range edge.GetFields() {
			fieldFn(edge.NameDatabase(), f, "")
		}

		statements = append(statements, "")
	}

	content := strings.Join(statements, "\n")

	tmpl := `%s

package %s

import(
	"fmt"
)
	
func (c *Client) ApplySchema() error {
	_, err := c.db.Query(tmpl, nil)
	if err != nil {
		return fmt.Errorf("could not apply schema: %%v", err)
	}

	return nil
}

var tmpl = %s
`

	data := []byte(fmt.Sprintf(tmpl, codegenComment, b.basePkgName(), "`"+content+"`"))

	err := os.WriteFile(path.Join(b.basePath(), "schema.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *build) buildBaseFile(node *field.NodeTable) error {
	pkgQuery := b.subPkg(def.PkgQuery)
	pkgConv := b.subPkg(def.PkgConv)

	f := jen.NewFile(b.basePkgName())

	f.PackageComment(codegenComment)

	f.Func().
		Params(jen.Id("c").Op("*").Id("Client")).
		Id(node.NameGo()).Params().
		Op("*").Id(node.NameGoLower()).
		Block(
			jen.Return(
				jen.Op("&").Id(node.NameGoLower()).
					Values(jen.Id("client").Op(":").Id("c")),
			),
		)

	f.Type().Id(node.NameGoLower()).Struct(
		jen.Id("client").Op("*").Id("Client"),
	)

	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("Query").Params().
		Op("*").Qual(pkgQuery, node.NameGo()).
		Block(
			jen.Return(jen.Qual(pkgQuery, "New"+node.NameGo()).Call(jen.Id("n").Dot("client").Dot("db"))),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("Create").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			// TODO: maybe add an option to the generator: "strict" vs. "non-strict" mode (regarding custom IDs)?
			jen.Id("key").Op(":=").Lit(node.NameDatabase()),
			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Op("!=").Lit("")).
				Block(
					jen.Id("key").Op("+=").
						Lit(":").Op("+").Lit("⟨").Op("+").Id(node.NameGoLower()).Dot("ID").Op("+").Lit("⟩"),
				),
			jen.Id("data").Op(":=").Qual(pkgConv, "From"+node.NameGo()).Call(jen.Op("*").Id(node.NameGoLower())),
			jen.Id("raw").Op(",").Err().Op(":=").
				Id("n").Dot("client").Dot("db").Dot("Create").
				Call(jen.Id("key"), jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.If(
				jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").
					Id("raw").Op(".").Call(jen.Index().Any()),
				jen.Op("!").Id("ok"),
			).Block(
				jen.Id("raw").Op("=").Index().Any().Values(jen.Id("raw")).Comment("temporary fix"),
			),

			jen.Var().Id("convNode").Qual(b.subPkg(def.PkgConv), node.NameGo()),
			jen.Err().Op("=").Qual(def.PkgSurrealDB, "Unmarshal").
				Call(jen.Id("raw"), jen.Op("&").Id("convNode")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.Op("*").Id(node.NameGoLower()).Op("=").
				Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).
				Call(jen.Id("convNode")),

			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("Read").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
		).
		Params(jen.Op("*").Add(b.input.SourceQual(node.NameGo())), jen.Bool(), jen.Error()).
		Block(
			jen.List(jen.Id("raw"), jen.Err()).Op(":=").
				Id("n").Dot("client").Dot("db").Dot("Select").
				Call(jen.Lit(node.NameDatabase()+":⟨").Op("+").Id("id").Op("+").Lit("⟩")),

			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.If(jen.Qual("errors", "As").Call(jen.Err(), jen.Op("&").Qual(def.PkgSurrealDB, "PermissionError").Values())).
					Block(jen.Return(jen.Nil(), jen.False(), jen.Nil())),
				jen.Return(jen.Nil(), jen.False(), jen.Err()),
			),

			jen.If(
				jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").
					Id("raw").Op(".").Call(jen.Index().Any()),
				jen.Op("!").Id("ok"),
			).Block(
				jen.Id("raw").Op("=").Index().Any().Values(jen.Id("raw")).Comment("temporary fix"),
			),

			jen.Var().Id("convNode").Qual(b.subPkg(def.PkgConv), node.NameGo()),
			jen.Err().Op("=").Qual(def.PkgSurrealDB, "Unmarshal").Call(jen.Id("raw"), jen.Op("&").Id("convNode")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Nil(), jen.False(), jen.Err()),
			),

			jen.Id("node").Op(":=").Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).Call(jen.Id("convNode")),
			jen.Return(jen.Op("&").Id("node"), jen.True(), jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("Update").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot update "+node.NameGo()+" without existing record ID"))),
				),

			jen.Id("data").Op(":=").Qual(pkgConv, "From"+node.NameGo()).Call(jen.Op("*").Id(node.NameGoLower())),

			jen.Id("raw").Op(",").Err().Op(":=").
				Id("n").Dot("client").Dot("db").Dot("Update").
				Call(jen.Lit(node.NameDatabase()+":⟨").Op("+").Id(node.NameGoLower()).Dot("ID").Op("+").Lit("⟩"), jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.Var().Id("convNode").Qual(b.subPkg(def.PkgConv), node.NameGo()),
			jen.Err().Op("=").Qual(def.PkgSurrealDB, "Unmarshal").
				Call(jen.Index().Any().Values(jen.Id("raw")), jen.Op("&").Id("convNode")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),

			jen.Op("*").Id(node.NameGoLower()).Op("=").
				Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).Call(jen.Id("convNode")),

			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("Delete").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.List(jen.Id("_"), jen.Err()).Op(":=").
				Id("n").Dot("client").Dot("db").Dot("Delete").
				Call(jen.Lit(node.NameDatabase()+":⟨").Op("+").Id(node.NameGoLower()).Dot("ID").Op("+").Lit("⟩")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Err()),
			),
			jen.Return(jen.Nil()),
		)

	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("Relate").Params().
		Op("*").Qual(b.subPkg(def.PkgRelate), node.NameGo()).
		Block(
			jen.Return(jen.Qual(b.subPkg(def.PkgRelate), "New"+node.NameGo()).
				Call(jen.Id("n").Dot("client").Dot("db"))),
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
