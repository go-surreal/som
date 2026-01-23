# Caching

SOM provides optional request-scoped caching for repository Read operations. The cache is stored in the context and is completely opt-in.

## Overview

- Cache is stored in context (request-scoped)
- Two modes: lazy (on-demand) and eager (pre-load all)
- Thread-safe via `sync.RWMutex`
- Each model has isolated cache (caching Group records doesn't affect User reads)
- No automatic invalidation on Update/Delete

## Methods

### WithCache

Creates an empty lazy cache. Records are fetched from the database on first access and stored in the cache. Subsequent reads for the same ID return the cached pointer.

```go
func (r *UserRepo) WithCache(ctx context.Context) context.Context
```

**Behavior:**
- First `Read()` for an ID: queries database, stores result in cache
- Subsequent `Read()` for same ID: returns cached pointer (no database query)
- Records created/updated after cache creation: fetched from database and cached on first access

### WithCacheAll

Pre-loads all records from the table into an eager cache. This is useful when you know you'll need most or all records.

```go
func (r *UserRepo) WithCacheAll(ctx context.Context) (context.Context, error)
```

**Behavior:**
- Executes `Query().All()` to load all records
- Returns error if query fails
- `Read()` calls only check the cache (no database queries)
- Records not in cache return `(nil, false, nil)` without querying database
- Records created after cache load are not visible through the cached context

### DropCache

Removes the cache for this model from the context. Subsequent reads will query the database directly.

```go
func (r *UserRepo) DropCache(ctx context.Context) context.Context
```

## Usage Examples

### Lazy Cache

```go
// Create an empty lazy cache
ctx = client.UserRepo().WithCache(ctx)

// First read: queries database, populates cache
user1, exists, err := client.UserRepo().Read(ctx, id1)

// Second read with same ID: returns cached pointer (no DB query)
user1Again, _, _ := client.UserRepo().Read(ctx, id1)

// user1 == user1Again (same pointer)
```

### Eager Cache

```go
// Pre-load all users into cache
ctx, err := client.UserRepo().WithCacheAll(ctx)
if err != nil {
    return err
}

// All reads come from cache (no DB queries)
user1, exists, _ := client.UserRepo().Read(ctx, id1)
user2, exists, _ := client.UserRepo().Read(ctx, id2)

// Reading a non-existent or newly-created ID
newUser := &model.User{Name: "New"}
client.UserRepo().Create(ctx, newUser)

// Returns (nil, false, nil) - eager cache doesn't query DB
cached, exists, _ := client.UserRepo().Read(ctx, newUser.ID())
// exists == false
```

### Dropping Cache

```go
// Enable cache
ctx = client.UserRepo().WithCache(ctx)
user, _, _ := client.UserRepo().Read(ctx, id) // cached

// Drop cache to get fresh data
ctx = client.UserRepo().DropCache(ctx)
user, _, _ := client.UserRepo().Read(ctx, id) // queries DB
```

### Multiple Models

Caches are isolated per model type:

```go
// Cache for users only
ctx = client.UserRepo().WithCache(ctx)

// User reads use cache
user, _, _ := client.UserRepo().Read(ctx, userID) // cached after first read

// Post reads still query database (no cache for Post)
post, _, _ := client.PostRepo().Read(ctx, postID) // always queries DB

// Add cache for posts too
ctx = client.PostRepo().WithCache(ctx)

// Now both use caching
user, _, _ := client.UserRepo().Read(ctx, userID) // uses User cache
post, _, _ := client.PostRepo().Read(ctx, postID) // uses Post cache
```

## When to Use

### Lazy Cache (`WithCache`)

Best for:
- Request handlers that may read the same record multiple times
- Graph traversals that might revisit nodes
- Operations where you don't know which records will be accessed

### Eager Cache (`WithCacheAll`)

Best for:
- Batch operations that need most/all records from a table
- Small reference tables (roles, categories, settings)
- Reports or exports that iterate over all records

### When Not to Use

- Long-running processes (cache may become stale)
- Write-heavy operations (cache doesn't auto-invalidate)
- Tables with many records and eager caching (memory concerns)

## Important Notes

1. **No Automatic Invalidation**: If you `Update` or `Delete` a record, the cache is not automatically updated. Either drop the cache or use a fresh context for subsequent reads.

2. **Per-Model Isolation**: Each model type has its own cache. Caching `User` records has no effect on `Post` reads.

3. **Context Scoped**: The cache lives in the context. Different contexts have different caches (or no cache).

4. **Thread Safety**: The cache uses `sync.RWMutex` and is safe for concurrent reads.

5. **Pointer Identity**: Repeated reads of the same ID return the same pointer, which can be useful for equality checks but means mutations affect all references.
