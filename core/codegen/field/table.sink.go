package field

import (
	"github.com/iancoleman/strcase"
)

// SinkTable is a write-only ingestion table (DEFINE TABLE ... DROP). It is
// structurally node-like (a set of columns), but has no read operations,
// no ID type and no features: records are discarded immediately after
// write.
type SinkTable struct {
	Name   string
	Fields []Field
}

func (t *SinkTable) FileName() string {
	return "sink." + strcase.ToSnake(t.Name) + ".go"
}

func (t *SinkTable) GetFields() []Field {
	return t.Fields
}

func (t *SinkTable) NameGo() string {
	return t.Name
}

func (t *SinkTable) NameGoLower() string {
	return strcase.ToLowerCamel(t.Name)
}

func (t *SinkTable) NameDatabase() string {
	return strcase.ToSnake(t.Name)
}
