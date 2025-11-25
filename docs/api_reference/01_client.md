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

user, err := client.UserRepo().Read(ctx, id)
```

## Error Handling

Check for errors on all operations:

```go
user, err := client.UserRepo().Read(ctx, id)
if err != nil {
    // Handle connection errors, not found, etc.
    return err
}
```

## Thread Safety

The client is safe for concurrent use from multiple goroutines.
