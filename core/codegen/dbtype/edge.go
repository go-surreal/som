package dbtype

import (
	"github.com/iancoleman/strcase"
)

type Edge struct {
	Name   string
	In     Field
	Out    Field
	Fields []Field
}

func (n *Edge) FileName() string {
	return "edge." + strcase.ToSnake(n.Name) + ".go"
}

func (n *Edge) GetFields() []Field {
	return n.Fields
}

func (n *Edge) NameGo() string {
	return n.Name
}

func (n *Edge) NameDatabase() string {
	return strcase.ToSnake(n.Name) // TODO
}
