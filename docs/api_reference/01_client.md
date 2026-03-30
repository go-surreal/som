# Client API

The generated client provides the main entry point for database operations.

## Creating a Client

```go
client, err := som.NewClient(ctx, som.Config{
    Address:   "ws://localhost:8000",
    Username:  "root",
    Password:  "root",
    Namespace: "myapp",
    Database:  "production",
})
if err != nil {
    log.Fatal(err)
}
```

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
}
```

## Version Verification

When creating a client, SOM automatically verifies that the connected SurrealDB server meets the minimum required version (currently **3.0.5**). If the version check fails, `NewClient` returns a `som.ErrUnsupportedVersion` error:

```go
client, err := som.NewClient(ctx, config)
if err != nil {
    if errors.Is(err, som.ErrUnsupportedVersion) {
        log.Fatal("SurrealDB server version too old, please upgrade to 3.0.5+")
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

// For an edge type
followsRepo := client.FollowsRepo()
```

## Connection Management

### Closing the Client

```go
err := client.Close()
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

var users []map[string]any
err = result.Scan(&users)
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
| `som.ErrUnsupportedVersion` | Server version below minimum required |

## Thread Safety

The client is safe for concurrent use from multiple goroutines.
