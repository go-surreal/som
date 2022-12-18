package field

import (
	"github.com/marcbinz/som/core/parser"
)

const (
	funcBuildDatabaseID = "buildDatabaseID"
	funcParseDatabaseID = "parseDatabaseID"
)

type Edge struct {
	Name   string
	In     Field
	Out    Field
	Fields []Field
}

type ElemGetter func(name string) (Element, bool)

type Field interface {
	NameGo() string
	NameDatabase() string

	CodeGen() *CodeGen
}

type Element interface {
	FileName() string
	GetFields() []Field
	NameGo() string
	NameGoLower() string
	NameDatabase() string
}

type BuildConfig struct {
	ToDatabaseName func(base string) string
}

type baseField struct {
	*BuildConfig

	source parser.Field
}

func (f *baseField) NameGo() string {
	return f.source.Name()
}

func (f *baseField) NameDatabase() string {
	return f.ToDatabaseName(f.source.Name())
}
