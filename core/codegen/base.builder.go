package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/def"
	"github.com/marcbinz/som/core/codegen/field"
	"github.com/marcbinz/som/core/embed"
	"github.com/marcbinz/som/core/parser"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	filenameClient     = "som.client.go"
	filenameDatabase   = "som.database.go"
	filenameFunctions  = "som.functions.go"
	filenameInterfaces = "som.interfaces.go"
	filenameSchema     = "som.schema.go"
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

	if err := b.copyInternalPackage(); err != nil {
		return err
	}

	if err := b.buildClientFile(); err != nil {
		return err
	}

	if err := b.buildDatabaseFile(); err != nil {
		return err
	}

	if err := b.buildFunctionsFile(); err != nil {
		return err
	}

	if err := b.buildInterfaceFile(); err != nil {
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

func (b *build) copyInternalPackage() error {
	files, err := embed.Lib()
	if err != nil {
		return err
	}

	dir := filepath.Join(b.outDir, "internal")

	err = os.MkdirAll(filepath.Join(dir, "lib"), os.ModePerm)
	if err != nil {
		return err
	}

	for _, file := range files {
		content := string(file.Content)
		content = strings.Replace(content, "//go:build embed", codegenComment, 1)

		err := os.WriteFile(filepath.Join(dir, file.Path), []byte(content), os.ModePerm)
		if err != nil {
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
	"github.com/surrealdb/surrealdb.go/pkg/gorilla"
	"github.com/surrealdb/surrealdb.go/pkg/logger"
	"time"
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

type ClientImpl struct {
	db Database
}

func NewClient(conf Config) (*ClientImpl, error) {
	url := conf.Address + "/rpc"

	logData, err := logger.New().Make()
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %v", err)
	}

	ws, err := gorilla.Create().Logger(logData).SetTimeOut(time.Minute).Connect(url)
	if err != nil {
		return nil, fmt.Errorf("could not create websocket: %v", err)
	}
	
	surreal, err := surrealdb.New("<unused>", ws)
	if err != nil {
		return nil, fmt.Errorf("could not create surrealdb client: %v", err)
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

	return &ClientImpl{db: &database{DB: surreal}}, nil
}

func (c *ClientImpl) Close() {
	c.db.Close()
}
`

	data := []byte(codegenComment + "\n\npackage " + b.basePkgName() + content)

	err := os.WriteFile(path.Join(b.basePath(), filenameClient), data, os.ModePerm)
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

	err := os.WriteFile(path.Join(b.basePath(), filenameDatabase), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *build) buildFunctionsFile() error {
	content := `

import (
	"encoding/json"
	"fmt"
)

const statusOK = "OK"

type RawQuery[T any] struct {
	Status string
	Time   string
	Result T
	Detail string 
}

func Unmarshal[M any](respond interface{}) (model M, err error) {
	var bytes []byte

	if arrResp, isArr := respond.([]interface{}); len(arrResp) > 0 {
		if dataMap, ok := arrResp[0].(map[string]interface{}); ok && isArr {
			if _, ok := dataMap["status"]; ok {
				if bytes, err = json.Marshal(respond); err == nil {
					var raw []RawQuery[M]
					if err = json.Unmarshal(bytes, &raw); err == nil {
						if raw[0].Status != statusOK {
							err = fmt.Errorf("%s: %s", raw[0].Status, raw[0].Detail)
						}
						model = raw[0].Result
					}
				}
				return model, err
			}
		}
	}

	if bytes, err = json.Marshal(respond); err == nil {
		err = json.Unmarshal(bytes, &model)
	}

	return model, err
}
`

	data := []byte(codegenComment + "\n\npackage " + b.basePkgName() + content)

	err := os.WriteFile(path.Join(b.basePath(), filenameFunctions), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *build) buildInterfaceFile() error {
	f := jen.NewFile(b.basePkgName())

	f.PackageComment(codegenComment)

	f.Type().Id("Client").InterfaceFunc(func(g *jen.Group) {
		for _, node := range b.input.nodes {
			g.Id(node.NameGo() + "Repo").Call().Id(node.NameGo() + "Repo")
		}

		g.Id("Close").Call()
	})

	if err := f.Save(path.Join(b.basePath(), filenameInterfaces)); err != nil {
		return err
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
	
func (c *ClientImpl) ApplySchema() error {
	_, err := c.db.Query(tmpl, nil)
	if err != nil {
		return fmt.Errorf("could not apply schema: %%v", err)
	}

	return nil
}

var tmpl = %s
`

	data := []byte(fmt.Sprintf(tmpl, codegenComment, b.basePkgName(), "`"+content+"`"))

	err := os.WriteFile(path.Join(b.basePath(), filenameSchema), data, os.ModePerm)
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

	//
	// type {NodeName}Repo interface {...}
	//
	f.Type().Id(node.NameGo()+"Repo").Interface(
		jen.Id("Query").Call().Qual(pkgQuery, node.NameGo()),

		jen.Id("Create").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("user").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("CreateWithID").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
			jen.Id("user").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Read").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
		).Parens(jen.List(
			jen.Op("*").Add(b.input.SourceQual(node.NameGo())),
			jen.Bool(),
			jen.Error(),
		)),

		jen.Id("Update").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("user").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Delete").Call(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("user").Op("*").Add(b.input.SourceQual(node.NameGo())),
		).Error(),

		jen.Id("Relate").Call().Op("*").Qual(b.subPkg(def.PkgRelate), node.NameGo()),
	)

	f.Line()
	f.Func().
		Params(jen.Id("c").Op("*").Id("ClientImpl")).
		Id(node.NameGo() + "Repo").Params().Id(node.NameGo() + "Repo").
		Block(
			jen.Return(
				jen.Op("&").Id(node.NameGoLower()).
					Values(jen.Id("db").Op(":").Id("c").Dot("db")),
			),
		)

	f.Line()
	f.Type().Id(node.NameGoLower()).Struct(
		jen.Id("db").Id("Database"),
	)

	f.Line()
	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("Query").Params().
		Qual(pkgQuery, node.NameGo()).
		Block(
			jen.Return(jen.Qual(pkgQuery, "New"+node.NameGo()).Call(jen.Id("n").Dot("db"))),
		)

	onCreatedAt := jen.Empty()
	onUpdatedAt := jen.Empty()
	if node.HasTimestamps() {
		onCreatedAt = jen.Id("data").Dot("CreatedAt").Op("=").Qual("time", "Now").Call()
		onUpdatedAt = jen.Id("data").Dot("UpdatedAt").Op("=").Id("data").Dot("CreatedAt")
	}

	f.Line()
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

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("!=").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("given node already has an id"))),
				),

			jen.Id("key").Op(":=").Lit(node.NameDatabase()),
			jen.Id("data").Op(":=").Qual(pkgConv, "From"+node.NameGo()).Call(jen.Op("*").Id(node.NameGoLower())),

			jen.Add(onCreatedAt),
			jen.Add(onUpdatedAt),

			jen.Id("raw").Op(",").Err().Op(":=").
				Id("n").Dot("db").Dot("Create").
				Call(jen.Id("key"), jen.Id("data")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not create entity: %w"), jen.Err())),
			),

			jen.Var().Id("convNodes").Index().Qual(b.subPkg(def.PkgConv), node.NameGo()),
			jen.Err().Op("=").Qual(def.PkgSurrealDB, "Unmarshal").
				Call(jen.Id("raw"), jen.Op("&").Id("convNodes")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not unmarshal response: %w"), jen.Err())),
			),

			jen.If(jen.Len(jen.Id("convNodes")).Op("<").Lit(1)).Block(
				jen.Return(jen.Qual("errors", "New").Call(jen.Lit("response is empty"))),
			),

			jen.Op("*").Id(node.NameGoLower()).Op("=").
				Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).
				Call(jen.Id("convNodes").Index(jen.Lit(0))),

			jen.Return(jen.Nil()),
		)

	f.Line()
	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("CreateWithID").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(), // TODO: name clash if node/model is named "id"!
			jen.Id(node.NameGoLower()).Op("*").Add(b.input.SourceQual(node.NameGo())),
		).
		Error().
		Block(
			jen.If(jen.Id(node.NameGoLower()).Op("==").Nil()).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("the passed node must not be nil"))),
				),

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("!=").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(
						jen.Lit("creating node with preset ID not allowed, use CreateWithID for that")),
					),
				),

			jen.Id("key").Op(":=").Lit(node.NameDatabase()+":").Op("+").
				Lit("⟨").Op("+").Id("id").Op("+").Lit("⟩"),
			jen.Id("data").Op(":=").Qual(pkgConv, "From"+node.NameGo()).Call(jen.Op("*").Id(node.NameGoLower())),

			jen.Add(onCreatedAt),
			jen.Add(onUpdatedAt),

			jen.List(jen.Id("convNode"), jen.Err()).Op(":=").
				Qual(def.PkgSurrealDB, "SmartUnmarshal").Types(jen.Qual(b.subPkg(def.PkgConv), node.NameGo())).
				Call(
					jen.Id("n").Dot("db").Dot("Create").
						Call(jen.Id("key"), jen.Id("data")),
				),

			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not create entity: %w"), jen.Err())),
			),

			jen.Op("*").Id(node.NameGoLower()).Op("=").
				Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).
				Call(jen.Id("convNode")),

			jen.Return(jen.Nil()),
		)

	f.Line()
	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("Read").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("id").String(),
		).
		Params(jen.Op("*").Add(b.input.SourceQual(node.NameGo())), jen.Bool(), jen.Error()).
		Block(
			jen.List(jen.Id("convNode"), jen.Err()).Op(":=").
				Qual(def.PkgSurrealDB, "SmartUnmarshal").Types(jen.Qual(b.subPkg(def.PkgConv), node.NameGo())).
				Call(
					jen.Id("n").Dot("db").Dot("Select").
						Call(jen.Lit(node.NameDatabase()+":⟨").Op("+").Id("id").Op("+").Lit("⟩")),
				),

			jen.If(jen.Qual("errors", "Is").Call(jen.Err(), jen.Qual(def.PkgSurrealDB, "ErrNoRow"))).
				Block(jen.Return(jen.Nil(), jen.False(), jen.Nil())),

			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(
					jen.Nil(),
					jen.False(),
					jen.Qual("fmt", "Errorf").Call(jen.Lit("could not read entity: %w"), jen.Err()),
				),
			),

			jen.Id("node").Op(":=").Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).Call(jen.Id("convNode")),
			jen.Return(jen.Op("&").Id("node"), jen.True(), jen.Nil()),
		)

	if node.HasTimestamps() {
		onUpdatedAt = jen.Id("data").Dot("UpdatedAt").Op("=").Qual("time", "Now").Call()
	}

	f.Line()
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

			jen.If(jen.Id(node.NameGoLower()).Dot("ID").Call().Op("==").Lit("")).
				Block(
					jen.Return(jen.Qual("errors", "New").Call(jen.Lit("cannot update "+node.NameGo()+" without existing record ID"))),
				),

			jen.Id("data").Op(":=").Qual(pkgConv, "From"+node.NameGo()).Call(jen.Op("*").Id(node.NameGoLower())),

			jen.Add(onUpdatedAt),

			jen.List(jen.Id("convNode"), jen.Err()).Op(":=").
				Qual(def.PkgSurrealDB, "SmartUnmarshal").Types(jen.Qual(b.subPkg(def.PkgConv), node.NameGo())).
				Call(
					jen.Id("n").Dot("db").Dot("Update").
						Call(jen.Lit(node.NameDatabase()+":⟨").Op("+").Id(node.NameGoLower()).Dot("ID").Call().
							Op("+").Lit("⟩"), jen.Id("data")),
				),

			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not update entity: %w"), jen.Err())),
			),

			jen.Op("*").Id(node.NameGoLower()).Op("=").
				Qual(b.subPkg(def.PkgConv), "To"+node.NameGo()).Call(jen.Id("convNode")),

			jen.Return(jen.Nil()),
		)

	f.Line()
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
				Id("n").Dot("db").Dot("Delete").
				Call(jen.Lit(node.NameDatabase()+":⟨").Op("+").Id(node.NameGoLower()).Dot("ID").Call().Op("+").Lit("⟩")),
			jen.If(jen.Err().Op("!=").Nil()).Block(
				jen.Return(jen.Qual("fmt", "Errorf").Call(jen.Lit("could not delete entity: %w"), jen.Err())),
			),
			jen.Return(jen.Nil()),
		)

	f.Line()
	f.Func().
		Params(jen.Id("n").Op("*").Id(node.NameGoLower())).
		Id("Relate").Params().
		Op("*").Qual(b.subPkg(def.PkgRelate), node.NameGo()).
		Block(
			jen.Return(jen.Qual(b.subPkg(def.PkgRelate), "New"+node.NameGo()).
				Call(jen.Id("n").Dot("db"))),
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
