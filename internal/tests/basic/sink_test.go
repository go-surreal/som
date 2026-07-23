package basic

import (
	"context"
	"testing"

	"som.test/model"
	"gotest.tools/v3/assert"
)

// TestSinkPopulatesView verifies the sink→view ingestion pattern: records
// written to a write-only EventLog sink are discarded (DROP table) but still
// feed the EventSummary view they populate.
func TestSinkPopulatesView(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// Single create.
	err := client.EventLogRepo().Create(ctx, &model.EventLog{
		Category: "alpha",
		Value:    1,
	})
	assert.NilError(t, err)

	// Batch insert.
	err = client.EventLogRepo().Insert(ctx, []*model.EventLog{
		{Category: "alpha", Value: 3},
		{Category: "beta", Value: 10},
	})
	assert.NilError(t, err)

	rows, err := client.EventSummaryRepo().Query().All(ctx)
	assert.NilError(t, err)

	byCategory := make(map[string]*model.EventSummary, len(rows))
	for _, r := range rows {
		byCategory[r.Category] = r
	}

	assert.Equal(t, len(byCategory), 2)

	assert.Equal(t, byCategory["alpha"].Total, 2)
	assert.Equal(t, byCategory["alpha"].AvgValue, 2.0)

	assert.Equal(t, byCategory["beta"].Total, 1)
	assert.Equal(t, byCategory["beta"].AvgValue, 10.0)
}

// TestSinkNilGuards verifies the write-only repo rejects nil inputs and
// treats an empty batch as a no-op.
func TestSinkNilGuards(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	assert.Error(t, client.EventLogRepo().Create(ctx, nil), "the passed record must not be nil")
	assert.Error(t, client.EventLogRepo().Insert(ctx, []*model.EventLog{nil}), "slice contains nil record")
	assert.NilError(t, client.EventLogRepo().Insert(ctx, nil))
}
