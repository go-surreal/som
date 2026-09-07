# Client API

The generated client provides the main entry point for database operations. It lives in the
generated `repo` package (`yourproject/gen/som/repo`), while shared types and errors live in
the root `som` package (`yourproject/gen/som`).

## Creating a Client

```go
import "yourproject/gen/som/repo"

client, err := repo.NewClient(ctx, repo.Config{
    Address:   "ws://localhost:8000",
    Username:  "root",
    Password:  "root",
    Namespace: "myapp",
    Database:  "production",
})
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

`NewClient` creates the namespace and database if they do not exist yet (`DEFINE NAMESPACE /
DATABASE IF NOT EXISTS`), so no manual bootstrapping is required.

## Configuration Options

```go
type Config struct {
    // Address is the SurrealDB server URL
    // Supports: ws://, wss://, http://, https://
    Address string

    // Username for authentication
    Username string

    // Password for authentication
    Password string

    // Namespace to use
    Namespace string

    // Database within the namespace
    Database string

    // ExpiryPurgeInterval controls how often expired records are purged from
    // tables with an expiry (TTL) configured. Defaults to one minute when unset.
    // It has no effect if no model embeds som.Expiry.
    ExpiryPurgeInterval time.Duration
}
```

See [Expiry (TTL)](../models/09_expiry.md) for details on the purge behaviour.

## Applying the Schema

The generated schema (tables, fields, indexes, analyzers, views) is applied with `ApplySchema`:

```go
if err := client.ApplySchema(ctx); err != nil {
    log.Fatal(err)
}
```

Call this once on startup (or as part of a deployment step) after connecting. The statements
are idempotent, so repeated calls are safe.

Tables, fields and analyzers are defined with `OVERWRITE`, so a changed definition always
reaches an existing database. Indexes and views are guarded by a hash stored in the
`__som__schema` table and are only rebuilt once their generated definition actually changes —
re-applying an unchanged schema does not reindex or recompute anything.

`ApplySchema` converges definitions, it does not migrate data. Renames, backfills and
removals of tables or fields that no longer exist in the models are not performed, since
that intent cannot be derived from the model structs.

## Version Verification

When creating a client, SOM automatically verifies that the connected SurrealDB server meets the minimum required version (currently **3.2.0**). If the version check fails, `NewClient` returns a `som.ErrUnsupportedVersion` error:

```go
client, err := repo.NewClient(ctx, config)
if err != nil {
    if errors.Is(err, som.ErrUnsupportedVersion) {
        log.Fatal("SurrealDB server version too old, please upgrade to 3.2.0+")
    }
    log.Fatal(err)
}
```

## Accessing Repositories

The client provides typed repository access for each model:

```go
// For a User model
userRepo := client.UserRepo()

// For a Post model
postRepo := client.PostRepo()
```

Repositories are generated for nodes, [views](../models/07_views.md) and
[sinks](../models/08_sinks.md). Edges have no repository of their own — they are created and
queried through the node they start from (see [Relationships](../relationships/README.md)).

## Connection Management

### Closing the Client

`Close` takes no arguments and returns nothing. It closes the connection and stops background
workers (such as the expiry purge goroutine):

```go
client.Close()
```

### Context Usage

All operations accept a context for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, exists, err := client.UserRepo().Read(ctx, id)
```

## Raw Queries

Execute arbitrary SurrealQL with parameter binding:

```go
result, err := client.Raw(ctx, "SELECT * FROM user WHERE age > $min", som.Params{"min": 18})
if err != nil {
    return err
}

users, err := result.Scan[map[string]any]()
```

See [Raw Queries](../querying/06_raw_queries.md) for full documentation.

## Error Handling

Check for errors on all operations:

```go
user, exists, err := client.UserRepo().Read(ctx, id)
if err != nil {
    // Handle connection errors, etc.
    return err
}
if !exists {
    // Record not found
}
```

### Structured Server Errors

SurrealDB v3 returns structured error responses. SOM exposes these as `som.ServerError`, which you can extract using `errors.As`:

```go
err := client.UserRepo().Update(ctx, user)
if err != nil {
    var se som.ServerError
    if errors.As(err, &se) {
        fmt.Println(se.Kind, se.Message, se.Details)
    }
}
```

SOM automatically recognizes common domain errors:

| Error | Description |
|-------|-------------|
| `som.ErrOptimisticLock` | Update failed due to version mismatch |
| `som.ErrAlreadyDeleted` | Soft delete on already-deleted record |
| `som.ErrNotFound` | Record not found |
| `som.ErrNilID` / `som.ErrEmptyID` | Operation requires a valid record ID |
| `som.ErrEmptyResponse` | Database returned an unexpected empty response |
| `som.ErrCacheNotSupported` | Caching enabled for a node with a complex ID |
| `som.ErrUnsupportedVersion` | Server version below minimum required |

### Error Kinds and Helpers

`ServerError.Kind` can be compared against the re-exported kind constants:

`som.KindValidation`, `som.KindConfiguration`, `som.KindThrown`, `som.KindQuery`,
`som.KindSerialization`, `som.KindNotAllowed`, `som.KindNotFound`, `som.KindAlreadyExists`,
`som.KindConnection`, `som.KindInternal`.

For common cases there are classification helpers that unwrap the error for you:

```go
if som.IsTransactionConflict(err) {
    // retry the transaction
}
```

Available helpers: `IsNotFound`, `IsNotAllowed`, `IsTransactionConflict`, `IsTimedOut`,
`IsNotExecuted`, `IsCancelled`, `IsParseError`, `IsDeserialization`, `IsLiveQueryNotSupported`,
`IsScriptingBlocked`, `IsTokenExpired`, `IsInvalidAuth`.

## Thread Safety

The client is safe for concurrent use from multiple goroutines.
