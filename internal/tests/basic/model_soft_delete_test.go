package basic

import (
	"context"
	"errors"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/gen/som/with"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestSoftDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// Create a new SoftDeleteUser record
	user := model.SoftDeleteNode{
		Name: "Test User",
	}

	err := client.SoftDeleteNodeRepo().Create(ctx, &user)
	assert.NilError(t, err)
	assert.Assert(t, user.ID() != nil)

	// Verify initial state - not deleted
	assert.Assert(t, !user.SoftDelete.IsDeleted(), "newly created record should not be deleted")
	assert.Assert(t, user.SoftDelete.DeletedAt().IsZero(), "DeletedAt should be zero for non-deleted record")

	// Try to restore a non-deleted record - should error
	err = client.SoftDeleteNodeRepo().Restore(ctx, &user)
	assert.Assert(t, err != nil, "restoring a non-deleted record should return an error")

	// Delete the record (soft delete) - auto-refreshes in-memory object
	err = client.SoftDeleteNodeRepo().Delete(ctx, &user)
	assert.NilError(t, err)

	// Verify deleted state (no manual Refresh needed - Delete auto-refreshes)
	assert.Assert(t, user.SoftDelete.IsDeleted(), "record should be marked as deleted")
	assert.Assert(t, !user.SoftDelete.DeletedAt().IsZero(), "DeletedAt should be set after deletion")

	// Query without WithDeleted - should not find the record
	allUsers, err := client.SoftDeleteNodeRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(allUsers), "deleted record should not appear in normal queries")

	// Query with WithDeleted - should find the record
	allUsersIncludingDeleted, err := client.SoftDeleteNodeRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(allUsersIncludingDeleted), "deleted record should appear with WithDeleted()")
	assert.Equal(t, user.Name, allUsersIncludingDeleted[0].Name)

	// Try to delete again - should error
	err = client.SoftDeleteNodeRepo().Delete(ctx, &user)
	assert.Assert(t, err != nil, "deleting an already-deleted record should return an error")

	// Restore the record
	err = client.SoftDeleteNodeRepo().Restore(ctx, &user)
	assert.NilError(t, err)

	// Verify restored state
	assert.Assert(t, !user.SoftDelete.IsDeleted(), "restored record should not be marked as deleted")
	assert.Assert(t, user.SoftDelete.DeletedAt().IsZero(), "DeletedAt should be zero after restore")

	// Query without WithDeleted - should find the record now
	allUsersAfterRestore, err := client.SoftDeleteNodeRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(allUsersAfterRestore), "restored record should appear in normal queries")

	// Delete again and then Erase (hard delete)
	err = client.SoftDeleteNodeRepo().Delete(ctx, &user)
	assert.NilError(t, err)

	err = client.SoftDeleteNodeRepo().Erase(ctx, &user)
	assert.NilError(t, err)

	// Query with WithDeleted - should not find the record (it's gone)
	allUsersAfterErase, err := client.SoftDeleteNodeRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(allUsersAfterErase), "erased record should not appear even with WithDeleted()")
}

func TestSoftDeleteFetchRelation(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// Create a SoftDeleteUser (author)
	author := model.SoftDeleteNode{
		Name: "Author User",
	}
	err := client.SoftDeleteNodeRepo().Create(ctx, &author)
	assert.NilError(t, err)

	// Create a SoftDeletePost with the author
	post := model.SoftDeletePost{
		Title:  "Test Post",
		Author: &author,
	}
	err = client.SoftDeletePostRepo().Create(ctx, &post)
	assert.NilError(t, err)

	// Soft delete the author
	err = client.SoftDeleteNodeRepo().Delete(ctx, &author)
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
	user1 := model.SoftDeleteNode{Name: "Author 1"}
	user2 := model.SoftDeleteNode{Name: "Author 2"}
	user3 := model.SoftDeleteNode{Name: "Author 3"}
	err := client.SoftDeleteNodeRepo().Create(ctx, &user1)
	assert.NilError(t, err)
	err = client.SoftDeleteNodeRepo().Create(ctx, &user2)
	assert.NilError(t, err)
	err = client.SoftDeleteNodeRepo().Create(ctx, &user3)
	assert.NilError(t, err)

	// Create BlogPost with all authors
	post := model.SoftDeleteBlogPost{
		Title:   "Test Post",
		Authors: []*model.SoftDeleteNode{&user1, &user2, &user3},
	}
	err = client.SoftDeleteBlogPostRepo().Create(ctx, &post)
	assert.NilError(t, err)

	// Soft delete user2
	err = client.SoftDeleteNodeRepo().Delete(ctx, &user2)
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
	var activeAuthors []*model.SoftDeleteNode
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
	record := model.SoftDeleteNode{
		Name: "Complete Record",
	}

	err := client.SoftDeleteNodeRepo().Create(ctx, &record)
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
	err = client.SoftDeleteNodeRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	// Refresh and verify
	err = client.SoftDeleteNodeRepo().Refresh(ctx, &record)
	assert.NilError(t, err)

	assert.Assert(t, record.SoftDelete.IsDeleted(), "record should be deleted")
	assert.Assert(t, !record.Timestamps.CreatedAt().IsZero(), "CreatedAt should still be set")
	assert.Assert(t, !record.Timestamps.UpdatedAt().IsZero(), "UpdatedAt should still be set")

	// Restore and verify all features still work
	err = client.SoftDeleteNodeRepo().Restore(ctx, &record)
	assert.NilError(t, err)

	assert.Assert(t, !record.SoftDelete.IsDeleted(), "record should be restored")

	// Update should still work and increment version
	// Version is now 4 because:
	// 1. Create: version 1
	// 2. Delete (soft delete UPDATE): version 2
	// 3. Restore (UPDATE SET deleted_at = NONE): version 3
	// 4. Update: version 4
	record.Name = "Updated Complete Record"
	err = client.SoftDeleteNodeRepo().Update(ctx, &record)
	assert.NilError(t, err)
	assert.Equal(t, 4, record.Version())
}

