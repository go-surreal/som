package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/som/core/codegen/field"
	"os"
	"path"
)

type convBuilder struct {
	*baseBuilder
}

func newConvBuilder(input *input, basePath, basePkg, pkgName string) *convBuilder {
	return &convBuilder{
		baseBuilder: newBaseBuilder(input, basePath, basePkg, pkgName),
	}
}

func (b *convBuilder) build() error {
	if err := b.createDir(); err != nil {
		return err
	}

	// Generate the base file.
	if err := b.buildBaseFile(); err != nil {
		return err
	}

	for _, node := range b.nodes {
		if err := b.buildFile(node); err != nil {
			return err
		}
	}

	for _, edge := range b.edges {
		if err := b.buildFile(edge); err != nil {
			return err
		}
	}

	for _, object := range b.objects {
		if err := b.buildFile(object); err != nil {
			return err
		}
	}

	return nil
}

func (b *convBuilder) buildBaseFile() error {
	content := `

package conv

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

func parseDatabaseID(node string, id string) string {
	return strings.TrimPrefix(id, node+":")
}

func buildDatabaseID(node string, id string) string {
	return node + ":" + id
}

func parseTime(val any) time.Time {
	res, err := time.Parse(time.RFC3339, val.(string))
	if err != nil {
		return time.Time{}
	}
	return res
}

func parseUUID(val string) uuid.UUID {
	res, err := uuid.Parse(val)
	if err != nil {
		// TODO: add logging!
		return uuid.UUID{}
	}
	return res
}
	
func mapRecords[I, O any](in []I, fn func(*I) O) []O {
	var out []O
	for _, i := range in {
		out = append(out, fn(&i))
	}
	return out
}
	
func convertEnum[I, O ~string](in []I) []O {
	var out []O
	for _, i := range in {
		out = append(out, O(i))
	}
	return out
}

// func extract[T any](val any, to func(map[string]any) T) T {
//	var t T
//	if val == nil {
//		return t
//	}
//	return to(val.(map[string]any))
// }
`

	data := []byte(codegenComment + content)

	err := os.WriteFile(path.Join(b.path(), "conv.go"), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *convBuilder) buildFile(elem field.Element) error {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		Elem:      elem,
	}

	f := jen.NewFile(b.pkgName)

	f.PackageComment(codegenComment)

	f.Type().Id(elem.NameGo()).StructFunc(func(g *jen.Group) {
		for _, f := range elem.GetFields() {
			if code := f.CodeGen().FieldDef(fieldCtx); code != nil {
				g.Add(code)
			}
		}
	})

	f.Add(b.buildFrom(elem))
	f.Add(b.buildTo(elem))

	if node, ok := elem.(*field.DatabaseNode); ok {
		f.Type().Id(node.NameGo() + "Field").Id(node.NameGo())

		f.Func().Params(jen.Id("f").Op("*").Id(node.NameGo()+"Field")).
			Id("MarshalJSON").Params().
			Params(jen.Index().Byte(), jen.Error()).
			Block(
				jen.If(jen.Id("f").Op("==").Nil()).Block(
					jen.Return(jen.Nil(), jen.Nil()),
				),
				jen.Return(jen.Qual("encoding/json", "Marshal").Call(jen.Id("f").Dot("ID"))),
			)

		f.Func().Params(jen.Id("f").Op("*").Id(node.NameGo()+"Field")).
			Id("UnmarshalJSON").Params(jen.Id("data").Index().Byte()).
			Error().
			Block(
				jen.Id("raw").Op(":=").String().Call(jen.Id("data")),
				jen.If(
					jen.Qual("strings", "HasPrefix").Call(jen.Id("raw"), jen.Lit("\"")).
						Op("&&").Qual("strings", "HasSuffix").Call(jen.Id("raw"), jen.Lit("\"")),
				).
					Block(
						jen.Id("raw").Op("=").Id("raw").Index(jen.Lit(1).Op(":").Len(jen.Id("raw")).Op("-").Lit(1)),
						jen.Id("f").Dot("ID").Op("=").Id("parseDatabaseID").Call(jen.Lit(node.NameDatabase()), jen.Id("raw")),
						jen.Return(jen.Nil()),
					),

				jen.Type().Id("fieldAlias").Id(node.NameGo()+"Field"),
				jen.Var().Id("field").Id("fieldAlias"),

				jen.Err().Op(":=").Qual("encoding/json", "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("field")),
				jen.If(jen.Err().Op("==").Nil()).Block(
					jen.Op("*").Id("f").Op("=").Id(node.NameGo()+"Field").Call(jen.Id("field")),
				),

				jen.Return(jen.Err()),
			)

		f.Add(b.buildFromField(node))
		f.Add(b.buildToField(node))
	}

	if err := f.Save(path.Join(b.path(), elem.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *convBuilder) buildFrom(elem field.Element) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		Elem:      elem,
	}

	return jen.Func().
		Id("From"+elem.NameGo()).
		Params(jen.Id("data").Op("*").Add(b.SourceQual(elem.NameGo()))).
		Op("*").Id(elem.NameGo()).
		Block(
			jen.If(jen.Id("data").Op("==").Nil()).Block(
				jen.Return(jen.Op("&").Id(elem.NameGo()).Values()),
			),
			jen.Return(jen.Op("&").Id(elem.NameGo()).Values(jen.DictFunc(func(d jen.Dict) {
				for _, f := range elem.GetFields() {
					if code := f.CodeGen().ConvFrom(fieldCtx); code != nil {
						d[jen.Id(f.NameGo())] = code
					}
				}
			}))),
		)
}

func (b *convBuilder) buildTo(elem field.Element) jen.Code {
	fieldCtx := field.Context{
		SourcePkg: b.sourcePkgPath,
		Elem:      elem,
	}

	return jen.Func().
		Id("To" + elem.NameGo()).
		Params(jen.Id("data").Op("*").Id(elem.NameGo())).
		Op("*").Add(b.SourceQual(elem.NameGo())).
		Block(
			jen.Return(jen.Id("&").Add(b.SourceQual(elem.NameGo())).Values(jen.DictFunc(func(d jen.Dict) {
				for _, f := range elem.GetFields() {
					if code := f.CodeGen().ConvTo(fieldCtx); code != nil {
						d[jen.Id(f.NameGo())] = code
					}
				}
			}))))
}

func (b *convBuilder) buildFromField(node *field.DatabaseNode) jen.Code {
	return jen.Func().
		Id("from"+node.NameGo()+"Field").
		Params(jen.Id("field").Id(node.NameGo()+"Field")).
		Op("*").Add(b.SourceQual(node.NameGo())).
		Block(
			jen.Id("node").Op(":=").Id(node.NameGo()).Call(jen.Id("field")),
			jen.Return(jen.Id("To"+node.NameGo()).Call(jen.Op("&").Id("node"))),
		)
}

func (b *convBuilder) buildToField(node *field.DatabaseNode) jen.Code {
	return jen.Func().
		Id("to" + node.NameGo() + "Field").
		Params(jen.Id("node").Op("*").Add(b.SourceQual(node.NameGo()))).
		Id(node.NameGo() + "Field").
		Block(
			jen.Return(jen.Id(node.NameGo() + "Field").Call(jen.Op("*").Id("From" + node.NameGo()).Call(jen.Id("node")))),
		)
}
