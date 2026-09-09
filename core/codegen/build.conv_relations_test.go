package codegen

import (
	"strings"
	"testing"
)

func TestCheckRelationLimit(t *testing.T) {
	if err := checkRelationLimit("User", maxRelations); err != nil {
		t.Errorf("expected %d relations to be allowed, got: %v", maxRelations, err)
	}

	err := checkRelationLimit("User", maxRelations+1)
	if err == nil {
		t.Fatalf("expected %d relations to be rejected", maxRelations+1)
	}

	if !strings.Contains(err.Error(), "User") {
		t.Errorf("expected the model name in the error, got: %v", err)
	}
}
