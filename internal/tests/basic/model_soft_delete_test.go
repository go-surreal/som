package basic

import (
	"context"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/with"
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

	// Try to restore a non-deleted record - should error
	err = client.SoftDeleteUserRepo().Restore(ctx, &user)
	assert.Assert(t, err != nil, "restoring a non-deleted record should return an error")

	// Delete the record (soft delete) - auto-refreshes in-memory object
	err = client.SoftDeleteUserRepo().Delete(ctx, &user)
	assert.NilError(t, err)

	// Verify deleted state (no manual Refresh needed - Delete auto-refreshes)
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

func TestSoftDeleteFetchRelation(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// Create a SoftDeleteUser (author)
	author := model.SoftDeleteUser{
		Name: "Author User",
	}
	err := client.SoftDeleteUserRepo().Create(ctx, &author)
	assert.NilError(t, err)

	// Create a SoftDeletePost with the author
	post := model.SoftDeletePost{
		Title:  "Test Post",
		Author: &author,
	}
	err = client.SoftDeletePostRepo().Create(ctx, &post)
	assert.NilError(t, err)

	// Soft delete the author
	err = client.SoftDeleteUserRepo().Delete(ctx, &author)
	assert.NilError(t, err)

	// Query the post with Fetch(Author) - soft-delete filtering does NOT apply to fetched relations
	// All related records are returned regardless of their soft-delete status
	posts, err := client.SoftDeletePostRepo().Query().
		Fetch(with.SoftDeletePost.Author()).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(posts))
	// The author is still returned even though it's soft-deleted
	assert.Assert(t, posts[0].Author != nil, "fetched relations return all records regardless of soft-delete status")
	assert.Equal(t, "Author User", posts[0].Author.Name)
	// We can check if it's soft-deleted using IsDeleted()
	assert.Assert(t, posts[0].Author.SoftDelete.IsDeleted(), "the fetched author should be marked as soft-deleted")
}

func TestSoftDeleteFetchSliceRelation(t *testing.T) {
	ctx := context.Background()
	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// Create multiple SoftDeleteUsers
	user1 := model.SoftDeleteUser{Name: "Author 1"}
	user2 := model.SoftDeleteUser{Name: "Author 2"}
	user3 := model.SoftDeleteUser{Name: "Author 3"}
	err := client.SoftDeleteUserRepo().Create(ctx, &user1)
	assert.NilError(t, err)
	err = client.SoftDeleteUserRepo().Create(ctx, &user2)
	assert.NilError(t, err)
	err = client.SoftDeleteUserRepo().Create(ctx, &user3)
	assert.NilError(t, err)

	// Create BlogPost with all authors
	post := model.SoftDeleteBlogPost{
		Title:   "Test Post",
		Authors: []*model.SoftDeleteUser{&user1, &user2, &user3},
	}
	err = client.SoftDeleteBlogPostRepo().Create(ctx, &post)
	assert.NilError(t, err)

	// Soft delete user2
	err = client.SoftDeleteUserRepo().Delete(ctx, &user2)
	assert.NilError(t, err)

	// Fetch BlogPost with Authors - soft-delete filtering does NOT apply to fetched relations
	// All related records are returned regardless of their soft-delete status
	posts, err := client.SoftDeleteBlogPostRepo().Query().
		Fetch(with.SoftDeleteBlogPost.Authors()).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, 3, len(posts[0].Authors), "all authors returned regardless of soft-delete status")

	// Verify all authors are present
	authorNames := make([]string, len(posts[0].Authors))
	for i, a := range posts[0].Authors {
		authorNames[i] = a.Name
	}
	assert.Assert(t, contains(authorNames, "Author 1"))
	assert.Assert(t, contains(authorNames, "Author 2"), "soft-deleted Author 2 should still be returned")
	assert.Assert(t, contains(authorNames, "Author 3"))

	// Users can filter soft-deleted records themselves using IsDeleted()
	var activeAuthors []*model.SoftDeleteUser
	for _, author := range posts[0].Authors {
		if !author.SoftDelete.IsDeleted() {
			activeAuthors = append(activeAuthors, author)
		}
	}
	assert.Equal(t, 2, len(activeAuthors), "manual filtering by IsDeleted() should work")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
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
	// Version is now 4 because:
	// 1. Create: version 1
	// 2. Delete (soft delete UPDATE): version 2
	// 3. Restore (UPDATE SET deleted_at = NONE): version 3
	// 4. Update: version 4
	record.Name = "Updated Complete Record"
	err = client.SoftDeleteCompleteRepo().Update(ctx, &record)
	assert.NilError(t, err)
	assert.Equal(t, 4, record.Version())
}
