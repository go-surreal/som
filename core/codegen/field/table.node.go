package field

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
)

type NodeTable struct {
	Name   string
	Fields []Field
	Source *parser.Node // Reference to source parser.Node

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
	return strcase.ToLowerCamel(t.Name)
}

func (t *NodeTable) NameDatabase() string {
	return strcase.ToSnake(t.Name) // TODO
}

func (t *NodeTable) HasComplexID() bool {
	return t.Source != nil && t.Source.ComplexID != nil
}
