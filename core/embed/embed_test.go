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

func TestReadVersionedQueryMethod(t *testing.T) {
	t.Parallel()

	const versionMethod = "func (b builder[M, R]) Version(at time.Time)"

	for _, versioned := range []bool{false, true} {
		files, err := Read(&Template{GenerateOutPath: "som.test/gen/som", Versioned: versioned})
		if err != nil {
			t.Fatal(err)
		}

		builder := findFile(t, files, "query/builder.go")

		if hasMethod := strings.Contains(string(builder.Content), versionMethod); hasMethod != versioned {
			t.Errorf("versioned=%v: Version method present=%v", versioned, hasMethod)
		}
	}
}

func findFile(t *testing.T, files []*File, path string) *File {
	t.Helper()

	for _, file := range files {
		if file.Path == path {
			return file
		}
	}

	t.Fatalf("file %s not found", path)

	return nil
}
