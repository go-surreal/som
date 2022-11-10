package codegen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/marcbinz/sdb/core/codegen/dbtype"
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

	for _, object := range b.objects {
		if err := b.buildFile(object); err != nil {
			return err
		}
	}

	return nil
}

func (b *convBuilder) buildBaseFile() error {
	content := `package conv

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

func prepareID(node string, id string) string {
	return strings.TrimPrefix(id, node+":")
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

// func extract[T any](val any, to func(map[string]any) T) T {
//	var t T
//	if val == nil {
//		return t
//	}
//	return to(val.(map[string]any))
// }
`

	err := os.WriteFile(path.Join(b.path(), "conv.go"), []byte(content), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write base file: %v", err)
	}

	return nil
}

func (b *convBuilder) buildFile(elem dbtype.Element) error {
	f := jen.NewFile(b.pkgName)

	f.Type().Id(elem.NameGo()).StructFunc(func(g *jen.Group) {
		for _, f := range elem.GetFields() {
			if code := f.FieldDef(); code != nil {
				g.Add(code)
			}
		}
	})

	f.Add(b.buildFrom(elem))
	f.Add(b.buildTo(elem))

	if node, ok := elem.(*dbtype.Node); ok {
		f.Add(b.buildFromRecord(node))
		f.Add(b.buildToRecord(node))
	}

	if err := f.Save(path.Join(b.path(), elem.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *convBuilder) buildFrom(elem dbtype.Element) jen.Code {
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
					if code := f.ConvFrom(); code != nil {
						d[jen.Id(f.NameGo())] = code
					}
				}
			}))),
		)
}

func (b *convBuilder) buildTo(elem dbtype.Element) jen.Code {
	return jen.Func().
		Id("To" + elem.NameGo()).
		Params(jen.Id("data").Op("*").Id(elem.NameGo())).
		Op("*").Add(b.SourceQual(elem.NameGo())).
		Block(
			jen.Return(jen.Id("&").Add(b.SourceQual(elem.NameGo())).Values(jen.DictFunc(func(d jen.Dict) {
				for _, f := range elem.GetFields() {
					if code := f.ConvTo(elem.NameGo()); code != nil {
						d[jen.Id(f.NameGo())] = code
					}
				}
			}))))
}

func (b *convBuilder) buildFromRecord(node *dbtype.Node) jen.Code {
	return jen.Func().
		Id("from"+node.NameGo()+"Record").
		Params(jen.Id("data").Any()).
		Op("*").Add(b.SourceQual(node.NameGo())).
		Block(
			jen.If(
				jen.Id("node").Op(",").Id("ok").Op(":=").Id("data").Op(".").Parens(jen.Op("*").Id(node.NameGo())),
				jen.Id("ok"),
			).Block(
				jen.Return(jen.Id("To"+node.NameGo()).Call(jen.Id("node"))),
			),
			jen.Return(jen.Op("&").Add(b.SourceQual(node.NameGo())).Values()),
		)
}

func (b *convBuilder) buildToRecord(node *dbtype.Node) jen.Code {
	return jen.Func().
		Id("to"+node.NameGo()+"Record").
		Params(jen.Id("node").Add(b.SourceQual(node.NameGo()))).
		String().
		Block(
			jen.If(jen.Id("node").Dot("ID").Op("==").Lit("")).
				Block(
					jen.Return(jen.Lit("")),
				),
			jen.Return(jen.Lit(node.NameDatabase()+":").Op("+").Id("node").Dot("ID")),
		)
}
