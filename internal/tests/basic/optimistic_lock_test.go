package basic

import (
	"context"
	"errors"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestOptimisticLock(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{
		Name: "Test Record",
	}

	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, record.ID() != "")

	assert.Equal(t, 1, record.Version())

	record.Name = "Updated Record"
	err = client.SpecialTypesRepo().Update(ctx, &record)
	assert.NilError(t, err)
	assert.Equal(t, 2, record.Version())

	staleRecord, exists, err := client.SpecialTypesRepo().Read(ctx, record.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, 2, staleRecord.Version())

	record.Name = "Updated Again"
	err = client.SpecialTypesRepo().Update(ctx, &record)
	assert.NilError(t, err)
	assert.Equal(t, 3, record.Version())

	staleRecord.Name = "Stale Update"
	err = client.SpecialTypesRepo().Update(ctx, staleRecord)
	assert.Assert(t, err != nil, "expected error from stale update")
	assert.Assert(t, errors.Is(err, som.ErrOptimisticLock),
		"expected ErrOptimisticLock, got: %v", err)

	finalRecord, exists, err := client.SpecialTypesRepo().Read(ctx, record.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "Updated Again", finalRecord.Name)
	assert.Equal(t, 3, finalRecord.Version())
}
