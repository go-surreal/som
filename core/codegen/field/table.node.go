package field

import (
	"github.com/iancoleman/strcase"
	"strings"
)

type NodeTable struct {
	Name       string
	Fields     []Field
	Timestamps bool

	// TODO: include source package path + method(s)
}

func (t *NodeTable) FileName() string {
	return "node." + strcase.ToSnake(t.Name) + ".go"
}

func (t *NodeTable) GetFields() []Field {
	return t.Fields
}

func (t *NodeTable) NameGo() string {
	return t.Name
}

func (t *NodeTable) NameGoLower() string {
	return strings.ToLower(t.Name[:1]) + t.Name[1:]
}

func (t *NodeTable) NameDatabase() string {
	return strcase.ToSnake(t.Name) // TODO
}

func (t *NodeTable) HasTimestamps() bool {
	return t.Timestamps
}
