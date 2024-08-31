package main

import (
	"context"
	sombase "github.com/go-surreal/som"
	"github.com/go-surreal/som/tests/basic/gen/som/where"
	"github.com/go-surreal/som/tests/basic/model"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
	"testing"
	"time"
)

func TestFilterCompareFields(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareTestDatabase(ctx, t)
	defer cleanup()

	str := "Some Value"
	date := time.Now()

	modelNew := model.AllFieldTypes{
		String:    str,
		StringPtr: &str,

		Time:    date.Add(-time.Hour),
		TimePtr: &date,

		Duration: time.Hour,
	}

	modelIn := modelNew

	err := client.AllFieldTypesRepo().Create(ctx, &modelIn)
	if err != nil {
		t.Fatal(err)
	}

	modelOut, err := client.AllFieldTypesRepo().Query().
		Filter(
			where.AllFieldTypes.StringPtr.Equal_(where.AllFieldTypes.String),
			where.AllFieldTypes.TimePtr.After_(where.AllFieldTypes.Time),
		).
		First(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, str, modelOut.String)
	assert.Equal(t, str, *modelOut.StringPtr)

	assert.DeepEqual(t,
		modelNew, *modelOut,
		cmpopts.IgnoreUnexported(sombase.Node{}, sombase.Timestamps{}),
	)
}
