package basic

import (
	"context"
	"testing"
	"time"

	"som.test/model"
	"gotest.tools/v3/assert"
)

func TestQueryView(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	data := []struct {
		category string
		value    float64
	}{
		{"alpha", 1},
		{"alpha", 3},
		{"beta", 10},
	}

	for _, d := range data {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldString:   d.category,
			FieldFloat64:  d.value,
			FieldTime:     time.Now(),
			FieldDuration: time.Second,
			FieldMonth:    time.January,
		})
		assert.NilError(t, err)
	}

	rows, err := client.AllTypesSummaryRepo().Query().All(ctx)
	assert.NilError(t, err)

	byCategory := make(map[string]*model.AllTypesSummary, len(rows))
	for _, r := range rows {
		byCategory[r.Category] = r
	}

	assert.Equal(t, len(byCategory), 2)

	assert.Equal(t, byCategory["alpha"].Total, 2)
	assert.Equal(t, byCategory["alpha"].AvgValue, 2.0)

	assert.Equal(t, byCategory["beta"].Total, 1)
	assert.Equal(t, byCategory["beta"].AvgValue, 10.0)
}
