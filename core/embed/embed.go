package embed

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"text/template"
)

const (
	embedDir = "_embed"
)

//go:embed _embed/*
var fs embed.FS

type Template struct {
	GenerateOutPath string
}

type File struct {
	Path    string
	Content []byte
}

func Som(tmpl *Template) ([]*File, error) {
	return readEmbed("som", tmpl)
}

func Conv(tmpl *Template) ([]*File, error) {
	return readEmbed("conv", tmpl)
}

func Fetch(tmpl *Template) ([]*File, error) {
	return readEmbed("fetch", tmpl)
}

func Query(tmpl *Template) ([]*File, error) {
	return readEmbed("query", tmpl)
}

func Relate(tmpl *Template) ([]*File, error) {
	return readEmbed("relate", tmpl)
}

func Sort(tmpl *Template) ([]*File, error) {
	return readEmbed("sort", tmpl)
}

func Define(tmpl *Template) ([]*File, error) {
	return readEmbed("define", tmpl)
}

func Lib(tmpl *Template) ([]*File, error) {
	return readEmbed("lib", tmpl)
}

func readEmbed(name string, tmpl *Template) ([]*File, error) {
	dir, err := fs.ReadDir(filepath.Join(embedDir, name))
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, entry := range dir {
		if entry.IsDir() {
			return nil, fmt.Errorf("%s package contains unexpected directory", name)
		}

		filePath := filepath.Join(embedDir, name, entry.Name())

		content, err := fs.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		fileTmpl, err := template.New(filePath).Parse(string(content))
		if err != nil {
			return nil, err
		}

		var parsedContent bytes.Buffer

		if err := fileTmpl.Execute(&parsedContent, tmpl); err != nil {
			return nil, err
		}

		files = append(files, &File{
			Path:    entry.Name(),
			Content: parsedContent.Bytes(),
		})
	}

	return files, nil
}
