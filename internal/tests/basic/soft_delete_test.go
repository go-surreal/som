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

	record := model.SpecialTypes{
		Name: "Test User",
	}

	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, record.ID() != "")

	assert.Assert(t, !record.SoftDelete.IsDeleted(), "newly created record should not be deleted")
	assert.Assert(t, record.SoftDelete.DeletedAt().IsZero(), "DeletedAt should be zero for non-deleted record")

	err = client.SpecialTypesRepo().Restore(ctx, &record)
	assert.Assert(t, err != nil, "restoring a non-deleted record should return an error")

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	assert.Assert(t, record.SoftDelete.IsDeleted(), "record should be marked as deleted")
	assert.Assert(t, !record.SoftDelete.DeletedAt().IsZero(), "DeletedAt should be set after deletion")

	allRecords, err := client.SpecialTypesRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(allRecords), "deleted record should not appear in normal queries")

	allRecordsIncludingDeleted, err := client.SpecialTypesRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(allRecordsIncludingDeleted), "deleted record should appear with WithDeleted()")
	assert.Equal(t, record.Name, allRecordsIncludingDeleted[0].Name)

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.Assert(t, err != nil, "deleting an already-deleted record should return an error")

	err = client.SpecialTypesRepo().Restore(ctx, &record)
	assert.NilError(t, err)

	assert.Assert(t, !record.SoftDelete.IsDeleted(), "restored record should not be marked as deleted")
	assert.Assert(t, record.SoftDelete.DeletedAt().IsZero(), "DeletedAt should be zero after restore")

	allRecordsAfterRestore, err := client.SpecialTypesRepo().Query().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(allRecordsAfterRestore), "restored record should appear in normal queries")

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Erase(ctx, &record)
	assert.NilError(t, err)

	allRecordsAfterErase, err := client.SpecialTypesRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(allRecordsAfterErase), "erased record should not appear even with WithDeleted()")
}

func TestSoftDeleteFetchRelation(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	author := model.SpecialTypes{
		Name: "Author User",
	}
	err := client.SpecialTypesRepo().Create(ctx, &author)
	assert.NilError(t, err)

	post := model.SpecialRelation{
		Title:  "Test Post",
		Author: &author,
	}
	err = client.SpecialRelationRepo().Create(ctx, &post)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &author)
	assert.NilError(t, err)

	posts, err := client.SpecialRelationRepo().Query().
		Fetch(with.SpecialRelation.Author()).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Assert(t, posts[0].Author != nil, "fetched relations return all records regardless of soft-delete status")
	assert.Equal(t, "Author User", posts[0].Author.Name)
	assert.Assert(t, posts[0].Author.SoftDelete.IsDeleted(), "the fetched author should be marked as soft-deleted")
}

func TestSoftDeleteFetchSliceRelation(t *testing.T) {
	ctx := context.Background()
	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user1 := model.SpecialTypes{Name: "Author 1"}
	user2 := model.SpecialTypes{Name: "Author 2"}
	user3 := model.SpecialTypes{Name: "Author 3"}
	err := client.SpecialTypesRepo().Create(ctx, &user1)
	assert.NilError(t, err)
	err = client.SpecialTypesRepo().Create(ctx, &user2)
	assert.NilError(t, err)
	err = client.SpecialTypesRepo().Create(ctx, &user3)
	assert.NilError(t, err)

	post := model.SpecialRelation{
		Title:   "Test Post",
		Authors: []*model.SpecialTypes{&user1, &user2, &user3},
	}
	err = client.SpecialRelationRepo().Create(ctx, &post)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &user2)
	assert.NilError(t, err)

	posts, err := client.SpecialRelationRepo().Query().
		Fetch(with.SpecialRelation.Authors()).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, 3, len(posts[0].Authors), "all authors returned regardless of soft-delete status")

	authorNames := make([]string, len(posts[0].Authors))
	for i, a := range posts[0].Authors {
		authorNames[i] = a.Name
	}
	assert.Assert(t, contains(authorNames, "Author 1"))
	assert.Assert(t, contains(authorNames, "Author 2"), "soft-deleted Author 2 should still be returned")
	assert.Assert(t, contains(authorNames, "Author 3"))

	var activeAuthors []*model.SpecialTypes
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

func TestSoftDeleteRead(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{
		Name: "Read Test User",
	}
	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)

	id := record.ID()

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	readRecord, exists, err := client.SpecialTypesRepo().Read(ctx, id)
	assert.NilError(t, err)
	assert.Assert(t, exists, "Read should find soft-deleted records")
	assert.Assert(t, readRecord != nil)
	assert.Assert(t, readRecord.SoftDelete.IsDeleted(), "Read result should show record as deleted")
	assert.Equal(t, "Read Test User", readRecord.Name)

	err = client.SpecialTypesRepo().Erase(ctx, &record)
	assert.NilError(t, err)

	all, err := client.SpecialTypesRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(all), "erased record should not exist")
}

