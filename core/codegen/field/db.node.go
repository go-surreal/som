package field

import (
	"github.com/iancoleman/strcase"
)

type DatabaseNode struct {
	Name   string
	Fields []Field
}

func (n *DatabaseNode) FileName() string {
	return "node." + strcase.ToSnake(n.Name) + ".go"
}

func (n *DatabaseNode) GetFields() []Field {
	return n.Fields
}

func (n *DatabaseNode) NameGo() string {
	return n.Name
}

func (n *DatabaseNode) NameGoLower() string {
	return strcase.ToLowerCamel(n.Name)
}

func (n *DatabaseNode) NameDatabase() string {
	return strcase.ToSnake(n.Name) // TODO
}
