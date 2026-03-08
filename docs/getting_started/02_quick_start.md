# Quick Start

This guide walks you through creating your first SOM-powered application.

> **Note**: This package is currently tested against SurrealDB version [3.0.0](https://surrealdb.com/releases).

## Disclaimer

This library is currently considered **HIGHLY EXPERIMENTAL** and under heavy development. While basic functionality works, there could be unknown bugs. Not recommended for production use without thorough testing.

## Step 1: Define Your Model

First, generate the static base files so types like `som.Node` are available:

```bash
go run github.com/go-surreal/som@latest
```

Create `model/user.go`:

```go
package model

import "yourproject/gen/som"

type User struct {
    som.Node[som.ULID]
    som.Timestamps

    Username string
    Email    string
    Age      int
    IsActive bool
}
```

Key points:
- `som.Node[som.ULID]` makes this a database record with ULID-based IDs (required)
- `som.Timestamps` adds auto-managed timestamp fields (optional)
- Field names become database field names (snake_case by default)
- You can override field names with the `som:"custom_name"` tag

## Step 2: Generate Code

```bash
go run github.com/go-surreal/som@latest -i ./model
```

This creates:
- `gen/som/` - Client and repository code
- `gen/som/filter/` - Type-safe filters
- `gen/som/by/` - Sorting helpers

## Step 3: Connect and Use

```go
package main

import (
    "context"
    "fmt"
    "log"

    "yourproject/gen/som"
    "yourproject/gen/som/filter"
    "yourproject/gen/som/by"
    "yourproject/model"
)

func main() {
    ctx := context.Background()

    // Connect to SurrealDB
    client, err := som.NewClient(ctx, som.Config{
        Address:   "ws://localhost:8000",
        Username:  "root",
        Password:  "root",
        Namespace: "myapp",
        Database:  "dev",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // CREATE
    user := &model.User{
        Username: "johndoe",
        Email:    "john@example.com",
        Age:      28,
        IsActive: true,
    }
    if err := client.UserRepo().Create(ctx, user); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created user with ID: %s\n", user.ID())

    // READ by ID
    found, exists, err := client.UserRepo().Read(ctx, string(user.ID()))
    if err != nil {
        log.Fatal(err)
    }
    if exists {
        fmt.Printf("Found: %s\n", found.Username)
    }

    // UPDATE
    user.Age = 29
    if err := client.UserRepo().Update(ctx, user); err != nil {
        log.Fatal(err)
    }

    // QUERY with filters
    users, err := client.UserRepo().Query().
        Where(filter.User.IsActive.IsTrue()).
        Where(filter.User.Age.GreaterThan(18)).
        Order(by.User.Username.Asc()).
        Limit(10).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d active adult users\n", len(users))

    // DELETE
    if err := client.UserRepo().Delete(ctx, user); err != nil {
        log.Fatal(err)
    }
}
```

## Step 4: Query Patterns

### Find First Match

```go
user, err := client.UserRepo().Query().
    Where(filter.User.Email.Equal("john@example.com")).
    First(ctx)

if err == nil {
    fmt.Println(user.Username)
}
```

### Check Existence

```go
exists, err := client.UserRepo().Query().
    Where(filter.User.Email.Equal("john@example.com")).
    Exists(ctx)
```

### Count Records

```go
count, err := client.UserRepo().Query().
    Where(filter.User.IsActive.IsTrue()).
    Count(ctx)
```

### Async Operations

```go
// Start query in background
result := client.UserRepo().Query().
    Where(filter.User.IsActive.IsTrue()).
    AllAsync(ctx)

// Do other work...

// Get results when ready
users := <-result.Val()
err := <-result.Err()
```

## Next Steps

- [Core Concepts](03_concepts.md) - Understand Nodes, Edges, and Repositories
- [Models](../models/README.md) - Define complex data models
- [Querying](../querying/README.md) - Master the query builder
- [Relationships](../relationships/README.md) - Work with graph edges
