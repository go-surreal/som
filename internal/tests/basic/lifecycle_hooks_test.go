package basic

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestHookBeforeCreate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	unregister := client.GroupRepo().OnBeforeCreate(func(ctx context.Context, group *model.Group) error {
		group.Name = "modified-by-hook"
		return nil
	})
	defer unregister()

	group := model.Group{Name: "original"}
	err := client.GroupRepo().Create(ctx, &group)
	assert.NilError(t, err)

	read, exists, err := client.GroupRepo().Read(ctx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "modified-by-hook", read.Name)
}

func TestHookAfterCreate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var called atomic.Bool

	unregister := client.GroupRepo().OnAfterCreate(func(ctx context.Context, group *model.Group) error {
		called.Store(true)
		return nil
	})
	defer unregister()

	group := model.Group{Name: "test"}
	err := client.GroupRepo().Create(ctx, &group)
	assert.NilError(t, err)
	assert.Assert(t, called.Load())

	_, exists, err := client.GroupRepo().Read(ctx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
}

func TestHookBeforeCreateAbort(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	hookErr := errors.New("abort create")

	unregister := client.GroupRepo().OnBeforeCreate(func(ctx context.Context, group *model.Group) error {
		return hookErr
	})
	defer unregister()

	group := model.Group{Name: "should-not-exist"}
	err := client.GroupRepo().Create(ctx, &group)
	assert.Assert(t, errors.Is(err, hookErr))

	results, err := client.GroupRepo().Query().Where(
		filter.Group.Name.Equal("should-not-exist"),
	).All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(results))
}

func TestHookBeforeUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := model.Group{Name: "original"}
	err := client.GroupRepo().Create(ctx, &group)
	assert.NilError(t, err)

	unregister := client.GroupRepo().OnBeforeUpdate(func(ctx context.Context, group *model.Group) error {
		group.Name = "modified-by-update-hook"
		return nil
	})
	defer unregister()

	group.Name = "updated"
	err = client.GroupRepo().Update(ctx, &group)
	assert.NilError(t, err)

	read, exists, err := client.GroupRepo().Read(ctx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "modified-by-update-hook", read.Name)
}

func TestHookAfterUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var called atomic.Bool

	group := model.Group{Name: "test"}
	err := client.GroupRepo().Create(ctx, &group)
	assert.NilError(t, err)

	unregister := client.GroupRepo().OnAfterUpdate(func(ctx context.Context, group *model.Group) error {
		called.Store(true)
		return nil
	})
	defer unregister()

	group.Name = "updated"
	err = client.GroupRepo().Update(ctx, &group)
	assert.NilError(t, err)
	assert.Assert(t, called.Load())
}

func TestHookBeforeDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := model.Group{Name: "test"}
	err := client.GroupRepo().Create(ctx, &group)
	assert.NilError(t, err)

	hookErr := errors.New("abort delete")

	unregister := client.GroupRepo().OnBeforeDelete(func(ctx context.Context, group *model.Group) error {
		return hookErr
	})
	defer unregister()

	err = client.GroupRepo().Delete(ctx, &group)
	assert.Assert(t, errors.Is(err, hookErr))

	_, exists, err := client.GroupRepo().Read(ctx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
}

func TestHookAfterDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var called atomic.Bool

	group := model.Group{Name: "test"}
	err := client.GroupRepo().Create(ctx, &group)
	assert.NilError(t, err)
	id := group.ID()

	unregister := client.GroupRepo().OnAfterDelete(func(ctx context.Context, group *model.Group) error {
		called.Store(true)
		return nil
	})
	defer unregister()

	err = client.GroupRepo().Delete(ctx, &group)
	assert.NilError(t, err)
	assert.Assert(t, called.Load())

	_, exists, err := client.GroupRepo().Read(ctx, id)
	assert.NilError(t, err)
	assert.Assert(t, !exists)
}

