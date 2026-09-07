# Changefeed

A changefeed makes SurrealDB retain the history of a table for a given duration. SOM exposes it
as a typed `SHOW CHANGES` query, letting you replay creates, updates and deletes after the fact —
unlike [live queries](04_live_queries.md), which only push changes while you are connected.

## Enabling a Changefeed

Tag the embedded `som.Node` (or `som.Edge`) with the retention duration:

```go
type Order struct {
    som.Node[som.ULID] `som:"changefeed=1d"`

    Number string
    Total  float64
}
```

This appends `CHANGEFEED 1d` to the table definition:

```surql
DEFINE TABLE order SCHEMAFULL TYPE NORMAL CHANGEFEED 1d PERMISSIONS FULL;
```

The duration determines how far back changes are kept. Models without the tag have no `Changes()`
method.

## Querying Changes

`Changes()` returns a builder. A starting point is mandatory — either a timestamp or a
versionstamp:

```go
entries, err := client.OrderRepo().Changes().
    Since(time.Now().Add(-1 * time.Hour)).
    Show(ctx)
```

```go
// Versionstamp 0 replays everything the feed still holds
entries, err := client.OrderRepo().Changes().
    SinceVersionstamp(0).
    Limit(100).
    Show(ctx)
```

| Method | Description |
|--------|-------------|
| `Since(t time.Time)` | Start at a timestamp |
| `SinceVersionstamp(v uint64)` | Start at a versionstamp (monotonically increasing) |
| `Limit(n int)` | Maximum number of change batches |
| `Show(ctx)` | Execute, returns `[]ChangeEntry[*Model]` |
| `Describe()` | Render the statement for debugging |

Calling `Show` without `Since` or `SinceVersionstamp` returns an error.

## Change Entries

```go
type ChangeEntry[M any] struct {
    Versionstamp uint64
    Creates      []M
    Updates      []M
    Deletes      []M
}
```

Each entry is one batch of changes at a versionstamp, with the records already decoded into your
model type:

```go
for _, entry := range entries {
    for _, order := range entry.Updates {
        fmt.Printf("v%d updated: %s\n", entry.Versionstamp, order.Number)
    }
    for _, order := range entry.Deletes {
        fmt.Printf("v%d deleted: %s\n", entry.Versionstamp, order.Number)
    }
}
```

Store the highest `Versionstamp` you have processed and pass it to `SinceVersionstamp` on the next
run to continue where you left off.

> SurrealDB reports most record changes — including newly created records — under `update`.
> Do not rely on `Creates` being populated for every insert; treat `Updates` as "current state
> of the record at that versionstamp".

Table definition changes contained in the feed are skipped.

## Limitations

- The pre-image (`INCLUDE ORIGINAL`) is not supported — you only get the resulting record.
- Changes older than the configured retention are dropped by the database.
- The feed is per table; there is no cross-table ordering guarantee beyond the versionstamp.

## Changefeed vs Live Queries vs Timestamps

| Need | Use |
|------|-----|
| Push updates to a connected client | [Live Queries](04_live_queries.md) |
| Replay what happened while you were offline | Changefeed |
| Know when a record was last modified | `som.Timestamps` |
