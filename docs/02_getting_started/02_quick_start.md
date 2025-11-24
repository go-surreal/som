# Quick Start

This guide walks you through creating your first SOM-powered application.

> **Note**: This package is currently tested against SurrealDB version [2.3.10](https://surrealdb.com/releases#v2-3-10).

## Disclaimer

This library is currently considered **HIGHLY EXPERIMENTAL** and under heavy development. While basic functionality works, there could be unknown bugs. Not recommended for production use without thorough testing.

## Step 1: Define Your Model

Create `model/user.go`:

```go
package model

import "github.com/go-surreal/som"

type User struct {
    som.Node
    som.Timestamps  // Adds CreatedAt, UpdatedAt

    Username string
    Email    string
    Age      int
    IsActive bool
}
```

Key points:
- `som.Node` makes this a database record (required)
- `som.Timestamps` adds auto-managed timestamp fields (optional)
- Field names become database field names (lowercase by default)

## Step 2: Generate Code

```bash
go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som
```

This creates:
- `gen/som/` - Client and repository code
- `gen/som/where/` - Type-safe filters
- `gen/som/by/` - Sorting helpers

## Step 3: Connect and Use

```go
package main

import (
    "context"
    "fmt"
    "log"

    "yourproject/gen/som"
    "yourproject/gen/som/where"
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
    fmt.Printf("Created user with ID: %s\n", user.ID)

    // READ by ID
    found, exists, err := client.UserRepo().Read(ctx, user.ID)
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
        Filter(where.User.IsActive.IsTrue()).
        Filter(where.User.Age.GreaterThan(18)).
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
user, exists, err := client.UserRepo().Query().
    Filter(where.User.Email.Equal("john@example.com")).
    First(ctx)

if exists {
    fmt.Println(user.Username)
}
```

### Check Existence

```go
exists, err := client.UserRepo().Query().
    Filter(where.User.Email.Equal("john@example.com")).
    Exists(ctx)
```

### Count Records

```go
count, err := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    Count(ctx)
```

### Async Operations

```go
// Start query in background
result := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    AllAsync(ctx)

// Do other work...

// Get results when ready
users := <-result.Val()
err := <-result.Err()
```

## Complete Example

Here's a complete working example:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "yourproject/gen/som"
    "yourproject/gen/som/where"
    "yourproject/model"
)

func main() {
    ctx := context.Background()

    client, _ := som.NewClient(ctx, som.Config{
        Address:   "ws://localhost:8000",
        Username:  "root",
        Password:  "root",
        Namespace: "demo",
        Database:  "quickstart",
    })
    defer client.Close()

    // Create some users
    for i := 1; i <= 5; i++ {
        user := &model.User{
            Username: fmt.Sprintf("user%d", i),
            Email:    fmt.Sprintf("user%d@example.com", i),
            Age:      20 + i,
            IsActive: i%2 == 0,
        }
        client.UserRepo().Create(ctx, user)
    }

    // Query active users over 21
    users, _ := client.UserRepo().Query().
        Filter(
            where.User.IsActive.IsTrue(),
            where.User.Age.GreaterThan(21),
        ).
        All(ctx)

    for _, u := range users {
        fmt.Printf("%s (age %d)\n", u.Username, u.Age)
    }
}
```

## Next Steps

- [Core Concepts](03_concepts.md) - Understand Nodes, Edges, and Repositories
- [Models](../03_models/README.md) - Define complex data models
- [Querying](../06_querying/README.md) - Master the query builder
- [Relationships](../07_relationships/README.md) - Work with graph edges
