package codegen

import (
	"testing"

	"github.com/go-surreal/som/core/codegen/field"
	"github.com/go-surreal/som/core/parser"
)

func TestBuildViewStatement(t *testing.T) {
	source := &field.NodeTable{Name: "AllTypes"}
	view := &field.ViewTable{Name: "AllTypesSummary"}

	b := &build{
		input: &input{
			nodes: []*field.NodeTable{source},
			views: []*field.ViewTable{view},
			define: &parser.DefineOutput{
				Views: []parser.ViewDef{{
					View:   "AllTypesSummary",
					Source: "AllTypes",
					Projections: []string{
						"field_string AS category",
						"count(field_string) AS total",
						"math::mean(field_float_64) AS avg_value",
					},
					Where:   "field_int > 0",
					GroupBy: []string{"field_string"},
				}},
			},
		},
	}

	got, err := b.buildViewStatement(view)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "DEFINE TABLE all_types_summary TYPE NORMAL AS SELECT " +
		"field_string AS category, count(field_string) AS total, math::mean(field_float_64) AS avg_value " +
		"FROM all_types WHERE field_int > 0 GROUP BY field_string;"

	if got != want {
		t.Errorf("view DDL mismatch:\n got: %s\nwant: %s", got, want)
	}
}

func TestBuildViewStatement_MissingDefinition(t *testing.T) {
	view := &field.ViewTable{Name: "Orphan"}

	b := &build{input: &input{
		views:  []*field.ViewTable{view},
		define: &parser.DefineOutput{},
	}}

	stmt, err := b.buildViewStatement(view)
	if err != nil {
		t.Fatalf("missing definition should skip, not error: %v", err)
	}
	if stmt != "" {
		t.Fatalf("expected empty statement for view without definition, got %q", stmt)
	}
}

func TestBuildViewStatement_DuplicateDefinition(t *testing.T) {
	source := &field.NodeTable{Name: "AllTypes"}
	view := &field.ViewTable{Name: "Dup"}

	b := &build{input: &input{
		nodes: []*field.NodeTable{source},
		views: []*field.ViewTable{view},
		define: &parser.DefineOutput{Views: []parser.ViewDef{
			{View: "Dup", Source: "AllTypes", Projections: []string{"count(field_string) AS total"}},
			{View: "Dup", Source: "AllTypes", Projections: []string{"count(field_int) AS total"}},
		}},
	}}

	if _, err := b.buildViewStatement(view); err == nil {
		t.Fatal("expected error for duplicate view definition, got nil")
	}
}
