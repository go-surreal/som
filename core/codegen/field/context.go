package field

import (
	"github.com/dave/jennifer/jen"
)

type Context struct {
	SourcePkg string
	Table     Table
	Receiver  *jen.Statement
}
