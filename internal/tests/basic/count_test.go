package basic

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

const (
	randMin = 5
	randMax = 20
)

func TestQueryCount(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	count := rand.Intn(randMax-randMin) + randMin

	for i := 0; i < count; i++ {
		err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
			FieldTime:     time.Now(),
			FieldDuration: time.Second,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	dbCount, err := client.AllTypesRepo().Query().Count(ctx)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, count, dbCount)
}
