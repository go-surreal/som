package dbtype

import (
	"github.com/iancoleman/strcase"
)

type Object struct {
	Name   string
	Fields []Field
}

func (o *Object) FileName() string {
	return "object." + strcase.ToSnake(o.Name) + ".go"
}

func (o *Object) GetFields() []Field {
	return o.Fields
}

func (o *Object) NameGo() string {
	return o.Name
}

func (o *Object) NameDatabase() string {
	return strcase.ToSnake(o.Name) // TODO
}
