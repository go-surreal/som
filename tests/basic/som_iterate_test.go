package basic

import (
	"context"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
	"testing"
)

func TestQueryIterate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	if err := client.ApplySchema(ctx); err != nil {
		t.Fatal(err)
	}

	query := client.AllFieldTypesRepo().Query()

	//for range query.Iterate(ctx, 10) {
	//	t.Fatal("there should be no results to iterate over")
	//}

	wanted := 5

	for range wanted {
		if err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{}); err != nil {
			t.Fatal(err)
		}
	}

	got := 0

	for _, err := range query.Iterate(ctx, 10) {
		if err != nil {
			t.Fatal(err)
		}

		got++
	}

	assert.Equal(t, wanted, got)

	got = 0

	for _, err := range query.Iterate(ctx, 2) {
		if err != nil {
			t.Fatal(err)
		}

		got++
	}

	assert.Equal(t, wanted, got)
}
