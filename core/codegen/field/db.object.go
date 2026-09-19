package field

import (
	"github.com/iancoleman/strcase"
)

type DatabaseObject struct {
	Name           string
	Fields         []Field
	IsArrayIndexed bool

	// HasValidate reports whether the struct type provides a
	// "Validate() error" method of its own.
	HasValidate bool
}

func (o *DatabaseObject) NameGo() string {
	return o.Name
}

func (o *DatabaseObject) NameGoLower() string {
	return strcase.ToLowerCamel(o.Name)
}

func (o *DatabaseObject) NameDatabase() string {
	return strcase.ToSnake(o.Name) // TODO
}

func (o *DatabaseObject) FileName() string {
	return "object." + strcase.ToSnake(o.Name) + ".go"
}

func (o *DatabaseObject) GetFields() []Field {
	return o.Fields
}