func TestSoftDeleteCount(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	user1 := model.SpecialTypes{Name: "User 1"}
	user2 := model.SpecialTypes{Name: "User 2"}
	user3 := model.SpecialTypes{Name: "User 3"}

	err := client.SpecialTypesRepo().Create(ctx, &user1)
	assert.NilError(t, err)
	err = client.SpecialTypesRepo().Create(ctx, &user2)
	assert.NilError(t, err)
	err = client.SpecialTypesRepo().Create(ctx, &user3)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &user2)
	assert.NilError(t, err)

	count, err := client.SpecialTypesRepo().Query().Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 2, count)

	countAll, err := client.SpecialTypesRepo().Query().WithDeleted().Count(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 3, countAll)
}

func TestSoftDeleteFirstAndExists(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{Name: "Only User"}
	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)

	exists, err := client.SpecialTypesRepo().Query().Exists(ctx)
	assert.NilError(t, err)
	assert.Assert(t, exists, "Exists should be true for non-deleted record")

	first, err := client.SpecialTypesRepo().Query().First(ctx)
	assert.NilError(t, err)
	assert.Equal(t, "Only User", first.Name)

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	exists, err = client.SpecialTypesRepo().Query().Exists(ctx)
	assert.NilError(t, err)
	assert.Assert(t, !exists, "Exists should be false after soft delete")

	existsAll, err := client.SpecialTypesRepo().Query().WithDeleted().Exists(ctx)
	assert.NilError(t, err)
	assert.Assert(t, existsAll, "WithDeleted().Exists() should be true")

	_, err = client.SpecialTypesRepo().Query().First(ctx)
	assert.Assert(t, err != nil, "First should error when no non-deleted records exist")

	firstAll, err := client.SpecialTypesRepo().Query().WithDeleted().First(ctx)
	assert.NilError(t, err)
	assert.Equal(t, "Only User", firstAll.Name)
}

func TestSoftDeleteFilter(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	alice := model.SpecialTypes{Name: "Alice"}
	bob := model.SpecialTypes{Name: "Bob"}
	err := client.SpecialTypesRepo().Create(ctx, &alice)
	assert.NilError(t, err)
	err = client.SpecialTypesRepo().Create(ctx, &bob)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &alice)
	assert.NilError(t, err)

	results, err := client.SpecialTypesRepo().Query().
		Where(filter.SpecialTypes.Name.Equal("Alice")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(results), "soft-deleted Alice should not appear in filtered query")

	resultsAll, err := client.SpecialTypesRepo().Query().
		WithDeleted().
		Where(filter.SpecialTypes.Name.Equal("Alice")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(resultsAll))
	assert.Equal(t, "Alice", resultsAll[0].Name)

	resultsBob, err := client.SpecialTypesRepo().Query().
		Where(filter.SpecialTypes.Name.Equal("Bob")).
		All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(resultsBob))
	assert.Equal(t, "Bob", resultsBob[0].Name)
}

func TestSoftDeleteErrorTypes(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{Name: "Error Test User"}
	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.Assert(t, errors.Is(err, som.ErrAlreadyDeleted), "double delete should return ErrAlreadyDeleted, got: %v", err)

	err = client.SpecialTypesRepo().Restore(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Restore(ctx, &record)
	assert.Assert(t, err != nil, "restoring a non-deleted record should error")
}

func TestSoftDeleteEraseNonDeleted(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	record := model.SpecialTypes{Name: "Erase Direct User"}
	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, !record.SoftDelete.IsDeleted())

	err = client.SpecialTypesRepo().Erase(ctx, &record)
	assert.NilError(t, err)

	all, err := client.SpecialTypesRepo().Query().WithDeleted().All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(all), "erased record should not exist even with WithDeleted")
}

func TestSoftDeleteOptimisticLockConflict(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	// -- Delete conflict --

	record := model.SpecialTypes{Name: "Lock Test"}
	err := client.SpecialTypesRepo().Create(ctx, &record)
	assert.NilError(t, err)
	assert.Equal(t, 1, record.Version())

	stale := record

	record.Name = "Lock Test Updated"
	err = client.SpecialTypesRepo().Update(ctx, &record)
	assert.NilError(t, err)
	assert.Equal(t, 2, record.Version())

	err = client.SpecialTypesRepo().Delete(ctx, &stale)
	assert.Assert(t, errors.Is(err, som.ErrOptimisticLock),
		"deleting with stale version should return ErrOptimisticLock, got: %v", err)

	// -- Restore conflict --

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, record.SoftDelete.IsDeleted())

	staleDeleted := record

	err = client.SpecialTypesRepo().Restore(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, !record.SoftDelete.IsDeleted())

	record.Name = "Lock Test Updated Again"
	err = client.SpecialTypesRepo().Update(ctx, &record)
	assert.NilError(t, err)

	err = client.SpecialTypesRepo().Delete(ctx, &record)
	assert.NilError(t, err)
	assert.Assert(t, record.SoftDelete.IsDeleted())

	err = client.SpecialTypesRepo().Restore(ctx, &staleDeleted)
	assert.Assert(t, errors.Is(err, som.ErrOptimisticLock),
		"restoring with stale version should return ErrOptimisticLock, got: %v", err)
}
