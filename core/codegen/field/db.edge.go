package field

import (
	"github.com/iancoleman/strcase"
)

type DatabaseEdge struct {
	Name   string
	In     Field
	Out    Field
	Fields []Field
}

func (e *DatabaseEdge) FileName() string {
	return "edge." + strcase.ToSnake(e.Name) + ".go"
}

func (e *DatabaseEdge) GetFields() []Field {
	return e.Fields
}

func (e *DatabaseEdge) NameGo() string {
	return e.Name
}

func (e *DatabaseEdge) NameGoLower() string {
	return strcase.ToLowerCamel(e.Name)
}

func (e *DatabaseEdge) NameDatabase() string {
	return strcase.ToSnake(e.Name) // TODO
}
