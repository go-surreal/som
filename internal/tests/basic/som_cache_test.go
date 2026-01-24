package basic

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	som "github.com/go-surreal/som/tests/basic/gen/som"
	"github.com/go-surreal/som/tests/basic/model"
	"gotest.tools/v3/assert"
)

func TestCacheLazy(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := &model.Group{Name: "Test Group"}
	err := client.GroupRepo().Create(ctx, group)
	assert.NilError(t, err)

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx)
	defer cacheCleanup()

	read1, exists1, err := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists1)
	assert.Equal(t, "Test Group", read1.Name)

	read2, exists2, err := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists2)
	assert.Equal(t, "Test Group", read2.Name)

	assert.Assert(t, read1 == read2, "expected same pointer from cache")
}

func TestCacheLazyExplicit(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := &model.Group{Name: "Test Group"}
	err := client.GroupRepo().Create(ctx, group)
	assert.NilError(t, err)

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx, som.Lazy())
	defer cacheCleanup()

	read1, exists1, err := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists1)
	assert.Equal(t, "Test Group", read1.Name)

	read2, exists2, err := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists2)
	assert.Equal(t, "Test Group", read2.Name)

	assert.Assert(t, read1 == read2, "expected same pointer from cache")
}

func TestCacheEager(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group1 := &model.Group{Name: "Group 1"}
	group2 := &model.Group{Name: "Group 2"}
	group3 := &model.Group{Name: "Group 3"}

	for _, g := range []*model.Group{group1, group2, group3} {
		err := client.GroupRepo().Create(ctx, g)
		assert.NilError(t, err)
	}

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx, som.Eager())
	defer cacheCleanup()

	read1, exists1, err := client.GroupRepo().Read(cachedCtx, group1.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists1)
	assert.Equal(t, "Group 1", read1.Name)

	read2, exists2, err := client.GroupRepo().Read(cachedCtx, group2.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists2)
	assert.Equal(t, "Group 2", read2.Name)

	read3, exists3, err := client.GroupRepo().Read(cachedCtx, group3.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists3)
	assert.Equal(t, "Group 3", read3.Name)

	group4 := &model.Group{Name: "Group 4"}
	err = client.GroupRepo().Create(ctx, group4)
	assert.NilError(t, err)

	read4, exists4, err := client.GroupRepo().Read(cachedCtx, group4.ID())
	assert.NilError(t, err)
	assert.Assert(t, !exists4, "eager cache should return false for record created after cache load")
	assert.Assert(t, read4 == nil, "eager cache should return nil for missing record")
}

func TestCacheEagerWithMaxSize(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	for i := 0; i < 5; i++ {
		group := &model.Group{Name: fmt.Sprintf("Group %d", i)}
		err := client.GroupRepo().Create(ctx, group)
		assert.NilError(t, err)
	}

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx, som.Eager(), som.WithMaxSize(3))
	defer cacheCleanup()

	_, _, err := client.GroupRepo().Read(cachedCtx, som.MakeID("group", "test"))
	assert.ErrorIs(t, err, som.ErrCacheSizeLimitExceeded)
}

func TestCacheCleanup(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := &model.Group{Name: "Test Group"}
	err := client.GroupRepo().Create(ctx, group)
	assert.NilError(t, err)

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx)

	read1, exists1, err := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists1)
	assert.Equal(t, "Test Group", read1.Name)

	cacheCleanup()

	_, _, err = client.GroupRepo().Read(cachedCtx, group.ID())
	assert.ErrorIs(t, err, som.ErrCacheAlreadyCleaned)
}

func TestCacheCleanupThenNewCache(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := &model.Group{Name: "Test Group"}
	err := client.GroupRepo().Create(ctx, group)
	assert.NilError(t, err)

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx)

	read1, exists1, err := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists1)
	assert.Equal(t, "Test Group", read1.Name)

	cacheCleanup()

	group.Name = "Updated Group"
	err = client.GroupRepo().Update(ctx, group)
	assert.NilError(t, err)

	newCachedCtx, newCacheCleanup := som.WithCache[model.Group](ctx)
	defer newCacheCleanup()

	read2, exists2, err := client.GroupRepo().Read(newCachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists2)
	assert.Equal(t, "Updated Group", read2.Name, "should get fresh data with new cache")

	assert.Assert(t, read1 != read2, "expected different pointers with new cache")
}