func TestSoftDeleteRead(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user := model.SoftDeleteNode{
		Name: "Read Test User",
	}
	err := client.SoftDeleteNodeRepo().Create(ctx, &user)
	assert.NilError(t, err)

	id := user.ID()

	// Soft delete the user
	err = client.SoftDeleteNodeRepo().Delete(ctx, &user)
	assert.NilError(t, err)

	// Read by ID should still return the soft-deleted record
	readUser, exists, err := client.SoftDeleteNodeRepo().Read(ctx, id)
	assert.NilError(t, err)
	assert.Assert(t, exists, "Read should find soft-deleted records")
	assert.Assert(t, readUser != nil)
	assert.Assert(t, readUser.SoftDelete.IsDeleted(), "Read result should show record as deleted")
	assert.Equal(t, "Read Test User", readUser.Name)

	// Erase the record
	err = client.SoftDeleteNodeRepo().Erase(ctx, &user)
	assert.NilError(t, err)

	// Verify via query that the record is gone
	all, err := client.SoftDeleteNodeRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(all), "erased record should not exist")
}

func TestSoftDeleteCount(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user1 := model.SoftDeleteNode{Name: "User 1"}
	user2 := model.SoftDeleteNode{Name: "User 2"}
	user3 := model.SoftDeleteNode{Name: "User 3"}

	err := client.SoftDeleteNodeRepo().Create(ctx, &user1)
	assert.NilError(t, err)
	err = client.SoftDeleteNodeRepo().Create(ctx, &user2)
	assert.NilError(t, err)
	err = client.SoftDeleteNodeRepo().Create(ctx, &user3)
	assert.NilError(t, err)

	// Soft delete one user
	err = client.SoftDeleteNodeRepo().Delete(ctx, &user2)
	assert.NilError(t, err)

	// Count without WithDeleted should exclude soft-deleted
	count, err := client.SoftDeleteNodeRepo().Query().Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 2, count)

	// Count with WithDeleted should include all
	countAll, err := client.SoftDeleteNodeRepo().Query().WithDeleted().Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 3, countAll)
}

func TestSoftDeleteFirstAndExists(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user := model.SoftDeleteNode{Name: "Only User"}
	err := client.SoftDeleteNodeRepo().Create(ctx, &user)
	assert.NilError(t, err)

	// Exists should be true before deletion
	exists, err := client.SoftDeleteNodeRepo().Query().Exists(ctx)
	assert.NilError(t, err)
	assert.Assert(t, exists, "Exists should be true for non-deleted record")

	// First should return the record
	first, err := client.SoftDeleteNodeRepo().Query().First(ctx)
	assert.NilError(t, err)
	assert.Equal(t, "Only User", first.Name)

	// Soft delete
	err = client.SoftDeleteNodeRepo().Delete(ctx, &user)
	assert.NilError(t, err)

	// Exists should be false after soft delete
	exists, err = client.SoftDeleteNodeRepo().Query().Exists(ctx)
	assert.NilError(t, err)
	assert.Assert(t, !exists, "Exists should be false after soft delete")

	// WithDeleted Exists should still be true
	existsAll, err := client.SoftDeleteNodeRepo().Query().WithDeleted().Exists(ctx)
	assert.NilError(t, err)
	assert.Assert(t, existsAll, "WithDeleted().Exists() should be true")

	// First should return error after soft delete (empty result)
	_, err = client.SoftDeleteNodeRepo().Query().First(ctx)
	assert.Assert(t, err != nil, "First should error when no non-deleted records exist")

	// WithDeleted First should return the record
	firstAll, err := client.SoftDeleteNodeRepo().Query().WithDeleted().First(ctx)
	assert.NilError(t, err)
	assert.Equal(t, "Only User", firstAll.Name)
}

