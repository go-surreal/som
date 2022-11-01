package parser

import (
	"errors"
	"golang.org/x/mod/modfile"
	"os"
	"path"
)

func findAndReadModFile(dir string) ([]byte, string, error) {
	for dir != "" {
		data, err := os.ReadFile(path.Join(dir, fileGoMod))

		if err == nil {
			return data, dir, nil
		}

		if !errors.Is(err, os.ErrNotExist) {
			return nil, "", err
		}

		dir = path.Dir(dir)
	}

	return nil, "", errors.New("could not find go.mod in worktree")
}

func parseMod(dir string) (string, string, error) {
	data, filePath, err := findAndReadModFile(dir)
	if err != nil {
		return "", "", err
	}

	f, err := modfile.Parse(fileGoMod, data, nil)
	if err != nil {
		return "", "", err
	}

	return f.Module.Mod.Path, filePath, nil
}
