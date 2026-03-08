package field

import (
	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
)

type EdgeTable struct {
	Name   string
	In     *Node
	Out    *Node
	Fields []Field
	Source *parser.Edge
}

func (t *EdgeTable) NameGo() string {
	return t.Name
}

func (t *EdgeTable) NameGoLower() string {
	return strcase.ToLowerCamel(t.Name)
}

func (t *EdgeTable) NameDatabase() string {
	return strcase.ToSnake(t.Name) // TODO
}

func (t *EdgeTable) FileName() string {
	return "edge." + strcase.ToSnake(t.Name) + ".go"
}

func (t *EdgeTable) GetFields() []Field {
	return t.Fields
}
