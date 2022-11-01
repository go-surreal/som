package dbtype

import (
	"github.com/marcbinz/sdb/core/codegen/field"
)

type Element interface {
	FileName() string
	GetFields() []field.Field
	NameGo() string
	NameDatabase() string
}
