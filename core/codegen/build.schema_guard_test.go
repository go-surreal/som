package codegen

import (
	"strings"
	"testing"
)

func TestGuardRebuild(t *testing.T) {
	got := guardRebuild("index:__som__all_types_count",
		"DEFINE INDEX OVERWRITE __som__all_types_count ON all_types COUNT;",
	)

	want := `IF (SELECT VALUE hash FROM ONLY type::record("__som__schema", "index:__som__all_types_count")) != "bfc3c43ba8d53a6934ef91706d4e904638a65845f538ed9039036ace83077509" {
	DEFINE INDEX OVERWRITE __som__all_types_count ON all_types COUNT;
	UPSERT type::record("__som__schema", "index:__som__all_types_count") SET hash = "bfc3c43ba8d53a6934ef91706d4e904638a65845f538ed9039036ace83077509";
};`

	if got != want {
		t.Errorf("unexpected statement:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestGuardRebuildHashDependsOnStatements(t *testing.T) {
	base := guardRebuild("view:summary", "REMOVE TABLE IF EXISTS summary;", "DEFINE TABLE summary TYPE NORMAL AS SELECT a FROM b;")
	changed := guardRebuild("view:summary", "REMOVE TABLE IF EXISTS summary;", "DEFINE TABLE summary TYPE NORMAL AS SELECT a, c FROM b;")

	if hashOf(t, base) == hashOf(t, changed) {
		t.Error("expected differing definitions to produce differing hashes")
	}

	if again := guardRebuild("view:summary", "REMOVE TABLE IF EXISTS summary;", "DEFINE TABLE summary TYPE NORMAL AS SELECT a FROM b;"); again != base {
		t.Error("expected the same definition to produce a stable statement")
	}
}

func hashOf(tb testing.TB, statement string) string {
	tb.Helper()

	_, rest, ok := strings.Cut(statement, "!= \"")
	if !ok {
		tb.Fatalf("no hash found in statement: %s", statement)
	}

	hash, _, ok := strings.Cut(rest, "\"")
	if !ok {
		tb.Fatalf("no hash found in statement: %s", statement)
	}

	return hash
}
