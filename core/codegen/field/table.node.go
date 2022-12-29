package field

import (
	"github.com/iancoleman/strcase"
)

type NodeTable struct {
	Name       string
	Fields     []Field
	Timestamps bool
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
	return strcase.ToLowerCamel(t.Name)
}

func (t *NodeTable) NameDatabase() string {
	return strcase.ToSnake(t.Name) // TODO
}

func (t *NodeTable) HasTimestamps() bool {
	return t.Timestamps
}
