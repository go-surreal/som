package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"path/filepath"
)

type Context struct {
	SourcePkg   string
	TargetPkg   string
	Table       Table
	Receiver    *jen.Statement
	isFromSlice bool
}

func (c Context) pkgLib() string {
	return filepath.Join(c.TargetPkg, def.PkgLib)
}

func (c Context) pkgTypes() string {
	return filepath.Join(c.TargetPkg, def.PkgTypes)
}

func (c Context) fromSlice() Context {
	c.isFromSlice = true
	return c
}
