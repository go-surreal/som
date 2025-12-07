package basic

import (
	"context"
	"testing"

	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestSoftDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// Create a new SoftDeleteUser record
	user := model.SoftDeleteUser{
		Name: "Test User",
	}

	err := client.SoftDeleteUserRepo().Create(ctx, &user)
	assert.NilError(t, err)
	assert.Assert(t, user.ID() != nil)

	// Verify initial state - not deleted
	assert.Assert(t, !user.SoftDelete.IsDeleted(), "newly created record should not be deleted")
	assert.Assert(t, user.SoftDelete.DeletedAt().IsZero(), "DeletedAt should be zero for non-deleted record")

	// Delete the record (soft delete)
	err = client.SoftDeleteUserRepo().Delete(ctx, &user)
	assert.NilError(t, err)

	// Refresh to get the updated state from DB
	err = client.SoftDeleteUserRepo().Refresh(ctx, &user)
	assert.NilError(t, err)

	// Verify deleted state
	assert.Assert(t, user.SoftDelete.IsDeleted(), "record should be marked as deleted")
	assert.Assert(t, !user.SoftDelete.DeletedAt().IsZero(), "DeletedAt should be set after deletion")

	// Query without WithDeleted - should not find the record
	allUsers, err := client.SoftDeleteUserRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(allUsers), "deleted record should not appear in normal queries")

	// Query with WithDeleted - should find the record
	allUsersIncludingDeleted, err := client.SoftDeleteUserRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(allUsersIncludingDeleted), "deleted record should appear with WithDeleted()")
	assert.Equal(t, user.Name, allUsersIncludingDeleted[0].Name)

	// Try to delete again - should error
	err = client.SoftDeleteUserRepo().Delete(ctx, &user)
	assert.Assert(t, err != nil, "deleting an already-deleted record should return an error")

	// Restore the record
	err = client.SoftDeleteUserRepo().Restore(ctx, &user)
	assert.NilError(t, err)

	// Verify restored state
	assert.Assert(t, !user.SoftDelete.IsDeleted(), "restored record should not be marked as deleted")
	assert.Assert(t, user.SoftDelete.DeletedAt().IsZero(), "DeletedAt should be zero after restore")

	// Query without WithDeleted - should find the record now
	allUsersAfterRestore, err := client.SoftDeleteUserRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(allUsersAfterRestore), "restored record should appear in normal queries")

	// Delete again and then Erase (hard delete)
	err = client.SoftDeleteUserRepo().Delete(ctx, &user)
	assert.NilError(t, err)

	err = client.SoftDeleteUserRepo().Erase(ctx, &user)
	assert.NilError(t, err)

	// Query with WithDeleted - should not find the record (it's gone)
	allUsersAfterErase, err := client.SoftDeleteUserRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(allUsersAfterErase), "erased record should not appear even with WithDeleted()")
}

func TestSoftDeleteWithTimestamps(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// Create a SoftDeleteComplete record (has Timestamps + OptimisticLock + SoftDelete)
	record := model.SoftDeleteComplete{
		Name: "Complete Record",
	}

	err := client.SoftDeleteCompleteRepo().Create(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, record.ID() != nil)

	// Verify Timestamps are set
	assert.Assert(t, !record.Timestamps.CreatedAt().IsZero(), "CreatedAt should be set")
	assert.Assert(t, !record.Timestamps.UpdatedAt().IsZero(), "UpdatedAt should be set")

	// Verify OptimisticLock
	assert.Equal(t, 1, record.Version())

	// Verify not deleted
	assert.Assert(t, !record.SoftDelete.IsDeleted())

	// Soft delete the record
	err = client.SoftDeleteCompleteRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	// Refresh and verify
	err = client.SoftDeleteCompleteRepo().Refresh(ctx, &record)
	assert.NilError(t, err)

	assert.Assert(t, record.SoftDelete.IsDeleted(), "record should be deleted")
	assert.Assert(t, !record.Timestamps.CreatedAt().IsZero(), "CreatedAt should still be set")
	assert.Assert(t, !record.Timestamps.UpdatedAt().IsZero(), "UpdatedAt should still be set")

	// Restore and verify all features still work
	err = client.SoftDeleteCompleteRepo().Restore(ctx, &record)
	assert.NilError(t, err)

	assert.Assert(t, !record.SoftDelete.IsDeleted(), "record should be restored")

	// Update should still work and increment version
	record.Name = "Updated Complete Record"
	err = client.SoftDeleteCompleteRepo().Update(ctx, &record)
	assert.NilError(t, err)
	assert.Equal(t, 2, record.Version())
}
