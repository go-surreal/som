# Expiry (TTL)

Embedding `som.Expiry` gives a node a time-to-live: every record gets an `expires_at` timestamp,
expired records are filtered out of queries, and a background worker purges them.

## Enabling Expiry

Embed `som.Expiry` and declare the lifetime as a bare duration in the `som` tag:

```go
type Session struct {
    som.Node[som.ULID]
    som.Expiry `som:"24h"`

    Token string
}
```

The duration is a SurrealDB duration literal — combinations of `ns`, `us`/`µs`, `ms`, `s`, `m`,
`h`, `d`, `w`, `y` (e.g. `90s`, `24h`, `7d`, `1y`). The tag is required; omitting it is a
generation error.

Expiry can be combined with `som.Timestamps`, `som.OptimisticLock` and `som.SoftDelete`.
It is **not** supported on edges, views or sinks.

## Generated Schema

```surql
DEFINE FIELD expires_at ON TABLE session TYPE option<datetime>
    VALUE $before OR (time::now() + 24h) READONLY;

DEFINE INDEX __som__session_expires_at ON session FIELDS expires_at CONCURRENTLY;
```

`expires_at` is set by the database on creation and is read-only — it does not shift on update.
The index keeps the purge deletes efficient.

## Reading the Expiry

```go
session := &model.Session{Token: "abc"}
err := client.SessionRepo().Create(ctx, session)

fmt.Println(session.Expiry.ExpiresAt())  // populated after create
```

## Filter on Read

Expired records are excluded from queries automatically, even before they are physically
removed:

```go
// Only records whose expires_at is still in the future
active, err := client.SessionRepo().Query().All(ctx)

// Include records that have passed their expiry but are not purged yet
all, err := client.SessionRepo().Query().WithExpired().All(ctx)
```

`WithExpired()` only has an effect on models embedding `som.Expiry`.

## Background Purge

When at least one model uses expiry, the client starts a purge goroutine on `NewClient`. It runs
`DELETE <table> WHERE expires_at < time::now()` for every expiry-enabled table on a fixed
interval and stops on `Close()`.

The interval defaults to one minute and can be configured:

```go
client, err := repo.NewClient(ctx, repo.Config{
    Address:             "ws://localhost:8000",
    Namespace:           "myapp",
    Database:            "prod",
    ExpiryPurgeInterval: 10 * time.Second,
})
```

Purge failures are logged via `slog` and do not stop the worker.

> Because the purge is client-side, records live slightly longer than their `expires_at` in the
> database. Query results are unaffected — the read filter hides them immediately.

## Expiry vs Soft Delete

| | Expiry | Soft delete |
|---|--------|-------------|
| Trigger | Time-based, automatic | Explicit `Delete()` call |
| Data retention | Deleted permanently by the purge | Kept until `Erase()` |
| Query escape hatch | `WithExpired()` | `WithDeleted()` |
| Restorable | No | Yes (`Restore()`) |

See [Soft Delete](06_soft_delete.md).