func TestSoftDeleteFilter(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	alice := model.SoftDeleteNode{Name: "Alice"}
	bob := model.SoftDeleteNode{Name: "Bob"}
	err := client.SoftDeleteNodeRepo().Create(ctx, &alice)
	assert.NilError(t, err)
	err = client.SoftDeleteNodeRepo().Create(ctx, &bob)
	assert.NilError(t, err)

	// Soft delete Alice
	err = client.SoftDeleteNodeRepo().Delete(ctx, &alice)
	assert.NilError(t, err)

	// Filter for Alice without WithDeleted - should be empty
	results, err := client.SoftDeleteNodeRepo().Query().
		Where(filter.SoftDeleteNode.Name.Equal("Alice")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(results), "soft-deleted Alice should not appear in filtered query")

	// Filter for Alice with WithDeleted - should return Alice
	resultsAll, err := client.SoftDeleteNodeRepo().Query().
		WithDeleted().
		Where(filter.SoftDeleteNode.Name.Equal("Alice")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(resultsAll))
	assert.Equal(t, "Alice", resultsAll[0].Name)

	// Filter for Bob without WithDeleted - should return Bob
	resultsBob, err := client.SoftDeleteNodeRepo().Query().
		Where(filter.SoftDeleteNode.Name.Equal("Bob")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(resultsBob))
	assert.Equal(t, "Bob", resultsBob[0].Name)
}

func TestSoftDeleteErrorTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user := model.SoftDeleteNode{Name: "Error Test User"}
	err := client.SoftDeleteNodeRepo().Create(ctx, &user)
	assert.NilError(t, err)

	// Delete the user
	err = client.SoftDeleteNodeRepo().Delete(ctx, &user)
	assert.NilError(t, err)

	// Double delete should return ErrAlreadyDeleted
	err = client.SoftDeleteNodeRepo().Delete(ctx, &user)
	assert.Assert(t, errors.Is(err, som.ErrAlreadyDeleted), "double delete should return ErrAlreadyDeleted, got: %v", err)

	// Restore, then try restoring a non-deleted record
	err = client.SoftDeleteNodeRepo().Restore(ctx, &user)
	assert.NilError(t, err)

	err = client.SoftDeleteNodeRepo().Restore(ctx, &user)
	assert.Assert(t, err != nil, "restoring a non-deleted record should error")
}

func TestSoftDeleteEraseNonDeleted(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user := model.SoftDeleteNode{Name: "Erase Direct User"}
	err := client.SoftDeleteNodeRepo().Create(ctx, &user)
	assert.NilError(t, err)
	assert.Assert(t, !user.SoftDelete.IsDeleted())

	// Erase without soft-deleting first should succeed (hard delete)
	err = client.SoftDeleteNodeRepo().Erase(ctx, &user)
	assert.NilError(t, err)

	// Verify record is gone even with WithDeleted
	all, err := client.SoftDeleteNodeRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(all), "erased record should not exist even with WithDeleted")
}

func TestSoftDeleteOptimisticLockConflict(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// -- Delete conflict --

	record := model.SoftDeleteNode{Name: "Lock Test"}
	err := client.SoftDeleteNodeRepo().Create(ctx, &record)
	assert.NilError(t, err)
	assert.Equal(t, 1, record.Version())

	// Take a stale copy at version 1
	stale := record

	// Update the original to bump version to 2
	record.Name = "Lock Test Updated"
	err = client.SoftDeleteNodeRepo().Update(ctx, &record)
	assert.NilError(t, err)
	assert.Equal(t, 2, record.Version())

	// Try to delete stale copy (sends version 1, DB has version 2)
	err = client.SoftDeleteNodeRepo().Delete(ctx, &stale)
	assert.Assert(t, errors.Is(err, som.ErrOptimisticLock),
		"deleting with stale version should return ErrOptimisticLock, got: %v", err)

	// -- Restore conflict --

	// Delete the original (now at version 2 → becomes version 3 after delete)
	err = client.SoftDeleteNodeRepo().Delete(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, record.SoftDelete.IsDeleted())

	// Take a stale deleted copy
	staleDeleted := record

	// Restore the original
	err = client.SoftDeleteNodeRepo().Restore(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, !record.SoftDelete.IsDeleted())

	// Update the original to bump version further
	record.Name = "Lock Test Updated Again"
	err = client.SoftDeleteNodeRepo().Update(ctx, &record)
	assert.NilError(t, err)

	// Delete the original again
	err = client.SoftDeleteNodeRepo().Delete(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, record.SoftDelete.IsDeleted())

	// Try to restore the stale deleted copy (old version, DB has newer version)
	err = client.SoftDeleteNodeRepo().Restore(ctx, &staleDeleted)
	assert.Assert(t, errors.Is(err, som.ErrOptimisticLock),
		"restoring with stale version should return ErrOptimisticLock, got: %v", err)
}
