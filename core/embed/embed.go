package embed

import (
	"embed"
	"errors"
	"path/filepath"
)

//go:embed lib/*
var libContent embed.FS

//go:embed conv/*
var convContent embed.FS

type File struct {
	Path    string
	Content []byte
}

func Lib() ([]*File, error) {
	dir, err := libContent.ReadDir("lib")
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, entry := range dir {
		if entry.IsDir() {
			return nil, errors.New("lib package contains unexpected directory")
		}

		filePath := filepath.Join("lib", entry.Name())

		content, err := libContent.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		files = append(files, &File{
			Path:    filePath,
			Content: content,
		})
	}

	return files, nil
}

func Conv() ([]*File, error) {
	dir, err := convContent.ReadDir("conv")
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, entry := range dir {
		if entry.IsDir() {
			return nil, errors.New("conv package contains unexpected directory")
		}

		filePath := filepath.Join("conv", entry.Name())

		content, err := convContent.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		files = append(files, &File{
			Path:    filePath,
			Content: content,
		})
	}

	return files, nil
}
