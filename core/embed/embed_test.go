package embed

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	localImport = "github.com/go-surreal/som/"
)

func TestEmbedNoLocalImports(t *testing.T) {
	t.Parallel()

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".go" || path == "embed_test.go" {
			return nil
		}

		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if strings.Contains(string(file), localImport) {
			t.Errorf("File %s contains invalid import of type '%s'", path, localImport)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
