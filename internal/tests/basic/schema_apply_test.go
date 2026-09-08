package basic

import (
	"context"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"som.test/model"
)

func TestApplySchemaIsIdempotent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	allTypes := &model.AllTypes{
		FieldString:   "alpha",
		FieldFloat64:  4,
		FieldTime:     time.Now(),
		FieldDuration: time.Second,
		FieldMonth:    time.January,
	}

	err := client.AllTypesRepo().Create(ctx, allTypes)
	assert.NilError(t, err)

	for range 2 {
		err = client.ApplySchema(ctx)
		assert.NilError(t, err)
	}

	// Existing records must survive re-applying the schema.
	found, ok, err := client.AllTypesRepo().Read(ctx, string(allTypes.ID()))
	assert.NilError(t, err)
	assert.Assert(t, ok)
	assert.Equal(t, found.FieldString, "alpha")

	count, err := client.AllTypesRepo().Query().Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, count, 1)

	// Views are removed and redefined on change, so they must still resolve.
	rows, err := client.AllTypesSummaryRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, len(rows), 1)
	assert.Equal(t, rows[0].Total, 1)
}
