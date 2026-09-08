package field

import (
	"github.com/iancoleman/strcase"
)

// ViewTable is a read-only, pre-computed table view. It is structurally
// node-like (a set of projected columns), but has no write operations,
// no ID type and no features.
type ViewTable struct {
	Name   string
	Fields []Field
}

func (t *ViewTable) FileName() string {
	return "view." + strcase.ToSnake(t.Name) + ".go"
}

func (t *ViewTable) GetFields() []Field {
	return t.Fields
}

func (t *ViewTable) NameGo() string {
	return t.Name
}

func (t *ViewTable) NameGoLower() string {
	return strcase.ToLowerCamel(t.Name)
}

func (t *ViewTable) NameDatabase() string {
	return strcase.ToSnake(t.Name)
}
