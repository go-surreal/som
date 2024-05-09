package embed

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"text/template"
)

//go:embed som/*
var somContent embed.FS

//go:embed conv/*
var convContent embed.FS

//go:embed fetch/*
var fetchContent embed.FS

//go:embed query/*
var queryContent embed.FS

//go:embed relate/*
var relateContent embed.FS

//go:embed sort/*
var sortContent embed.FS

//go:embed lib/*
var libContent embed.FS

type Template struct {
	GenerateOutPath string
}

type File struct {
	Path    string
	Content []byte
}

func Som(tmpl *Template) ([]*File, error) {
	return readEmbed(somContent, "som", tmpl)
}

func Conv(tmpl *Template) ([]*File, error) {
	return readEmbed(convContent, "conv", tmpl)
}

func Fetch(tmpl *Template) ([]*File, error) {
	return readEmbed(fetchContent, "fetch", tmpl)
}

func Query(tmpl *Template) ([]*File, error) {
	return readEmbed(queryContent, "query", tmpl)
}

func Relate(tmpl *Template) ([]*File, error) {
	return readEmbed(relateContent, "relate", tmpl)
}

func Sort(tmpl *Template) ([]*File, error) {
	return readEmbed(sortContent, "sort", tmpl)
}

func Lib(tmpl *Template) ([]*File, error) {
	return readEmbed(libContent, "lib", tmpl)
}

func readEmbed(fs embed.FS, name string, tmpl *Template) ([]*File, error) {
	dir, err := fs.ReadDir(name)
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, entry := range dir {
		if entry.IsDir() {
			return nil, fmt.Errorf("%s package contains unexpected directory", name)
		}

		filePath := filepath.Join(name, entry.Name())

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
