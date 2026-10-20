package basic

import (
	"context"
	"errors"
	"testing"
	"time"

	som "som.test/gen/som"
	"som.test/model"
	"gotest.tools/v3/assert"
)

func TestReadForUpdateRequiresTransaction(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := &model.SpecialTypes{Name: "no-tx"}
	err := client.SpecialTypesRepo().Create(ctx, node)
	assert.NilError(t, err)

	_, _, err = client.SpecialTypesRepo().ReadForUpdate(ctx, string(node.ID()))
	assert.Assert(t, errors.Is(err, som.ErrTransactionRequired), "unexpected error: %v", err)
}

func TestReadForUpdateEmptyID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	_, _, err := client.SpecialTypesRepo().ReadForUpdate(txCtx, "")
	assert.Assert(t, errors.Is(err, som.ErrEmptyID), "unexpected error: %v", err)
}

func TestReadForUpdateCommit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := &model.SpecialTypes{Name: "locked-read"}
	err := client.SpecialTypesRepo().Create(ctx, node)
	assert.NilError(t, err)

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	read, exists, err := client.SpecialTypesRepo().ReadForUpdate(txCtx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "locked-read", read.Name)
	assert.Equal(t, node.ID(), read.ID())

	err = client.SpecialTypesRepo().Create(txCtx, &model.SpecialTypes{Name: "decided-on-" + read.Name})
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)
}

func TestReadForUpdateMissingRecord(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	read, exists, err := client.SpecialTypesRepo().ReadForUpdate(txCtx, "550e8400-e29b-41d4-a716-446655440000")
	assert.NilError(t, err)
	assert.Assert(t, !exists)
	assert.Assert(t, read == nil)

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)
}

func TestReadForUpdateConflictOnConcurrentUpdate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	node := &model.SpecialTypes{Name: "before"}
	err := client.SpecialTypesRepo().Create(ctx, node)
	assert.NilError(t, err)

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	_, exists, err := client.SpecialTypesRepo().ReadForUpdate(txCtx, string(node.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)

	node.Name = "after"
	err = client.SpecialTypesRepo().Update(ctx, node)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Create(txCtx, &model.SpecialTypes{Name: "decided-on-before"})
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.Assert(t, err != nil, "commit should fail after the read record was changed concurrently")
	assert.Assert(t, som.IsTransactionConflict(err), "unexpected error: %v", err)
}

func TestReadForUpdateConflictOnConcurrentCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	const id = "550e8400-e29b-41d4-a716-446655440001"

	_, exists, err := client.SpecialTypesRepo().ReadForUpdate(txCtx, id)
	assert.NilError(t, err)
	assert.Assert(t, !exists)

	err = client.SpecialTypesRepo().CreateWithID(ctx, id, &model.SpecialTypes{Name: "concurrent"})
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Create(txCtx, &model.SpecialTypes{Name: "decided-on-missing"})
	assert.NilError(t, err)

	err = som.TxCommit(txCtx)
	assert.Assert(t, err != nil, "commit should fail after the missing record was created concurrently")
	assert.Assert(t, som.IsTransactionConflict(err), "unexpected error: %v", err)
}

func TestReadForUpdateComplexID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	key := model.WeatherKey{City: "Berlin", Date: time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)}
	err := client.WeatherRepo().CreateWithID(ctx, &model.Weather{
		Node:        som.NewNode[model.WeatherKey](key),
		Temperature: 22.5,
	})
	assert.NilError(t, err)

	txCtx, cancel := som.TxStart(ctx)
	defer cancel()

	read, exists, err := client.WeatherRepo().ReadForUpdate(txCtx, key)
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, 22.5, read.Temperature)
	assert.Equal(t, "Berlin", read.ID().City)

	err = som.TxCommit(txCtx)
	assert.NilError(t, err)
}
