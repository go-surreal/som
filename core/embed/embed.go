package embed

import (
	"embed"
	"errors"
	"path/filepath"
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

type File struct {
	Path    string
	Content []byte
}

func Som() ([]*File, error) {
	dir, err := somContent.ReadDir("som")
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, entry := range dir {
		if entry.IsDir() {
			return nil, errors.New("som package contains unexpected directory")
		}

		filePath := filepath.Join("som", entry.Name())

		content, err := somContent.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		files = append(files, &File{
			Path:    entry.Name(),
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
			Path:    entry.Name(),
			Content: content,
		})
	}

	return files, nil
}

func Fetch() ([]*File, error) {
	dir, err := fetchContent.ReadDir("fetch")
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, entry := range dir {
		if entry.IsDir() {
			return nil, errors.New("fetch package contains unexpected directory")
		}

		filePath := filepath.Join("fetch", entry.Name())

		content, err := fetchContent.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		files = append(files, &File{
			Path:    entry.Name(),
			Content: content,
		})
	}

	return files, nil
}

func Query() ([]*File, error) {
	dir, err := queryContent.ReadDir("query")
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, entry := range dir {
		if entry.IsDir() {
			return nil, errors.New("query package contains unexpected directory")
		}

		filePath := filepath.Join("query", entry.Name())

		content, err := queryContent.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		files = append(files, &File{
			Path:    entry.Name(),
			Content: content,
		})
	}

	return files, nil
}

func Relate() ([]*File, error) {
	dir, err := relateContent.ReadDir("relate")
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, entry := range dir {
		if entry.IsDir() {
			return nil, errors.New("relate package contains unexpected directory")
		}

		filePath := filepath.Join("relate", entry.Name())

		content, err := relateContent.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		files = append(files, &File{
			Path:    entry.Name(),
			Content: content,
		})
	}

	return files, nil
}

func Sort() ([]*File, error) {
	dir, err := sortContent.ReadDir("sort")
	if err != nil {
		return nil, err
	}

	var files []*File

	for _, entry := range dir {
		if entry.IsDir() {
			return nil, errors.New("sort package contains unexpected directory")
		}

		filePath := filepath.Join("sort", entry.Name())

		content, err := sortContent.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		files = append(files, &File{
			Path:    entry.Name(),
			Content: content,
		})
	}

	return files, nil
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
