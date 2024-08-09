package fs

import (
	"bytes"
	"fmt"
	"golang.org/x/exp/maps"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
	//if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
	//	return fmt.Errorf("failed to remove %s: %w", dir, err)
	//}
	//
	//path := dir
	//
	//for {
	//	path = filepath.Dir(path)
	//
	//	entries, err := os.ReadDir(path)
	//	if err != nil && !os.IsNotExist(err) {
	//		return fmt.Errorf("failed to read dir %s: %w", path, err)
	//	}
	//
	//	if len(entries) > 0 {
	//		break
	//	}
	//
	//	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
	//		return fmt.Errorf("failed to remove dir %s: %w", path, err)
	//	}
	//}

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

	allFiles := fs.allFiles()

	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !slices.Contains(allFiles, strings.TrimPrefix(path, dir+string(os.PathSeparator))) {
			return os.Remove(path)
		}

		return nil
	})
}

func (fs *FS) Dry(dir string) error {
	dir = filepath.Clean(dir)

	allFiles := fs.allFiles()

	changes := make(map[string]*bool)

	valTrue := true
	valFalse := false

	for _, file := range allFiles {
		changes[file] = &valTrue
	}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil // TODO: mark directories that are removed?
		}

		path = strings.TrimPrefix(path, dir+string(os.PathSeparator))

		if slices.Contains(allFiles, path) {
			changes[path] = nil
			return nil
		}

		changes[path] = &valFalse

		allFiles = append(allFiles, path)

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk dir %s: %w", dir, err)
	}

	slices.Sort(allFiles)

	for _, file := range allFiles {
		switch changes[file] {
		case nil:
			fmt.Println("× " + filepath.Join(dir, file)) // TODO: only if content changed? (±)
		case &valTrue:
			fmt.Println("+ " + filepath.Join(dir, file))
		case &valFalse:
			fmt.Println("- " + filepath.Join(dir, file))
		}
	}

	return nil
}

func (fs *FS) allFiles() []string {
	return append(
		maps.Keys(fs.writes),
		maps.Keys(fs.writer)...,
	)
}
