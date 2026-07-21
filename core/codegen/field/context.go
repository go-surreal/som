package field

import (
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/go-surreal/som/core/codegen/def"
)

type Context struct {
	SourcePkg   string
	TargetPkg   string
	Table       Table
	Element     Element
	Receiver    *jen.Statement
	ArrayIndex  *int
	isFromSlice bool
}

func (c Context) pkgLib() string {
	return path.Join(c.TargetPkg, def.PkgLib)
}

func (c Context) pkgTypes() string {
	return path.Join(c.TargetPkg, def.PkgTypes)
}

func (c Context) pkgDistinct() string {
	return path.Join(c.TargetPkg, def.PkgDistinct)
}

func (c Context) pkgCBOR() string {
	return path.Join(c.TargetPkg, def.PkgCBORHelpers)
}

func (c Context) pkgInternal() string {
	return path.Join(c.TargetPkg, def.PkgInternal)
}

func (c Context) filterKeyCode(name string) jen.Code {
	if c.ArrayIndex != nil {
		return jen.Qual(c.pkgLib(), "Index").Call(jen.Id("key"), jen.Lit(*c.ArrayIndex))
	}
	return jen.Qual(c.pkgLib(), "Field").Call(jen.Id("key"), jen.Lit(name))
}

func (c Context) sortKeyCode(name string) jen.Code {
	if c.ArrayIndex != nil {
		return jen.Id("indexed").Call(jen.Id("key"), jen.Lit(*c.ArrayIndex))
	}
	return jen.Id("keyed").Call(jen.Id("key"), jen.Lit(name))
}

func (c Context) fromSlice() Context {
	c.isFromSlice = true
	return c
}
