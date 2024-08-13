package fs

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type FS struct {
	mem afero.Fs
}

func New() *FS {
	return &FS{
		mem: afero.NewMemMapFs(),
	}
}

func (fs *FS) Open(name string) (fs.File, error) {
	file, err := fs.mem.Open(name)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fs *FS) Write(path string, content []byte) error {
	file, err := fs.Writer(path)
	if err != nil {
		return err
	}

	if _, err := file.Write(content); err != nil {
		return err
	}

	return nil
}

func (fs *FS) Writer(path string) (io.Writer, error) {
	file, err := fs.mem.Create(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fs *FS) Flush(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", dir, err)
	}

	if err := os.CopyFS(dir, fs); err != nil {
		return fmt.Errorf("failed to write dir %s: %w", dir, err)
	}

	allFiles, err := fs.allFiles()
	if err != nil {
		return err
	}

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

	allFiles, err := fs.allFiles()
	if err != nil {
		return err
	}

	changes := make(map[string]*bool)

	valTrue := true
	valFalse := false

	for _, file := range allFiles {
		changes[file] = &valTrue
	}

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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

func (fs *FS) Clear(dir string) error {
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove %s: %w", dir, err)
	}

	path := dir

	for {
		path = filepath.Dir(path)

		entries, err := os.ReadDir(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to read dir %s: %w", path, err)
		}

		if len(entries) > 0 {
			break
		}

		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove dir %s: %w", path, err)
		}
	}

	return nil
}

func (fs *FS) allFiles() ([]string, error) {
	fileInfos, err := afero.ReadDir(fs.mem, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read files: %w", err)
	}

	var allFiles []string

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}

		allFiles = append(allFiles, fileInfo.Name())
	}

	return allFiles, nil
}
