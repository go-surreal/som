package define

import (
	"github.com/marcbinz/sdb/define/field"
)

type ObjectDef struct {
	name   string
	fields []*field.Field
}

func (d *ObjectDef) With(fields ...*field.Field) *ObjectDef {
	d.fields = fields
	return d
}

func (d *ObjectDef) render() string {
	return "object def"
}
