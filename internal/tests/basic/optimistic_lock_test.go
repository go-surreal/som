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

	// Create a new Group record
	group := model.Group{
		Name: "Test Group",
	}

	err := client.GroupRepo().Create(ctx, &group)
	assert.NilError(t, err)
	assert.Assert(t, group.ID() != nil)

	// Verify initial version is 1
	assert.Equal(t, 1, group.Version())

	// Update the record - should succeed and increment version
	group.Name = "Updated Group"
	err = client.GroupRepo().Update(ctx, &group)
	assert.NilError(t, err)
	assert.Equal(t, 2, group.Version())

	// Read a fresh copy of the same record (simulating another process)
	staleGroup, exists, err := client.GroupRepo().Read(ctx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, 2, staleGroup.Version())

	// Update the original copy - version becomes 3
	group.Name = "Updated Again"
	err = client.GroupRepo().Update(ctx, &group)
	assert.NilError(t, err)
	assert.Equal(t, 3, group.Version())

	// Try to update the stale copy (still has version 2)
	// This should fail with ErrOptimisticLock
	staleGroup.Name = "Stale Update"
	err = client.GroupRepo().Update(ctx, staleGroup)
	assert.Assert(t, err != nil, "expected error from stale update")
	assert.Assert(t, errors.Is(err, som.ErrOptimisticLock),
		"expected ErrOptimisticLock, got: %v", err)

	// Verify the record was not updated with stale data
	finalGroup, exists, err := client.GroupRepo().Read(ctx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "Updated Again", finalGroup.Name)
	assert.Equal(t, 3, finalGroup.Version())
}
