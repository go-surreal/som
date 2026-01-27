# Caching

SOM provides optional request-scoped caching for repository Read operations. The cache is completely opt-in and requires explicit cleanup.

## Overview

- Request-scoped caching for Read operations
- Two modes: lazy (on-demand) and eager (pre-load all)
- Thread-safe via `sync.RWMutex`
- Each model type has isolated cache
- Cleanup function must be called when done (typically via `defer`)
- No automatic invalidation on Update/Delete

## API Reference

### WithCache

Creates a cache for the specified model type. Returns a context with caching enabled and a cleanup function that must be called.

```go
func WithCache[T Model](ctx context.Context, opts ...CacheOption) (context.Context, func())
```

The cleanup function removes the cache from the global store and marks it as cleaned. After cleanup, any Read using that context returns `ErrCacheAlreadyCleaned`.

### Options

Configure cache behavior using functional options:

```go
som.Lazy()              // Lazy loading (default) - fetch on demand
som.Eager()             // Eager loading - load all records on first Read
som.WithTTL(duration)   // Entry expiration time
som.WithMaxSize(n)      // Max records for eager cache (default 1000)
```

### Errors

```go
som.ErrCacheSizeLimitExceeded  // Eager cache table has more records than MaxSize
som.ErrCacheAlreadyCleaned     // Using cache context after cleanup() was called
```

## Usage Examples

### Basic Lazy Cache

```go
// Enable lazy cache for Group model
ctx, cleanup := som.WithCache[model.Group](ctx)
defer cleanup()

// First read: queries database, populates cache
group1, exists, err := client.GroupRepo().Read(ctx, id1)

// Second read with same ID: returns cached pointer (no DB query)
group1Again, _, _ := client.GroupRepo().Read(ctx, id1)

// group1 == group1Again (same pointer)
```

### Explicit Lazy Mode

```go
ctx, cleanup := som.WithCache[model.Group](ctx, som.Lazy())
defer cleanup()

// Same behavior as default
```

### Eager Cache

```go
ctx, cleanup := som.WithCache[model.Group](ctx, som.Eager())
defer cleanup()

// First Read triggers loading ALL records from the table
group1, exists, err := client.GroupRepo().Read(ctx, group1ID)

// All subsequent reads come from cache only (no DB queries)
group2, exists, _ := client.GroupRepo().Read(ctx, group2ID)
group3, exists, _ := client.GroupRepo().Read(ctx, group3ID)

// Records created after cache load are not visible
newGroup := &model.Group{Name: "New"}
client.GroupRepo().Create(ctx, newGroup)

// Returns (nil, false, nil) - eager cache doesn't query DB for misses
cached, exists, _ := client.GroupRepo().Read(ctx, newGroup.ID())
// exists == false
```

### Eager Cache with MaxSize

```go
// Limit eager cache to 5000 records
ctx, cleanup := som.WithCache[model.Group](ctx, som.Eager(), som.WithMaxSize(5000))
defer cleanup()

// Returns ErrCacheSizeLimitExceeded if table has > 5000 records
group, exists, err := client.GroupRepo().Read(ctx, id)
if errors.Is(err, som.ErrCacheSizeLimitExceeded) {
    // Handle: table too large for eager cache
}
```

### Cache with TTL (Lazy Mode)

```go
ctx, cleanup := som.WithCache[model.Group](ctx, som.WithTTL(5*time.Minute))
defer cleanup()

// First read caches the record
group1, _, _ := client.GroupRepo().Read(ctx, id)

// Within TTL: returns cached pointer
group2, _, _ := client.GroupRepo().Read(ctx, id)
// group1 == group2

// After TTL expires: fetches fresh data from DB
time.Sleep(6 * time.Minute)
group3, _, _ := client.GroupRepo().Read(ctx, id)
// group1 != group3 (different pointer, fresh data)
```

### Eager Cache with TTL (Auto-Refresh)

```go
ctx, cleanup := som.WithCache[model.Group](ctx, som.Eager(), som.WithTTL(5*time.Minute))
defer cleanup()

// First Read loads ALL records into cache
group1, exists, _ := client.GroupRepo().Read(ctx, group1ID)

// New record created after cache load
newGroup := &model.Group{Name: "New"}
client.GroupRepo().Create(context.Background(), newGroup)

// Not visible yet (cache hasn't expired)
cached, exists, _ := client.GroupRepo().Read(ctx, newGroup.ID())
// exists == false

// After TTL expires, next Read triggers full cache refresh
time.Sleep(6 * time.Minute)
refreshed, exists, _ := client.GroupRepo().Read(ctx, newGroup.ID())
// exists == true (cache was refreshed, new record is now visible)
```

