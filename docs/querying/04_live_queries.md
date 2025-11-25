# Live Queries

SurrealDB supports real-time queries that push updates when data changes. SOM provides type-safe access to this feature through the `Live()` method.

## Basic Live Query

Subscribe to changes on a query:

```go
updates, err := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    Live(ctx)
if err != nil {
    return err
}

for update := range updates {
    if update.Error != nil {
        log.Println("Error:", update.Error)
        continue
    }

    switch update.Action {
    case "CREATE":
        fmt.Println("User created:", update.Data.Name)
    case "UPDATE":
        fmt.Println("User updated:", update.Data.Name)
    case "DELETE":
        fmt.Println("User deleted")
    }
}
```

## LiveResult Type

Each update is a `LiveResult` struct:

```go
type LiveResult[M any] struct {
    Action string  // "CREATE", "UPDATE", or "DELETE"
    Data   *M      // The affected record
    Error  error   // Any error that occurred
}
```

## Actions

| Action | Description |
|--------|-------------|
| `CREATE` | A new record matching the query was created |
| `UPDATE` | An existing record matching the query was updated |
| `DELETE` | A record matching the query was deleted |

## Filtered Live Queries

Live queries respect filters - you only receive updates for records that match:

```go
// Only receive updates for premium users
updates, err := client.UserRepo().Query().
    Filter(
        where.User.IsPremium.IsTrue(),
        where.User.IsActive.IsTrue(),
    ).
    Live(ctx)
```

## Cancellation

Use context cancellation to stop the live query:

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    updates, err := client.UserRepo().Query().Live(ctx)
    if err != nil {
        return
    }

    for update := range updates {
        // Handle updates
        fmt.Printf("%s: %+v\n", update.Action, update.Data)
    }
    // Channel closes when context is cancelled
    fmt.Println("Live query stopped")
}()

// Later, stop the subscription
cancel()
```

## Async Live Query

Start a live query without blocking:

```go
result := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    LiveAsync(ctx)

// Do setup work...

// Get the channel when ready
updates := <-result.Val()
err := <-result.Err()

for update := range updates {
    // Handle updates
}
```

## Error Handling

Handle errors that occur during the live subscription:

```go
for update := range updates {
    if update.Error != nil {
        // Connection error, parse error, etc.
        log.Printf("Live query error: %v", update.Error)

        // Decide whether to continue or break
        if isRecoverable(update.Error) {
            continue
        }
        break
    }

    // Process the update
    processUpdate(update)
}
```

## Real-World Example: Chat Application

```go
func SubscribeToMessages(ctx context.Context, roomID string) {
    updates, err := client.MessageRepo().Query().
        Filter(where.Message.RoomID.Equal(roomID)).
        Live(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for update := range updates {
        if update.Error != nil {
            continue
        }

        switch update.Action {
        case "CREATE":
            // New message - display to user
            displayMessage(update.Data)

        case "UPDATE":
            // Message edited - update display
            updateMessageDisplay(update.Data)

        case "DELETE":
            // Message deleted - remove from display
            removeMessageDisplay(update.Data.ID)
        }
    }
}
```

## Use Cases

Live queries are ideal for:

- **Real-time dashboards** - Metrics and KPIs that update automatically
- **Chat applications** - Instant message delivery
- **Collaborative editing** - See others' changes in real-time
- **Notification systems** - Push updates to users
- **Live activity feeds** - Social media style updates
- **Monitoring systems** - Real-time alerts and status updates
- **Gaming** - Player positions, scores, game state

## Best Practices

### Always Handle Errors

```go
for update := range updates {
    if update.Error != nil {
        log.Printf("Error: %v", update.Error)
        continue
    }
    // ...
}
```

### Cancel When Done

Always cancel the context when you no longer need updates:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // Ensure cleanup
```

### Use Specific Filters

Narrow your subscription to relevant records only:

```go
// Good - specific filter
updates, _ := client.UserRepo().Query().
    Filter(where.User.TeamID.Equal(currentTeamID)).
    Live(ctx)

// Avoid - subscribing to all records
updates, _ := client.UserRepo().Query().Live(ctx)
```

### Consider Connection Limits

Each live query maintains a connection. In production:

- Limit concurrent subscriptions per client
- Implement reconnection logic
- Monitor connection health
