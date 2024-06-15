package util

import (
	"fmt"
	"os"
	"testing"
)

func TestNewGoMod(t *testing.T) {
	data, err := os.ReadFile("testdata/go.mod")
	if err != nil {
		t.Fatal(err)
	}

	mod, err := NewGoMod("go.mod", data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(mod)
}
