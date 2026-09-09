package field

import (
	"github.com/iancoleman/strcase"
)

// FragmentTable is a projection of a node table: a subset of the node's fields
// that is fetched instead of the full record. It has no table of its own, so
// it neither takes part in the schema nor has write operations.
type FragmentTable struct {
	Name   string
	Parent *NodeTable
	Fields []Field
}

func (t *FragmentTable) FileName() string {
	return "fragment." + strcase.ToSnake(t.Name) + ".go"
}

func (t *FragmentTable) GetFields() []Field {
	return t.Fields
}

func (t *FragmentTable) NameGo() string {
	return t.Name
}

func (t *FragmentTable) NameGoLower() string {
	return strcase.ToLowerCamel(t.Name)
}

// NameDatabase returns the database name of the table the fragment is
// projected from, since a fragment has no table of its own.
func (t *FragmentTable) NameDatabase() string {
	return t.Parent.NameDatabase()
}

// Projection returns the database fields the fragment selects, id included.
func (t *FragmentTable) Projection() []string {
	fields := make([]string, 0, len(t.Fields)+1)
	fields = append(fields, "id")

	for _, f := range t.Fields {
		if f.NameDatabase() == "id" {
			continue
		}
		fields = append(fields, f.NameDatabase())
	}

	return fields
}
