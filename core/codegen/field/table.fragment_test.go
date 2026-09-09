package field

import (
	"testing"

	"github.com/go-surreal/som/core/parser"
	"github.com/iancoleman/strcase"
	"gotest.tools/v3/assert"
)

func TestFragmentTable(t *testing.T) {
	conf := &BuildConfig{ToDatabaseName: strcase.ToSnake}

	fragment := FragmentTable{
		Name:   "PersonCard",
		Parent: &NodeTable{Name: "Person"},
		Fields: []Field{
			&String{baseField: &baseField{BuildConfig: conf, source: parser.NewFieldString("Name")}},
			&String{baseField: &baseField{BuildConfig: conf, source: parser.NewFieldString("HomeTown")}},
		},
	}

	assert.Equal(t, "fragment.person_card.go", fragment.FileName())
	assert.Equal(t, "PersonCard", fragment.NameGo())
	assert.Equal(t, "personCard", fragment.NameGoLower())

	// A fragment has no table of its own, it queries the one of its parent.
	assert.Equal(t, "person", fragment.NameDatabase())

	// The record id is always projected, and it comes first.
	assert.DeepEqual(t, []string{"id", "name", "home_town"}, fragment.Projection())
}

// TestFragmentTableProjectionSkipsID makes sure the id is projected once, even
// if the parent node contributes an explicit id field.
func TestFragmentTableProjectionSkipsID(t *testing.T) {
	conf := &BuildConfig{ToDatabaseName: strcase.ToSnake}

	fragment := FragmentTable{
		Name:   "PersonCard",
		Parent: &NodeTable{Name: "Person"},
		Fields: []Field{
			&ID{baseField: &baseField{BuildConfig: conf, source: parser.NewFieldID("ID", parser.IDTypeULID)}, source: parser.NewFieldID("ID", parser.IDTypeULID)},
			&String{baseField: &baseField{BuildConfig: conf, source: parser.NewFieldString("Name")}},
		},
	}

	assert.DeepEqual(t, []string{"id", "name"}, fragment.Projection())
}
