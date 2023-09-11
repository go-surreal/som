package basic

import (
	"context"
	"github.com/go-surreal/som/examples/basic/model"
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
		err := client.UserRepo().Create(ctx, &model.User{})
		if err != nil {
			t.Fatal(err)
		}
	}

	dbCount, err := client.UserRepo().Query().Count(ctx)

	if err != nil {
		t.Fatal(err)
	}

	// TODO: add database cleanup?

	assert.Equal(t, count, dbCount)
}
