package codegen

import (
	"github.com/go-surreal/som/core/util/fs"
	"path"
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

func (b *baseBuilder) subPkg(pkg string) string {
	return path.Join(b.basePkg, pkg)
}
