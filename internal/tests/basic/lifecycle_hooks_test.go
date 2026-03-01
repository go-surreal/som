package basic

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-surreal/som/tests/basic/gen/som/filter"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestHookBeforeCreate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	unregister := client.SpecialTypesRepo().OnBeforeCreate(func(ctx context.Context, group *model.SpecialTypes) error {
		group.Name = "modified-by-hook"
		return nil
	})
	defer unregister()

	group := model.SpecialTypes{Name: "original"}
	err := client.SpecialTypesRepo().Create(ctx, &group)
	assert.NilError(t, err)

	read, exists, err := client.SpecialTypesRepo().Read(ctx, string(group.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "modified-by-hook", read.Name)
}

func TestHookAfterCreate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var called atomic.Bool

	unregister := client.SpecialTypesRepo().OnAfterCreate(func(ctx context.Context, group *model.SpecialTypes) error {
		called.Store(true)
		return nil
	})
	defer unregister()

	group := model.SpecialTypes{Name: "test"}
	err := client.SpecialTypesRepo().Create(ctx, &group)
	assert.NilError(t, err)
	assert.Assert(t, called.Load())

	_, exists, err := client.SpecialTypesRepo().Read(ctx, string(group.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
}

func TestHookBeforeCreateAbort(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	hookErr := errors.New("abort create")

	unregister := client.SpecialTypesRepo().OnBeforeCreate(func(ctx context.Context, group *model.SpecialTypes) error {
		return hookErr
	})
	defer unregister()

	group := model.SpecialTypes{Name: "should-not-exist"}
	err := client.SpecialTypesRepo().Create(ctx, &group)
	assert.Assert(t, errors.Is(err, hookErr))

	results, err := client.SpecialTypesRepo().Query().Where(
		filter.SpecialTypes.Name.Equal("should-not-exist"),
	).All(ctx)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(results))
}

func TestHookBeforeUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := model.SpecialTypes{Name: "original"}
	err := client.SpecialTypesRepo().Create(ctx, &group)
	assert.NilError(t, err)

	unregister := client.SpecialTypesRepo().OnBeforeUpdate(func(ctx context.Context, group *model.SpecialTypes) error {
		group.Name = "modified-by-update-hook"
		return nil
	})
	defer unregister()

	group.Name = "updated"
	err = client.SpecialTypesRepo().Update(ctx, &group)
	assert.NilError(t, err)

	read, exists, err := client.SpecialTypesRepo().Read(ctx, string(group.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "modified-by-update-hook", read.Name)
}

func TestHookAfterUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var called atomic.Bool

	group := model.SpecialTypes{Name: "test"}
	err := client.SpecialTypesRepo().Create(ctx, &group)
	assert.NilError(t, err)

	unregister := client.SpecialTypesRepo().OnAfterUpdate(func(ctx context.Context, group *model.SpecialTypes) error {
		called.Store(true)
		return nil
	})
	defer unregister()

	group.Name = "updated"
	err = client.SpecialTypesRepo().Update(ctx, &group)
	assert.NilError(t, err)
	assert.Assert(t, called.Load())
}

func TestHookBeforeDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := model.SpecialTypes{Name: "test"}
	err := client.SpecialTypesRepo().Create(ctx, &group)
	assert.NilError(t, err)

	hookErr := errors.New("abort delete")

	unregister := client.SpecialTypesRepo().OnBeforeDelete(func(ctx context.Context, group *model.SpecialTypes) error {
		return hookErr
	})
	defer unregister()

	err = client.SpecialTypesRepo().Delete(ctx, &group)
	assert.Assert(t, errors.Is(err, hookErr))

	_, exists, err := client.SpecialTypesRepo().Read(ctx, string(group.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
}

func TestHookAfterDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var called atomic.Bool

	group := model.SpecialTypes{Name: "test"}
	err := client.SpecialTypesRepo().Create(ctx, &group)
	assert.NilError(t, err)
	id := group.ID()

	unregister := client.SpecialTypesRepo().OnAfterDelete(func(ctx context.Context, group *model.SpecialTypes) error {
		called.Store(true)
		return nil
	})
	defer unregister()

	err = client.SpecialTypesRepo().Delete(ctx, &group)
	assert.NilError(t, err)
	assert.Assert(t, called.Load())

	read, exists, err := client.SpecialTypesRepo().Read(ctx, string(id))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Assert(t, read.SoftDelete.IsDeleted())
}

