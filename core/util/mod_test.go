package util

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

	msg, err = mod.CheckSOMVersion()
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