func TestCacheIsolation(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := &model.Group{Name: "Test Group"}
	err := client.GroupRepo().Create(ctx, group)
	assert.NilError(t, err)

	allFieldTypes := &model.AllFieldTypes{String: "Test"}
	err = client.AllFieldTypesRepo().Create(ctx, allFieldTypes)
	assert.NilError(t, err)

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx)
	defer cacheCleanup()

	readGroup, existsGroup, err := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, existsGroup)
	assert.Equal(t, "Test Group", readGroup.Name)

	readAFT, existsAFT, err := client.AllFieldTypesRepo().Read(cachedCtx, allFieldTypes.ID())
	assert.NilError(t, err)
	assert.Assert(t, existsAFT)
	assert.Equal(t, "Test", readAFT.String)

	allFieldTypes.String = "Updated"
	err = client.AllFieldTypesRepo().Update(ctx, allFieldTypes)
	assert.NilError(t, err)

	readAFT2, existsAFT2, err := client.AllFieldTypesRepo().Read(cachedCtx, allFieldTypes.ID())
	assert.NilError(t, err)
	assert.Assert(t, existsAFT2)
	assert.Equal(t, "Updated", readAFT2.String, "AllFieldTypes should not be affected by Group cache")

	readGroup2, _, _ := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.Assert(t, readGroup == readGroup2, "Group cache should still return cached pointer")
}

func TestCacheConcurrent(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	const numGroups = 10
	groups := make([]*model.Group, numGroups)
	for i := 0; i < numGroups; i++ {
		groups[i] = &model.Group{Name: "Group " + string(rune('A'+i))}
		err := client.GroupRepo().Create(ctx, groups[i])
		assert.NilError(t, err)
	}

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx, som.Eager())
	defer cacheCleanup()

	const numGoroutines = 50
	const readsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	errCh := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < readsPerGoroutine; j++ {
				groupIdx := (workerID + j) % numGroups
				read, exists, err := client.GroupRepo().Read(cachedCtx, groups[groupIdx].ID())
				if err != nil {
					errCh <- err
					return
				}
				if !exists {
					errCh <- fmt.Errorf("expected group to exist (worker %d, groupIdx %d)", workerID, groupIdx)
					return
				}
				_ = read.Name
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		assert.NilError(t, err)
	}
}

func TestCacheLazyPopulation(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group1 := &model.Group{Name: "Group 1"}
	group2 := &model.Group{Name: "Group 2"}
	err := client.GroupRepo().Create(ctx, group1)
	assert.NilError(t, err)
	err = client.GroupRepo().Create(ctx, group2)
	assert.NilError(t, err)

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx)
	defer cacheCleanup()

	read1a, _, _ := client.GroupRepo().Read(cachedCtx, group1.ID())
	read1b, _, _ := client.GroupRepo().Read(cachedCtx, group1.ID())
	assert.Assert(t, read1a == read1b, "same ID should return same cached pointer")

	read2a, _, _ := client.GroupRepo().Read(cachedCtx, group2.ID())
	assert.Assert(t, read1a != read2a, "different IDs should return different pointers")

	read2b, _, _ := client.GroupRepo().Read(cachedCtx, group2.ID())
	assert.Assert(t, read2a == read2b, "same ID should return same cached pointer")
}

func TestCacheWithTTL(t *testing.T) {
	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	group := &model.Group{Name: "Test Group"}
	err := client.GroupRepo().Create(ctx, group)
	assert.NilError(t, err)

	cachedCtx, cacheCleanup := som.WithCache[model.Group](ctx, som.WithTTL(100*time.Millisecond))
	defer cacheCleanup()

	read1, exists1, err := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists1)
	assert.Equal(t, "Test Group", read1.Name)

	read2, _, _ := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.Assert(t, read1 == read2, "should return cached pointer before TTL expires")

	time.Sleep(300 * time.Millisecond)

	read3, exists3, err := client.GroupRepo().Read(cachedCtx, group.ID())
	assert.NilError(t, err)
	assert.Assert(t, exists3)
	assert.Assert(t, read1 != read3, "should return fresh data after TTL expires")
}