func TestHookCleanup(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var called atomic.Bool

	unregister := client.GroupRepo().OnBeforeCreate(func(ctx context.Context, group *model.Group) error {
		called.Store(true)
		return nil
	})

	unregister()

	group := model.Group{Name: "test"}
	err := client.GroupRepo().Create(ctx, &group)
	assert.NilError(t, err)
	assert.Assert(t, !called.Load())
}

func TestHookMultiple(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var order []int

	unregister1 := client.GroupRepo().OnBeforeCreate(func(ctx context.Context, group *model.Group) error {
		order = append(order, 1)
		return nil
	})
	defer unregister1()

	unregister2 := client.GroupRepo().OnBeforeCreate(func(ctx context.Context, group *model.Group) error {
		order = append(order, 2)
		return nil
	})
	defer unregister2()

	group := model.Group{Name: "test"}
	err := client.GroupRepo().Create(ctx, &group)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(order))
	assert.Equal(t, 1, order[0])
	assert.Equal(t, 2, order[1])
}

func TestModelHookBeforeCreate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.FieldsLikeDBResponse{Status: "active"}
	err := client.FieldsLikeDBResponseRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	read, exists, err := client.FieldsLikeDBResponseRepo().Read(ctx, rec.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "[created]active", read.Status)
}

func TestModelHookAfterCreate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.FieldsLikeDBResponse{Detail: "info"}
	err := client.FieldsLikeDBResponseRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	assert.Assert(t, strings.HasSuffix(rec.Detail, "[after-create]"))
}

func TestModelHookBeforeUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.FieldsLikeDBResponse{Status: "init"}
	err := client.FieldsLikeDBResponseRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	rec.Status = "changed"
	err = client.FieldsLikeDBResponseRepo().Update(ctx, &rec)
	assert.NilError(t, err)

	read, exists, err := client.FieldsLikeDBResponseRepo().Read(ctx, rec.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "[updated]changed", read.Status)
}

func TestModelHookAfterUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.FieldsLikeDBResponse{Detail: "info"}
	err := client.FieldsLikeDBResponseRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	err = client.FieldsLikeDBResponseRepo().Update(ctx, &rec)
	assert.NilError(t, err)

	assert.Assert(t, strings.HasSuffix(rec.Detail, "[after-update]"))
}

func TestModelHookBeforeDeleteAbort(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.FieldsLikeDBResponse{Status: "keep"}
	err := client.FieldsLikeDBResponseRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	abort := true
	deleteCtx := context.WithValue(ctx, model.AbortDeleteKey, &abort)

	err = client.FieldsLikeDBResponseRepo().Delete(deleteCtx, &rec)
	assert.Assert(t, err != nil)

	_, exists, err := client.FieldsLikeDBResponseRepo().Read(ctx, rec.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists)
}

func TestModelHookAfterDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.FieldsLikeDBResponse{Status: "remove"}
	err := client.FieldsLikeDBResponseRepo().Create(ctx, &rec)
	assert.NilError(t, err)
	id := rec.ID()

	called := false
	deleteCtx := context.WithValue(ctx, model.AfterDeleteCalledKey, &called)

	err = client.FieldsLikeDBResponseRepo().Delete(deleteCtx, &rec)
	assert.NilError(t, err)
	assert.Assert(t, called)

	_, exists, err := client.FieldsLikeDBResponseRepo().Read(ctx, id)
	assert.NilError(t, err)
	assert.Assert(t, !exists)
}

func TestModelHookAndRepoHookOrder(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var repoSawStatus string
	unregister := client.FieldsLikeDBResponseRepo().OnBeforeCreate(func(_ context.Context, rec *model.FieldsLikeDBResponse) error {
		repoSawStatus = rec.Status
		return nil
	})
	defer unregister()

	rec := model.FieldsLikeDBResponse{Status: "orig"}
	err := client.FieldsLikeDBResponseRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	assert.Equal(t, "[created]orig", repoSawStatus)
}
