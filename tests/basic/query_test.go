package basic

import (
	"context"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
	"math/rand"
	"testing"
)

const (
	randMin = 5
	randMax = 20
)

func TestQueryCount(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	if err := client.ApplySchema(ctx); err != nil {
		t.Fatal(err)
	}

	count := rand.Intn(randMax-randMin) + randMin

	for i := 0; i < count; i++ {
		err := client.AllFieldTypesRepo().Create(ctx, &model.AllFieldTypes{})
		if err != nil {
			t.Fatal(err)
		}
	}

	dbCount, err := client.AllFieldTypesRepo().Query().Count(ctx)

	if err != nil {
		t.Fatal(err)
	}

	// TODO: add database cleanup?

	assert.Equal(t, count, dbCount)
}
