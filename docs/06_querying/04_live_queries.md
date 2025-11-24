# Live Queries

SurrealDB supports real-time queries that push updates when data changes. SOM provides type-safe access to this feature.

## Basic Live Query

Subscribe to changes on a query:

```go
live, err := client.UserRepo().Query().
    Filter(where.User.IsActive.Equal(true)).
    Live(ctx)
if err != nil {
    return err
}

for update := range live {
    switch res := update.(type) {
    case query.LiveCreate[model.User]:
        user, _ := res.Get()
        fmt.Println("User created:", user.Name)

    case query.LiveUpdate[model.User]:
        user, _ := res.Get()
        fmt.Println("User updated:", user.Name)

    case query.LiveDelete[model.User]:
        fmt.Println("User deleted")
    }
}
```

## Event Types

Live queries emit three types of events:

- `LiveCreate` - A new record matching the query was created
- `LiveUpdate` - An existing record was updated
- `LiveDelete` - A record was deleted

## Cancellation

Use context cancellation to stop the live query:

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    live, _ := client.UserRepo().Query().Live(ctx)
    for update := range live {
        // Handle updates
    }
}()

// Later, stop the subscription
cancel()
```

## Use Cases

Live queries are ideal for:

- Real-time dashboards
- Chat applications
- Collaborative editing
- Notification systems
- Live activity feeds

## Considerations

- Live queries maintain an open connection
- Cancel subscriptions when no longer needed
- Handle reconnection logic for production use
