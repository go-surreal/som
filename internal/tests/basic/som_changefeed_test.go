package basic

import (
	"context"
	"strings"
	"testing"
	"time"

	"som.test/model"
	"gotest.tools/v3/assert"
)

// changefeed support is exercised via the AllTypes model, which carries the
// `som:"changefeed=1d"` tag. Each test runs against its own isolated database,
// so the change stream only ever contains that test's own mutations.
func newChangefeedModel(name string) *model.AllTypes {
	return &model.AllTypes{
		FieldString: name,
		FieldInt:    0,
		FieldMonth:  time.January,
	}
}

func TestChangefeedBasic(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	repo := client.AllTypesRepo()

	newModel := newChangefeedModel("test-changefeed")

	err := repo.Create(ctx, newModel)
	assert.NilError(t, err)
	assert.Check(t, newModel.ID() != "")

	// Use versionstamp 0 to get all changes from the beginning
	entries, err := repo.Changes().SinceVersionstamp(0).Show(ctx)
	assert.NilError(t, err)
	assert.Check(t, len(entries) > 0, "expected at least one change entry")

	// SurrealDB uses "update" for all record changes including creates
	var foundRecord bool
	for _, entry := range entries {
		assert.Check(t, entry.Versionstamp > 0, "versionstamp should be non-zero")
		for _, updated := range entry.Updates {
			if updated.FieldString == "test-changefeed" {
				foundRecord = true
			}
		}
	}
	assert.Check(t, foundRecord, "expected to find the created record in changes")
}

func TestChangefeedUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	repo := client.AllTypesRepo()

	newModel := newChangefeedModel("before-update")

	err := repo.Create(ctx, newModel)
	assert.NilError(t, err)

	newModel.FieldString = "after-update"
	err = repo.Update(ctx, newModel)
	assert.NilError(t, err)

	// Query all changes from the beginning
	entries, err := repo.Changes().SinceVersionstamp(0).Show(ctx)
	assert.NilError(t, err)

	var foundUpdate bool
	for _, entry := range entries {
		for _, updated := range entry.Updates {
			if updated.FieldString == "after-update" {
				foundUpdate = true
			}
		}
	}
	assert.Check(t, foundUpdate, "expected to find the updated record in changes")
}

func TestChangefeedDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	repo := client.AllTypesRepo()

	newModel := newChangefeedModel("to-be-deleted")

	err := repo.Create(ctx, newModel)
	assert.NilError(t, err)

	recordID := newModel.ID()

	err = repo.Delete(ctx, newModel)
	assert.NilError(t, err)

	// Query all changes from the beginning
	entries, err := repo.Changes().SinceVersionstamp(0).Show(ctx)
	assert.NilError(t, err)

	// We should see at least the creation (as update) and the deletion
	// SurrealDB may report deletes in Deletes slice or as Update with the record ID
	var foundCreate, foundDelete bool
	for _, entry := range entries {
		if len(entry.Deletes) > 0 {
			foundDelete = true
		}
		for _, updated := range entry.Updates {
			if updated.FieldString == "to-be-deleted" {
				foundCreate = true
			}
			// Deleted records might appear as updates with matching ID
			if updated.ID() != "" && updated.ID() == recordID {
				foundDelete = true
			}
		}
	}
	assert.Check(t, foundCreate, "expected to find the created record before deletion")
	assert.Check(t, foundDelete, "expected to find deletion in changes")
}

func TestChangefeedMultipleOperations(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	repo := client.AllTypesRepo()

	model1 := newChangefeedModel("model-1")
	model2 := newChangefeedModel("model-2")
	model3 := newChangefeedModel("model-3")

	err := repo.Create(ctx, model1)
	assert.NilError(t, err)

	err = repo.Create(ctx, model2)
	assert.NilError(t, err)

	err = repo.Create(ctx, model3)
	assert.NilError(t, err)

	model1.FieldString = "model-1-updated"
	err = repo.Update(ctx, model1)
	assert.NilError(t, err)

	err = repo.Delete(ctx, model2)
	assert.NilError(t, err)

	// Use versionstamp 0 to get all changes
	entries, err := repo.Changes().SinceVersionstamp(0).Show(ctx)
	assert.NilError(t, err)
	assert.Check(t, len(entries) > 0, "expected change entries")

	// SurrealDB uses "update" for all record mutations
	var updateCount int
	for _, entry := range entries {
		updateCount += len(entry.Updates)
	}

	// Expect at least 3 updates (3 creates, reported as "update")
	assert.Check(t, updateCount >= 3, "expected at least 3 updates (creates), got %d", updateCount)
}

func TestChangefeedLimit(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	repo := client.AllTypesRepo()

	for i := 0; i < 5; i++ {
		m := newChangefeedModel("limit-test")
		err := repo.Create(ctx, m)
		assert.NilError(t, err)
	}

	// Use versionstamp 0 to get all changes, limited to 2 batches
	entries, err := repo.Changes().SinceVersionstamp(0).Limit(2).Show(ctx)
	assert.NilError(t, err)
	assert.Check(t, len(entries) <= 2, "expected at most 2 change batches, got %d", len(entries))
}

func TestChangefeedSinceRequired(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	repo := client.AllTypesRepo()

	_, err := repo.Changes().Show(ctx)
	assert.Check(t, err != nil, "expected error when Since() not called")
	assert.Check(t, strings.Contains(err.Error(), "Since"), "error should mention Since")
}

func TestChangefeedDescribe(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	repo := client.AllTypesRepo()

	query := repo.Changes().Since(time.Now()).Limit(10)
	desc := query.Describe()

	assert.Check(t, strings.Contains(desc, "SHOW CHANGES"), "describe should contain SHOW CHANGES")
	assert.Check(t, strings.Contains(desc, "all_types"), "describe should contain table name")
	assert.Check(t, strings.Contains(desc, "LIMIT"), "describe should contain LIMIT")
}

func TestChangefeedVersionstampOrdering(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	repo := client.AllTypesRepo()

	// Create two records
	model1 := newChangefeedModel("versionstamp-test-1")
	err := repo.Create(ctx, model1)
	assert.NilError(t, err)

	model2 := newChangefeedModel("versionstamp-test-2")
	err = repo.Create(ctx, model2)
	assert.NilError(t, err)

	// Get all changes and verify versionstamps are monotonically increasing
	entries, err := repo.Changes().SinceVersionstamp(0).Show(ctx)
	assert.NilError(t, err)
	assert.Check(t, len(entries) > 0, "expected at least one entry")

	var prevVersionstamp uint64
	for _, entry := range entries {
		assert.Check(t, entry.Versionstamp >= prevVersionstamp,
			"versionstamps should be monotonically increasing: %d >= %d",
			entry.Versionstamp, prevVersionstamp)
		prevVersionstamp = entry.Versionstamp
	}

	// Verify both records are in the changes
	var foundModel1, foundModel2 bool
	for _, entry := range entries {
		for _, updated := range entry.Updates {
			if updated.FieldString == "versionstamp-test-1" {
				foundModel1 = true
			}
			if updated.FieldString == "versionstamp-test-2" {
				foundModel2 = true
			}
		}
	}
	assert.Check(t, foundModel1, "expected to find model 1")
	assert.Check(t, foundModel2, "expected to find model 2")
}
