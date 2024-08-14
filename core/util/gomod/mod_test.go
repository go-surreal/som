package gomod

import (
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

func TestGoModValid(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/go.valid.mod")
	if err != nil {
		t.Fatal(err)
	}

	mod, err := NewGoMod("go.mod", data)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := mod.CheckGoVersion()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", msg)

	msg, err = mod.CheckSOMVersion(false)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", msg)

	msg, err = mod.CheckSDBCVersion()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", msg)
}

func TestGoModUnsupportedGoVersion(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/go.invalid.mod")
	if err != nil {
		t.Fatal(err)
	}

	mod, err := NewGoMod("go.mod", data)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := mod.CheckGoVersion()

	assert.ErrorContains(t, err, "go version 1.12 is not supported")
	assert.Equal(t, "", msg)
}

func TestGoModMissingSOMPackage(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/go.invalid.mod")
	if err != nil {
		t.Fatal(err)
	}

	mod, err := NewGoMod("go.mod", data)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := mod.CheckSOMVersion(false)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", msg)

	exists := false

	for _, req := range mod.file.Require {
		if req.Mod.Path != pkgSOM {
			continue
		}

		exists = true

		if req.Mod.Version != requiredSOMVersion {
			t.Fatal("som version not updated")
		}
	}

	assert.Assert(t, exists)
}

func TestGoModWrongSOMVersion(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/go.outdated.mod")
	if err != nil {
		t.Fatal(err)
	}

	mod, err := NewGoMod("go.mod", data)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := mod.CheckSOMVersion(false)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", msg)

	exists := false

	for _, req := range mod.file.Require {
		if req.Mod.Path != pkgSOM {
			continue
		}

		exists = true

		if req.Mod.Version != requiredSOMVersion {
			t.Fatal("som version not updated")
		}
	}

	assert.Assert(t, exists)
}

func TestGoModMissingSDBCPackage(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/go.invalid.mod")
	if err != nil {
		t.Fatal(err)
	}

	mod, err := NewGoMod("go.mod", data)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := mod.CheckSDBCVersion()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", msg)

	exists := false

	for _, req := range mod.file.Require {
		if req.Mod.Path != pkgSDBC {
			continue
		}

		exists = true

		if req.Mod.Version != requiredSDBCVersion {
			t.Fatal("som version not updated")
		}
	}

	assert.Assert(t, exists)
}

func TestGoModWrongSDBCVersion(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("testdata/go.outdated.mod")
	if err != nil {
		t.Fatal(err)
	}

	mod, err := NewGoMod("go.mod", data)
	if err != nil {
		t.Fatal(err)
	}

	msg, err := mod.CheckSDBCVersion()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", msg)

	exists := false

	for _, req := range mod.file.Require {
		if req.Mod.Path != pkgSDBC {
			continue
		}

		exists = true

		if req.Mod.Version != requiredSDBCVersion {
			t.Fatal("som version not updated")
		}
	}

	assert.Assert(t, exists)
}
