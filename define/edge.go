package define

import (
	"github.com/marcbinz/sdb/define/field"
)

type EdgeDef struct {
	name   string
	from   *TableDef
	to     *TableDef
	fields []*field.Field
}

func (d *EdgeDef) From(table *TableDef) *EdgeDef {
	d.from = table
	return d
}

func (d *EdgeDef) To(table *TableDef) *EdgeDef {
	d.to = table
	return d
}

func (d *EdgeDef) With(fields ...*field.Field) *EdgeDef {
	d.fields = fields
	return d
}

func (d *EdgeDef) render() string {
	return "edge def"
}
