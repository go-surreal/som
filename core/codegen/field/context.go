package field

import (
	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
	"path"
)

type Context struct {
	SourcePkg   string
	TargetPkg   string
	Table       Table
	Receiver    *jen.Statement
	isFromSlice bool
}

func (c Context) pkgLib() string {
	return path.Join(c.TargetPkg, def.PkgLib)
}

func (c Context) pkgTypes() string {
	return path.Join(c.TargetPkg, def.PkgTypes)
}

func (c Context) pkgCBOR() string {
	return path.Join(c.TargetPkg, def.PkgCBORHelpers)
}

func (c Context) fromSlice() Context {
	c.isFromSlice = true
	return c
}
