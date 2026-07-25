# Live Queries

SurrealDB supports real-time queries that push updates when data changes. SOM provides type-safe access to this feature through the `Live()` and `LiveCount()` methods.

## Basic Live Query

Subscribe to changes on a query:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

liveChan, err := client.UserRepo().Query().
    Where(filter.User.IsActive.IsTrue()).
    Live(ctx)
if err != nil {
    return err
}

for update := range liveChan {
    switch res := update.(type) {

    case query.LiveCreate[*model.User]:
        user, err := res.Get()
        if err != nil {
            log.Println("Error:", err)
            continue
        }
        fmt.Println("User created:", user.Name)

    case query.LiveUpdate[*model.User]:
        user, err := res.Get()
        if err != nil {
            log.Println("Error:", err)
            continue
        }
        fmt.Println("User updated:", user.Name)

    case query.LiveDelete[*model.User]:
        user, err := res.Get()
        if err != nil {
            log.Println("Error:", err)
            continue
        }
        fmt.Println("User deleted:", user.Name)

    case query.LiveKilled[*model.User]:
        fmt.Println("Live query was terminated by the server")
        return
    }
}
```

## Event Types

Each update received on the channel implements `LiveResult[M]`. Use type assertions to determine the event kind:

| Type | Description |
|------|-------------|
| `LiveCreate[M]` | A new record matching the query was created |
| `LiveUpdate[M]` | An existing record matching the query was updated |
| `LiveDelete[M]` | A record matching the query was deleted |
| `LiveKilled[M]` | The live query was terminated by the server |

`LiveCreate`, `LiveUpdate`, and `LiveDelete` each provide a `Get() (M, error)` method to access the affected record.

`LiveKilled` is emitted when the server terminates the live query (e.g., table removed, permissions changed). It does not carry data.

## Live Count

Track the number of matching records in real time:

```go
liveCount, err := client.UserRepo().Query().
    Where(filter.User.IsActive.IsTrue()).
    LiveCount(ctx)
if err != nil {
    return err
}

for count := range liveCount {
    fmt.Printf("Active users: %d\n", count)
}
```

`LiveCount` first executes a `Count()` query for the initial value, then adjusts the count as `LiveCreate` and `LiveDelete` events arrive.

## Filtered Live Queries

Live queries respect filters — you only receive updates for records that match:

```go
liveChan, err := client.UserRepo().Query().
    Where(
        filter.User.IsPremium.IsTrue(),
        filter.User.IsActive.IsTrue(),
    ).
    Live(ctx)
```

## Fetching Related Records

Live queries support the `Fetch()` clause to include related records:

```go
liveChan, err := client.UserRepo().Query().
    Where(filter.User.IsActive.IsTrue()).
    Fetch(with.User.Organization()).
    Live(ctx)
if err != nil {
    return err
}

for update := range liveChan {
    switch res := update.(type) {
    case query.LiveCreate[*model.User]:
        user, _ := res.Get()
        // user.Organization is fully populated
        fmt.Printf("New user %s in org %s\n", user.Name, user.Organization.Name)
    }
}
```

## Cancellation

Use context cancellation to stop the live query. When the context is cancelled, a best-effort attempt is made to kill the live query on the server and the channel is closed:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

liveChan, err := client.UserRepo().Query().Live(ctx)
if err != nil {
    return err
}

go func() {
    for update := range liveChan {
        // Handle updates...
    }
    // Channel closes when context is cancelled
    fmt.Println("Live query stopped")
}()

// Later, stop the subscription
cancel()
```

## Builder Constraints

Not all query builder methods are compatible with live queries. The following methods return `BuilderNoLive`, which does not expose `Live()` or `LiveCount()`:

- `Order()` / `OrderRandom()`
- `Start()` / `Limit()` / `Range()`
- `Timeout()` / `Parallel()` / `TempFiles()`

This is enforced at compile time — you cannot accidentally build an invalid live query.

## Use Cases

Live queries are ideal for:

- **Real-time dashboards** — Metrics and KPIs that update automatically
- **Chat applications** — Instant message delivery
- **Collaborative editing** — See others' changes in real-time
- **Notification systems** — Push updates to users
- **Live activity feeds** — Social media style updates
- **Monitoring systems** — Real-time alerts and status updates

## Best Practices

### Cancel When Done

Always cancel the context when you no longer need updates to clean up the server-side live query:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

### Use Specific Filters

Narrow your subscription to relevant records only:

```go
// Good — specific filter
liveChan, _ := client.UserRepo().Query().
    Where(filter.User.TeamID.Equal(currentTeamID)).
    Live(ctx)

// Avoid — subscribing to all records
liveChan, _ := client.UserRepo().Query().Live(ctx)
```

### Start Live Before Initial Query

If you need both the current result set and live updates, start the live query first to avoid missing updates between the initial query and subscription:

```go
liveChan, err := client.UserRepo().Query().Live(ctx)
// then
users, err := client.UserRepo().Query().All(ctx)
```