### Multiple Models

Caches are isolated per model type:

```go
// Cache for Group only
ctx, cleanupGroup := som.WithCache[model.Group](ctx)
defer cleanupGroup()

// Group reads use cache
group, _, _ := client.GroupRepo().Read(ctx, groupID) // cached after first read

// User reads still query database (no cache for User)
user, _, _ := client.UserRepo().Read(ctx, userID) // always queries DB

// Add cache for User too
ctx, cleanupUser := som.WithCache[model.User](ctx)
defer cleanupUser()

// Now both use caching
group, _, _ := client.GroupRepo().Read(ctx, groupID) // uses Group cache
user, _, _ := client.UserRepo().Read(ctx, userID)    // uses User cache
```

### Creating New Cache After Cleanup

```go
ctx, cleanup := som.WithCache[model.Group](ctx)

group1, _, _ := client.GroupRepo().Read(ctx, id)

cleanup() // Cache cleaned up

// Update the record
group1.Name = "Updated"
client.GroupRepo().Update(context.Background(), group1)

// Create a fresh cache to see updated data
ctx, newCleanup := som.WithCache[model.Group](context.Background())
defer newCleanup()

group2, _, _ := client.GroupRepo().Read(ctx, id)
// group2.Name == "Updated" (fresh data)
```

### Error After Cleanup

```go
ctx, cleanup := som.WithCache[model.Group](ctx)

group, _, _ := client.GroupRepo().Read(ctx, id) // works

cleanup()

_, _, err := client.GroupRepo().Read(ctx, id)
// err == som.ErrCacheAlreadyCleaned
```

## Behavior Details

### Lazy Mode (Default)

1. First `Read()` for an ID queries the database and stores the result in cache
2. Subsequent `Read()` calls for the same ID return the cached pointer
3. Cache misses always query the database
4. Records created/updated after cache creation are fetched and cached on first access

### Eager Mode

1. First `Read()` call loads ALL records from the table into cache
2. Before loading, checks record count against MaxSize (default 1000)
3. If count exceeds MaxSize, returns `ErrCacheSizeLimitExceeded`
4. After initial load, all reads only check the cache
5. Cache misses return `(nil, false, nil)` without querying the database
6. Records created after cache load are not visible through the cached context
7. With TTL enabled, the entire cache expires and refreshes automatically on next access
8. During TTL refresh, MaxSize is re-checked; returns `ErrCacheSizeLimitExceeded` if exceeded

### Cleanup Lifecycle

1. `WithCache` generates a unique cache ID and stores it in the context
2. Actual cache data is stored in a global map keyed by this ID
3. Calling cleanup removes the cache from the global map and marks the ID as dropped
4. Any subsequent Read using that context checks if the ID is dropped and returns `ErrCacheAlreadyCleaned`

## When to Use

### Lazy Cache

Best for:
- Request handlers that may read the same record multiple times
- Graph traversals that might revisit nodes
- Operations where you don't know which records will be accessed

### Eager Cache

Best for:
- Batch operations that need most/all records from a table
- Small reference tables (roles, categories, settings)
- Reports or exports that iterate over all records

### When Not to Use

- Long-running processes (cache may become stale)
- Write-heavy operations (cache doesn't auto-invalidate)
- Tables with many records and eager caching (memory concerns)

## Important Notes

1. **Cleanup is Required**: Always call the cleanup function, typically via `defer`. Failing to do so leaves stale entries in the global cache store.

2. **No Automatic Invalidation**: If you `Update` or `Delete` a record, the cache is not automatically updated. Create a new cache after writes to see fresh data.

3. **Per-Model Isolation**: Each model type has its own cache. Caching `Group` records has no effect on `User` reads.

4. **Context Carries ID, Not Data**: The context stores a cache ID and options. The actual cache data is stored in a global map, enabling cleanup from anyfilter.

5. **Thread Safety**: The cache uses `sync.RWMutex` and is safe for concurrent reads and writes.

6. **Pointer Identity**: Repeated reads of the same ID return the same pointer, which can be useful for equality checks but means mutations affect all references.
