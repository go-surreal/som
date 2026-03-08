package basic

import (
	"context"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestAsync(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{FieldMonth: time.January})
	if err != nil {
		t.Fatal(err)
	}

	resCh := client.AllTypesRepo().Query().
		Where().
		CountAsync(ctx)

	assert.NilError(t, <-resCh.Err())
	assert.Equal(t, 1, <-resCh.Val())

	err = client.AllTypesRepo().Create(ctx, &model.AllTypes{FieldMonth: time.January})
	if err != nil {
		t.Fatal(err)
	}

	resCh = client.AllTypesRepo().Query().
		Where().
		CountAsync(ctx)

	assert.NilError(t, <-resCh.Err())
	assert.Equal(t, 2, <-resCh.Val())
}

func TestAsyncQueries(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	err := client.AllTypesRepo().Create(ctx, &model.AllTypes{
		FieldString: "async_test",
		FieldMonth:  time.January,
	})
	if err != nil {
		t.Fatal(err)
	}

	existsRes := client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldString.Equal("async_test")).
		ExistsAsync(ctx)
	assert.NilError(t, <-existsRes.Err())
	assert.Equal(t, true, <-existsRes.Val())

	firstRes := client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldString.Equal("async_test")).
		FirstAsync(ctx)
	assert.NilError(t, <-firstRes.Err())
	first := <-firstRes.Val()
	assert.Check(t, first != nil)
	assert.Equal(t, "async_test", first.FieldString)

	allRes := client.AllTypesRepo().Query().
		Where(filter.AllTypes.FieldString.Equal("async_test")).
		AllAsync(ctx)
	assert.NilError(t, <-allRes.Err())
	all := <-allRes.Val()
	assert.Equal(t, 1, len(all))
	assert.Equal(t, "async_test", all[0].FieldString)
}
