package field

import (
	"github.com/iancoleman/strcase"
)

type EdgeTable struct {
	Name       string
	In         *Node
	Out        *Node
	Fields     []Field
	Changefeed string
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

func (t *EdgeTable) HasChangefeed() bool {
	return t.Changefeed != ""
}
