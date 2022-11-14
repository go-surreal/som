package dbtype

import (
	"github.com/iancoleman/strcase"
)

type Node struct {
	Name   string
	Fields []Field
}

func (n *Node) FileName() string {
	return "node." + strcase.ToSnake(n.Name) + ".go"
}

func (n *Node) GetFields() []Field {
	return n.Fields
}

func (n *Node) NameGo() string {
	return n.Name
}

func (n *Node) NameDatabase() string {
	return strcase.ToSnake(n.Name) // TODO
}
