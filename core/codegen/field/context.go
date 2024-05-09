package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"path/filepath"
)

type Context struct {
	SourcePkg string
	TargetPkg string
	Table     Table
	Receiver  *jen.Statement
}

func (c Context) pkgLib() string {
	return filepath.Join(c.TargetPkg, def.PkgLib)
}
