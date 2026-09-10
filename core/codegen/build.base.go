package codegen

import (
	"path"

	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/util/fs"
)

type builder interface {
	build() error
}

type baseBuilder struct {
	*input

	// fs is the in-memory file system all generated file should be written to.
	fs *fs.FS

	// basePkg holds the base package path for all the generated code.
	basePkg string

	// pkgName holds the name of the generated package.
	pkgName string
}

func newBaseBuilder(input *input, fs *fs.FS, basePkg, pkgName string) *baseBuilder {
	return &baseBuilder{
		input:   input,
		fs:      fs,
		basePkg: basePkg,
		pkgName: pkgName,
	}
}

// TODO: rename to pkg()
func (b *baseBuilder) path() string {
	return b.pkgName
}

func (b *baseBuilder) relativePkgPath(pkg ...string) string {
	return path.Join(append([]string{b.basePkg}, pkg...)...)
}

// definedFields returns the fields an accessor struct is built from. In
// contrast to GetFields, these include the fields of the embedded som types.
func definedFields(elem field.Element) []field.Field {
	switch elem := elem.(type) {
	case *field.NodeTable:
		return elem.Fields

	case *field.DatabaseObject:
		return elem.Fields

	default:
		return elem.GetFields()
	}
}

func fieldContextFor(sourcePkg, targetPkg string, elem field.Element, index int) field.Context {
	ctx := field.Context{
		SourcePkg: sourcePkg,
		TargetPkg: targetPkg,
		Table:     elem,
	}

	// The fields of an array based complex ID are addressed by their position.
	if object, ok := elem.(*field.DatabaseObject); ok && object.IsArrayIndexed {
		ctx.ArrayIndex = &index
	}

	return ctx
}