func TestHookCleanup(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var called atomic.Bool

	unregister := client.SpecialTypesRepo().OnBeforeCreate(func(ctx context.Context, group *model.SpecialTypes) error {
		called.Store(true)
		return nil
	})

	unregister()

	group := model.SpecialTypes{Name: "test"}
	err := client.SpecialTypesRepo().Create(ctx, &group)
	assert.NilError(t, err)
	assert.Assert(t, !called.Load())
}

func TestHookMultiple(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var order []int

	unregister1 := client.SpecialTypesRepo().OnBeforeCreate(func(ctx context.Context, group *model.SpecialTypes) error {
		order = append(order, 1)
		return nil
	})
	defer unregister1()

	unregister2 := client.SpecialTypesRepo().OnBeforeCreate(func(ctx context.Context, group *model.SpecialTypes) error {
		order = append(order, 2)
		return nil
	})
	defer unregister2()

	group := model.SpecialTypes{Name: "test"}
	err := client.SpecialTypesRepo().Create(ctx, &group)
	assert.NilError(t, err)
	assert.Equal(t, 2, len(order))
	assert.Equal(t, 1, order[0])
	assert.Equal(t, 2, order[1])
}

func TestModelHookBeforeCreate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.AllTypes{FieldHookStatus: "active", FieldMonth: time.January}
	err := client.AllTypesRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	read, exists, err := client.AllTypesRepo().Read(ctx, string(rec.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "[created]active", read.FieldHookStatus)
}

func TestModelHookAfterCreate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.AllTypes{FieldHookDetail: "info", FieldMonth: time.January}
	err := client.AllTypesRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	assert.Assert(t, strings.HasSuffix(rec.FieldHookDetail, "[after-create]"))
}

func TestModelHookBeforeUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.AllTypes{FieldHookStatus: "init", FieldMonth: time.January}
	err := client.AllTypesRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	rec.FieldHookStatus = "changed"
	err = client.AllTypesRepo().Update(ctx, &rec)
	assert.NilError(t, err)

	read, exists, err := client.AllTypesRepo().Read(ctx, string(rec.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
	assert.Equal(t, "[updated]changed", read.FieldHookStatus)
}

func TestModelHookAfterUpdate(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.AllTypes{FieldHookDetail: "info", FieldMonth: time.January}
	err := client.AllTypesRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	err = client.AllTypesRepo().Update(ctx, &rec)
	assert.NilError(t, err)

	assert.Assert(t, strings.HasSuffix(rec.FieldHookDetail, "[after-update]"))
}

func TestModelHookBeforeDeleteAbort(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.AllTypes{FieldHookStatus: "keep", FieldMonth: time.January}
	err := client.AllTypesRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	abort := true
	deleteCtx := context.WithValue(ctx, model.AbortDeleteKey, &abort)

	err = client.AllTypesRepo().Delete(deleteCtx, &rec)
	assert.Assert(t, err != nil)

	_, exists, err := client.AllTypesRepo().Read(ctx, string(rec.ID()))
	assert.NilError(t, err)
	assert.Assert(t, exists)
}

func TestModelHookAfterDelete(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	rec := model.AllTypes{FieldHookStatus: "remove", FieldMonth: time.January}
	err := client.AllTypesRepo().Create(ctx, &rec)
	assert.NilError(t, err)
	id := string(rec.ID())

	called := false
	deleteCtx := context.WithValue(ctx, model.AfterDeleteCalledKey, &called)

	err = client.AllTypesRepo().Delete(deleteCtx, &rec)
	assert.NilError(t, err)
	assert.Assert(t, called)

	_, exists, err := client.AllTypesRepo().Read(ctx, id)
	assert.NilError(t, err)
	assert.Assert(t, !exists)
}

func TestModelHookAndRepoHookOrder(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	var repoSawStatus string
	unregister := client.AllTypesRepo().OnBeforeCreate(func(_ context.Context, rec *model.AllTypes) error {
		repoSawStatus = rec.FieldHookStatus
		return nil
	})
	defer unregister()

	rec := model.AllTypes{FieldHookStatus: "orig", FieldMonth: time.January}
	err := client.AllTypesRepo().Create(ctx, &rec)
	assert.NilError(t, err)

	assert.Equal(t, "[created]orig", repoSawStatus)
}
