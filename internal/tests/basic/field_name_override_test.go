package basic

import (
	"context"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som/by"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestFieldNameOverride(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	created := &model.AllTypes{
		FieldRenamed: "test_value",
		FieldMonth:   time.January,
	}

	err := client.AllTypesRepo().Create(ctx, created)
	if err != nil {
		t.Fatal(err)
	}

	// Filter by the renamed field.
	result, err := client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldRenamed.Equal("test_value")).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "test_value", result.FieldRenamed)

	// Create a second record to test sorting.
	second := &model.AllTypes{
		FieldRenamed: "aaa_first",
		FieldMonth:   time.January,
	}

	err = client.AllTypesRepo().Create(ctx, second)
	if err != nil {
		t.Fatal(err)
	}

	// Sort ascending by the renamed field.
	sorted, err := client.AllTypesRepo().Query().
		OrderBy(by.AllTypes.FieldRenamed.Asc()).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, len(sorted))
	assert.Equal(t, "aaa_first", sorted[0].FieldRenamed)
	assert.Equal(t, "test_value", sorted[1].FieldRenamed)
}
