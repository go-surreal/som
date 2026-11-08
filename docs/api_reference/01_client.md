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
DATABASE IF NOT EXISTS`), so no manual bootstrapping is required. For a database som owns,
set `Strict: true` — see [Strict Mode](#strict-mode).

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

    // Strict marks the database as strict on creation.
    Strict bool

    // Reconnect controls how a lost WebSocket connection is recovered.
    // The zero value reconnects automatically with the driver defaults.
    Reconnect ReconnectConfig

    // OnStatus, if set, is called with each change of the connection status.
    OnStatus func(ConnectionEvent)
}
```

See [Expiry (TTL)](../models/09_expiry.md) for details on the purge behaviour and
[Reconnecting](#reconnecting) for the connection settings.

## Strict Mode

`Strict: true` appends `STRICT` to the `DEFINE DATABASE` statement, which makes SurrealDB
refuse any resource that has not been defined first — a `CREATE` against an unknown table
fails with `The table '...' does not exist` instead of creating it on the fly.

This is the preferred way to set up a database for som. som generates a complete schema for
every model, so nothing it queries is ever undefined, and strict mode turns the schema into
an enforced contract: a stray table name can no longer appear in the database by accident.

```go
client, err := repo.NewClient(ctx, repo.Config{
    Address:   "ws://localhost:8000",
    Username:  "root",
    Password:  "root",
    Namespace: "myapp",
    Database:  "production",
    Strict:    true,
})
```

Two limits are worth knowing:

- **It only applies to a database som creates.** The statement is
  `DEFINE DATABASE IF NOT EXISTS <db> STRICT`, so pointing a strict config at a database that
  already exists does nothing — it keeps the mode it was created with. To convert an existing
  database, run `DEFINE DATABASE OVERWRITE <db> STRICT` yourself, after checking what that
  does to your data.
- **It affects the whole database, not just som.** If other services or hand-written SurrealQL
  share the same database and rely on implicit table creation, they will break too. Only
  enable it where som owns the database.

Strict mode is a per-database setting as of SurrealDB 3.0. The older server-level
`--strict` / `SURREAL_STRICT` flag is ignored by current versions, which log
`Strict mode is no longer defined on the server level` on startup.

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

When creating a client, SOM automatically verifies that the connected SurrealDB server is
supported: at least **4.0.0**, and below the first version the underlying driver does not
support (currently **5.0.0**). If the version check fails, `NewClient` returns a
`som.ErrUnsupportedVersion` error:

```go
client, err := repo.NewClient(ctx, config)
if err != nil {
    if errors.Is(err, som.ErrUnsupportedVersion) {
        log.Fatal("SurrealDB server version too old, please upgrade to 4.0.0+")
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

### Reconnecting

A WebSocket client reconnects on its own when the connection is lost. Before any other request
is sent, the namespace, database and authentication are restored on the new connection;
requests made in the meantime wait for it until their context ends.

Authentication is restored with the token of the original sign-in. If that token has expired by
then (after one hour for system users by default), the client signs in again with the configured
`Username` and `Password`. Requests sent in the brief moment before that sign-in completes fail
with a permission error.

What cannot be restored is lost:

- **Open transactions** are gone. Retry them.
- **Live queries** end with a `LiveDisconnected` result. Subscribe again once the client is
  connected — see [Live Queries](../querying/04_live_queries.md).
- **Requests in flight** when the connection dropped fail. They may or may not have run on the
  server, so only retry them blindly if they are idempotent.

The behaviour is tuned via `Config.Reconnect`; the zero value keeps the defaults:

```go
client, err := repo.NewClient(ctx, repo.Config{
    // ...
    Reconnect: repo.ReconnectConfig{
        MinDelay:    500 * time.Millisecond, // base of the exponential backoff (default 1s)
        MaxDelay:    30 * time.Second,       // cap of the backoff (default 1m)
        MaxAttempts: -1,                     // attempts before giving up (default 5, negative = no limit)
    },
})
```

Set `Disabled: true` to turn reconnecting off, so that a lost connection closes the client for
good. HTTP connections are stateless and unaffected.

To log or react to connection changes, pass `OnStatus`. It is called with each change — the
connection being lost, every attempt to reconnect and its outcome, and the client closing:

```go
client, err := repo.NewClient(ctx, repo.Config{
    // ...
    OnStatus: func(event repo.ConnectionEvent) {
        slog.Info("surrealdb connection", "status", event.Status, "attempt", event.Attempt, "error", event.Err)
    },
})
```

`event.Status` is one of `repo.StatusConnected`, `repo.StatusDisconnected`,
`repo.StatusReconnecting` or `repo.StatusClosed`. On `StatusConnected` after a reconnect,
`event.Err` holds what could not be restored, or `nil` once the session is fully back — including
when the client had to sign in again. A `StatusClosed` event after failed
attempts means the client gave up and has to be recreated.

### Health Check

`Health` reports whether the SurrealDB server is up. It needs no authentication, which makes
it suitable for readiness probes:

```go
if err := client.Health(ctx); err != nil {
    // server not reachable or not healthy
}
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
| `som.ErrIndexBuilding` | Index rebuilt while its build is still running |
| `som.ErrUnsupportedVersion` | Server version outside the supported range |

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
