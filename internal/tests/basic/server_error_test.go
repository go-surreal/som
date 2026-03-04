package basic

import (
	"context"
	"errors"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestServerError_OptimisticLockUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{Name: "Test Record"}
	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)

	stale, exists, err := client.SpecialTypesRepo().Read(ctx, string(record.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)

	record.Name = "Updated"
	err = client.SpecialTypesRepo().Update(ctx, &record)
	assert.NilError(t, err)

	stale.Name = "Stale Update"
	err = client.SpecialTypesRepo().Update(ctx, stale)
	assert.Assert(t, err != nil)
	assert.Assert(t, errors.Is(err, som.ErrOptimisticLock))

	var se som.ServerError
	assert.Assert(t, errors.As(err, &se),
		"optimistic lock error from Update (RPC) should contain a ServerError, got: %v", err)
	assert.Assert(t, se.Message != "", "ServerError.Message should not be empty")
}

func TestServerError_OptimisticLockSoftDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{Name: "Test Record"}
	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)

	stale := record

	record.Name = "Updated"
	err = client.SpecialTypesRepo().Update(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &stale)
	assert.Assert(t, err != nil)
	assert.Assert(t, errors.Is(err, som.ErrOptimisticLock))
}

func TestServerError_AlreadyDeleted(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{Name: "Test Record"}
	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.Assert(t, err != nil)
	assert.Assert(t, errors.Is(err, som.ErrAlreadyDeleted))
}

func TestServerError_OptimisticLockRestore(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{Name: "Test Record"}
	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	staleDeleted := record

	err = client.SpecialTypesRepo().Restore(ctx, &record)
	assert.NilError(t, err)

	record.Name = "Updated After Restore"
	err = client.SpecialTypesRepo().Update(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Restore(ctx, &staleDeleted)
	assert.Assert(t, err != nil)
	assert.Assert(t, errors.Is(err, som.ErrOptimisticLock))
}
