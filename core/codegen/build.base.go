package codegen

import (
	"os"
	"path"
)

type builder interface {
	build() error
}

type baseBuilder struct {
	*input

	// basePath holds the base path for all the generated code.
	basePath string

	// basePkg holds the base package path for all the generated code.
	basePkg string

	// pkgName holds the name of the generated package.
	pkgName string
}

func newBaseBuilder(input *input, basePath, basePkg, pkgName string) *baseBuilder {
	return &baseBuilder{
		input:    input,
		basePath: basePath,
		basePkg:  basePkg,
		pkgName:  pkgName,
	}
}

func (b *baseBuilder) path() string {
	return path.Join(b.basePath, b.pkgName)
}

func (b *baseBuilder) pkg() string {
	return path.Join(b.basePkg, b.pkgName)
}

func (b *baseBuilder) subPkg(pkg string) string {
	return path.Join(b.basePkg, pkg)
}

// createDir creates the directory for the generated files.
func (b *baseBuilder) createDir() error {
	return os.MkdirAll(b.path(), os.ModePerm)
}
