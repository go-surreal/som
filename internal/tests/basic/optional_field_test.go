package basic

import (
	"context"
	"testing"
	"time"

	"som.test/gen/som/filter"
	"som.test/model"
	"gotest.tools/v3/assert"
)

// TestOptionalFieldClear verifies that clearing a previously-set pointer field
// (setting it to nil and updating) persists as NONE and reads back as nil.
// This exercises the omit-on-marshal + full-content-replace path.
func TestOptionalFieldClear(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	num := 42

	user := &model.AllTypes{
		FieldString: "clear-me",
		FieldMonth:  time.January,
		FieldIntPtr: &num,
	}

	err := client.AllTypesRepo().Create(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	out, err := client.AllTypesRepo().Query().
		Where(filter.AllTypes.ID.Equal(string(user.ID()))).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, out.FieldIntPtr != nil && *out.FieldIntPtr == num)

	out.FieldIntPtr = nil

	err = client.AllTypesRepo().Update(ctx, out)
	if err != nil {
		t.Fatal(err)
	}

	out, err = client.AllTypesRepo().Query().
		Where(filter.AllTypes.ID.Equal(string(user.ID()))).
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, out.FieldIntPtr == nil)
}

// TestOptionalFieldNotEqual verifies that NotEqual/NotIn on an optional field
// does NOT match records where the field is unset (NONE). In SurrealDB
// "NONE != value" evaluates to true, so without the NONE/NULL guard these
// filters would wrongly include unset records.
func TestOptionalFieldNotEqual(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	const marker = "not-equal-marker"

	num := 10

	set := &model.AllTypes{
		FieldString: marker,
		FieldMonth:  time.January,
		FieldIntPtr: &num,
	}

	unset := &model.AllTypes{
		FieldString: marker,
		FieldMonth:  time.January,
		FieldIntPtr: nil,
	}

	for _, r := range []*model.AllTypes{set, unset} {
		if err := client.AllTypesRepo().Create(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	other := 999

	// NotEqual must return only the set record, not the unset one.
	results, err := client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldString.Equal(marker),
			filter.AllTypes.FieldIntPtr.NotEqual(&other),
		).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(results) == 1, "expected 1 result, got %d", len(results))
	if len(results) == 1 {
		assert.Check(t, results[0].FieldIntPtr != nil && *results[0].FieldIntPtr == num)
	}

	// NotIn must likewise exclude the unset record.
	results, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldString.Equal(marker),
			filter.AllTypes.FieldIntPtr.NotIn([]*int{&other}),
		).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(results) == 1, "expected 1 result for NotIn, got %d", len(results))

	// Nil(true) must return only the unset record.
	results, err = client.AllTypesRepo().Query().
		Where(
			filter.AllTypes.FieldString.Equal(marker),
			filter.AllTypes.FieldIntPtr.Nil(true),
		).
		All(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Check(t, len(results) == 1, "expected 1 nil result, got %d", len(results))
	if len(results) == 1 {
		assert.Check(t, results[0].FieldIntPtr == nil)
	}
}
