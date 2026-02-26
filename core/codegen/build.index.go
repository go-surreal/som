package codegen

import (
	"fmt"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/embed"
	"github.com/go-surreal/som/core/parser"
	"github.com/go-surreal/som/core/util/fs"
	"github.com/iancoleman/strcase"
)

type indexBuilder struct {
	*baseBuilder
	noCountIndex bool
	define       *parser.DefineOutput
}

func newIndexBuilder(input *input, fs *fs.FS, basePkg, pkgName string, noCountIndex bool, define *parser.DefineOutput) *indexBuilder {
	return &indexBuilder{
		baseBuilder:  newBaseBuilder(input, fs, basePkg, pkgName),
		noCountIndex: noCountIndex,
		define:       define,
	}
}

type indexEntry struct {
	IndexName string
	GoName    string
	TableName string
}

func (b *indexBuilder) build() error {
	for _, node := range b.nodes {
		entries := b.collectIndexEntries(node.NameDatabase(), node.GetFields(), node.Source.SoftDelete)
		if len(entries) == 0 {
			continue
		}
		if err := b.buildFile(node.NameGo(), node.FileName(), entries); err != nil {
			return err
		}
	}

	return nil
}

func (b *indexBuilder) buildFile(nameGo, fileName string, entries []indexEntry) error {
	f := jen.NewFile(b.pkgName)

	f.PackageComment(string(embed.CodegenComment))

	f.Line()
	f.Func().Id("New" + nameGo).
		Params(jen.Id("db").Id("Database")).
		Op("*").Id(nameGo).
		Block(
			jen.Return(jen.Op("&").Id(nameGo).Values(
				jen.Id("db").Op(":").Id("db"),
			)),
		)

	f.Line()
	f.Type().Id(nameGo).Struct(
		jen.Id("db").Id("Database"),
	)

	for _, entry := range entries {
		f.Line()
		f.Func().Params(jen.Id("i").Op("*").Id(nameGo)).
			Id(entry.GoName).Params().
			Op("*").Id("rebuildable").
			Block(
				jen.Return(jen.Op("&").Id("rebuildable").Values(
					jen.Id("db").Op(":").Id("i").Dot("db"),
					jen.Id("table").Op(":").Lit(entry.TableName),
					jen.Id("name").Op(":").Lit(entry.IndexName),
				)),
			)
	}

	if err := f.Render(b.fs.Writer(path.Join(b.path(), fileName))); err != nil {
		return err
	}

	return nil
}

func (b *indexBuilder) collectIndexEntries(tableName string, fields []field.Field, softDelete bool) []indexEntry {
	var entries []indexEntry

	if !b.noCountIndex {
		indexName := fmt.Sprintf(def.IndexPrefix+"%s_count", tableName)
		entries = append(entries, indexEntry{
			IndexName: indexName,
			GoName:    "Count",
			TableName: tableName,
		})
	}

	compositeUnique := make(map[string]bool)

	b.collectFieldIndexEntries(tableName, "", fields, &entries, compositeUnique)

	if softDelete {
		indexName := fmt.Sprintf(def.IndexPrefix+"%s_index_deleted_at", tableName)
		entries = append(entries, indexEntry{
			IndexName: indexName,
			GoName:    "IndexDeletedAt",
			TableName: tableName,
		})
	}

	for uniqueName := range compositeUnique {
		indexName := fmt.Sprintf(def.IndexPrefix+"%s_unique_%s", tableName, uniqueName)
		entries = append(entries, indexEntry{
			IndexName: indexName,
			GoName:    "Unique" + strcase.ToCamel(uniqueName),
			TableName: tableName,
		})
	}

	return entries
}

func (b *indexBuilder) collectFieldIndexEntries(tableName, fieldPrefix string, fields []field.Field, entries *[]indexEntry, compositeUnique map[string]bool) {
	for _, f := range fields {
		fieldPath := f.NameDatabase()
		if fieldPrefix != "" {
			fieldPath = fieldPrefix + "." + fieldPath
		}

		indexInfo := f.IndexInfo()
		searchInfo := f.SearchInfo()

		if indexInfo != nil {
			if indexInfo.Unique && indexInfo.UniqueName != "" {
				compositeUnique[indexInfo.UniqueName] = true
			} else if indexInfo.Unique {
				indexName := indexInfo.Name
				if indexName == "" {
					indexName = fmt.Sprintf(def.IndexPrefix+"%s_unique_%s", tableName, strings.ReplaceAll(fieldPath, ".", "_"))
				}
				goName := indexGoName(indexName, tableName)
				*entries = append(*entries, indexEntry{
					IndexName: indexName,
					GoName:    goName,
					TableName: tableName,
				})
			} else {
				indexName := indexInfo.Name
				if indexName == "" {
					indexName = fmt.Sprintf(def.IndexPrefix+"%s_index_%s", tableName, strings.ReplaceAll(fieldPath, ".", "_"))
				}
				goName := indexGoName(indexName, tableName)
				*entries = append(*entries, indexEntry{
					IndexName: indexName,
					GoName:    goName,
					TableName: tableName,
				})
			}
		}

		if searchInfo != nil && searchInfo.ConfigName != "" {
			searchDef := b.findSearchConfig(searchInfo.ConfigName)
			if searchDef != nil {
				indexName := fmt.Sprintf(def.IndexPrefix+"%s_search_%s", tableName, strings.ReplaceAll(fieldPath, ".", "_"))
				goName := indexGoName(indexName, tableName)
				*entries = append(*entries, indexEntry{
					IndexName: indexName,
					GoName:    goName,
					TableName: tableName,
				})
			}
		}

		if nestedFields := f.NestedFields(); nestedFields != nil {
			b.collectFieldIndexEntries(tableName, fieldPath, nestedFields, entries, compositeUnique)
		}
	}
}

func indexGoName(indexName, tableName string) string {
	prefix := def.IndexPrefix + tableName + "_"
	remainder := strings.TrimPrefix(indexName, prefix)
	if remainder == indexName {
		return strcase.ToCamel(indexName)
	}
	return strcase.ToCamel(remainder)
}

func (b *indexBuilder) findSearchConfig(name string) *parser.SearchDef {
	if b.define == nil {
		return nil
	}
	for i := range b.define.Searches {
		if b.define.Searches[i].Name == name {
			return &b.define.Searches[i]
		}
	}
	return nil
}
