package fs

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/maps"
	"io"
	"os"
	"path/filepath"
	"slices"
)

type FS struct {
	writes map[string][]byte
	writer map[string]*bytes.Buffer
}

func New() *FS {
	return &FS{
		writes: make(map[string][]byte),
		writer: make(map[string]*bytes.Buffer),
	}
}

func (fs *FS) Write(path string, content []byte) {
	fs.writes[path] = content
}

func (fs *FS) Writer(path string) io.Writer {
	if _, ok := fs.writer[path]; !ok {
		fs.writer[path] = bytes.NewBuffer(nil)
	}

	return fs.writer[path]
}

func (fs *FS) Flush(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", dir, err)
	}

	dirs := make(map[string]struct{})

	for path := range fs.writes {
		dirs[filepath.Dir(filepath.Join(dir, path))] = struct{}{}
	}

	for path := range fs.writer {
		dirs[filepath.Dir(filepath.Join(dir, path))] = struct{}{}
	}

	for _, dir := range maps.Keys(dirs) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create dir %s: %w", dir, err)
		}
	}

	for path, content := range fs.writes {
		if err := os.WriteFile(filepath.Join(dir, path), content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", path, err)
		}
	}

	for path, writer := range fs.writer {
		if err := os.WriteFile(filepath.Join(dir, path), writer.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", path, err)
		}
	}

	return nil
}

func (fs *FS) Dry(dir string) error {
	dir = filepath.Clean(dir)

	var files []string

	files = append(files, maps.Keys(fs.writes)...)
	files = append(files, maps.Keys(fs.writer)...)

	slices.Sort(files)

	for _, file := range files {
		fmt.Println(filepath.Join(dir, file))
	}

	return nil
}
