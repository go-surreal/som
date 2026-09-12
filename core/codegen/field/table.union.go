package field

import (
	"strings"

	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
)

// UnionTable is a set of node tables reachable through a common interface. It
// is not a table of its own: it exists to type multi-table record links
// (record<a|b>) and relation endpoints, and to query all of its members at
// once.
type UnionTable struct {
	Name    string
	Members []*NodeTable
	Source  *parser.Union
}

func (t *UnionTable) FileName() string {
	return "union." + strcase.ToSnake(t.Name) + ".go"
}

func (t *UnionTable) GetFields() []Field {
	return nil
}

func (t *UnionTable) NameGo() string {
	return t.Name
}

func (t *UnionTable) NameGoLower() string {
	return strcase.ToLowerCamel(t.Name)
}

func (t *UnionTable) NameDatabase() string {
	return strcase.ToSnake(t.Name)
}

// TypeDatabase returns the record type covering all members of the union,
// e.g. "record<square|triangle>".
func (t *UnionTable) TypeDatabase() string {
	return "record<" + strings.Join(t.memberNames(), "|") + ">"
}

// QueryDatabase returns the query target selecting from all member tables,
// e.g. "square, triangle".
func (t *UnionTable) QueryDatabase() string {
	return strings.Join(t.memberNames(), ", ")
}

func (t *UnionTable) memberNames() []string {
	names := make([]string, len(t.Members))
	for i, member := range t.Members {
		names[i] = member.NameDatabase()
	}
	return names
}
