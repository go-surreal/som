# Basic CRUD Example

This example demonstrates a complete CRUD application using SOM.

## Model Definition

Create `model/user.go`:

```go
package model

import (
    "github.com/go-surreal/som"
)

type User struct {
    som.Node
    som.Timestamps  // Adds CreatedAt and UpdatedAt

    Name     string
    Email    string
    Age      int
    IsActive bool
}
```

## Generate Code

```bash
go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som
```

## Application Code

Create `main.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/where"
    "yourproject/model"
)

func main() {
    ctx := context.Background()

    // Connect to database
    client, err := som.NewClient(ctx, som.Config{
        Address:   "ws://localhost:8000",
        Username:  "root",
        Password:  "root",
        Namespace: "example",
        Database:  "crud",
    })
    if err != nil {
        log.Fatal(err)
    }

    // CREATE
    user := &model.User{
        Name:     "Alice",
        Email:    "alice@example.com",
        Age:      30,
        IsActive: true,
    }

    err = client.UserRepo().Create(ctx, user)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created user with ID: %s\n", user.ID())

    // READ
    // Note: Read returns (record, exists, error)
    retrieved, exists, err := client.UserRepo().Read(ctx, user.ID())
    if err != nil {
        log.Fatal(err)
    }
    if !exists {
        log.Fatal("User not found")
    }
    fmt.Printf("Retrieved user: %s\n", retrieved.Name)

    // UPDATE
    retrieved.Name = "Alice Smith"
    retrieved.Age = 31
    err = client.UserRepo().Update(ctx, retrieved)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Updated user")

    // QUERY
    activeUsers, err := client.UserRepo().Query().
        Filter(where.User.IsActive.IsTrue()).
        Order(by.User.Name.Asc()).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d active users\n", len(activeUsers))

    // DELETE
    err = client.UserRepo().Delete(ctx, retrieved)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Deleted user")
}
```

## Running the Example

1. Start SurrealDB:
   ```bash
   surreal start --user root --pass root memory
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

## Expected Output

```
Created user with ID: user:01HQ...
Retrieved user: Alice
Updated user
Found 1 active users
Deleted user
```

## Key Points

### Create

- Pass a pointer to your model
- The `ID()` method is populated after successful creation
- `CreatedAt` and `UpdatedAt` are set automatically (with `som.Timestamps`)

### Read

- Returns three values: `(record, exists, error)`
- Check `exists` before using the record
- Returns `nil, false, nil` if record doesn't exist (not an error)

### Update

- The record must have a valid ID from a previous Create or Read
- `UpdatedAt` is updated automatically (with `som.Timestamps`)

### Delete

- Pass the record (not just the ID)
- The record must have a valid ID

### Query

- Use the fluent builder pattern
- Filter with type-safe conditions from `where` package
- Order with helpers from `by` package
- Execute with `All()`, `First()`, `One()`, `Count()`, or `Exists()`

## Error Handling Pattern

```go
user, exists, err := client.UserRepo().Read(ctx, id)

// Always check error first
if err != nil {
    return fmt.Errorf("database error: %w", err)
}

// Then check existence
if !exists {
    return ErrUserNotFound
}

// Now safely use the record
fmt.Println(user.Name)
```

## Using CreateWithID

For custom IDs instead of auto-generated ULIDs:

```go
user := &model.User{
    Name:  "Bob",
    Email: "bob@example.com",
}

// Creates user:bob instead of user:01HQ...
err := client.UserRepo().CreateWithID(ctx, "bob", user)
```

## Refreshing Records

Reload a record to get the latest data:

```go
// After potential updates from other sources
err := client.UserRepo().Refresh(ctx, user)
if err != nil {
    log.Fatal(err)
}
// user now has current database values
```
