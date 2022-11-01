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

func prepareID(node string, id any) string {
	return strings.TrimPrefix(id.(string), node+":")
}

func parseTime(val any) time.Time {
	res, err := time.Parse(time.RFC3339, val.(string))
	if err != nil {
		return time.Time{}
	}
	return res
}

func parseUUID(val any) uuid.UUID {
	res, err := uuid.Parse(val.(string))
	if err != nil {
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

	f.Add(b.buildFrom(elem))
	f.Add(b.buildTo(elem))

	if err := f.Save(path.Join(b.path(), elem.FileName())); err != nil {
		return err
	}

	return nil
}

func (b *convBuilder) buildFrom(elem dbtype.Element) jen.Code {
	return jen.Func().
		Id("From" + elem.NameGo()).
		Params(jen.Id("data").Add(b.SourceQual(elem.NameGo()))).
		Map(jen.String()).Any().
		Block(
			jen.Return(jen.Map(jen.String()).Any().Values(jen.DictFunc(func(d jen.Dict) {
				for _, f := range elem.GetFields() {
					if code := f.ConvFrom(); code != nil {
						d[jen.Lit(f.NameDatabase())] = code
					}
				}
			}))))
}

func (b *convBuilder) buildTo(elem dbtype.Element) jen.Code {
	return jen.Func().
		Id("To" + elem.NameGo()).
		Params(jen.Id("data").Map(jen.String()).Any()).
		Add(b.SourceQual(elem.NameGo())).
		Block(
			jen.Return(jen.Add(b.SourceQual(elem.NameGo())).Values(jen.DictFunc(func(d jen.Dict) {
				for _, f := range elem.GetFields() {
					if code := f.ConvTo(elem.NameGo()); code != nil {
						d[jen.Id(f.NameGo())] = code
					}
				}
			}))))
}

// from:
// to:
