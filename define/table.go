package define

import (
	"github.com/marcbinz/sdb/define/field"
)

type TableDef struct {
	name   string
	fields []*field.Field
}

func (d *TableDef) With(fields ...*field.Field) *TableDef {
	d.fields = fields
	return d
}

func (d *TableDef) render() string {
	return "table def"
}
