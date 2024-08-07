package basic

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	genComment = `// Code generated by github.com/go-surreal/som, DO NOT EDIT.`
)

func TestGenComments(t *testing.T) {
	err := filepath.WalkDir("./gen/som", func(path string, d fs.DirEntry, subErr error) error {
		if d.IsDir() {
			return nil
		}

		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if !strings.HasPrefix(string(file), genComment) {
			t.Errorf("file %s does not contain codegen comment", path)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
