# Basic CRUD Example

This example demonstrates a complete CRUD application using SOM.

## Model Definition

Create `model/user.go`:

```go
package model

import (
    "time"
    "github.com/go-surreal/som"
)

type User struct {
    som.Node

    Name      string
    Email     string
    Age       int
    IsActive  bool
    CreatedAt time.Time
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
    "time"

    "yourproject/gen/som"
    "yourproject/gen/som/where"
    "yourproject/gen/som/order"
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
        Name:      "Alice",
        Email:     "alice@example.com",
        Age:       30,
        IsActive:  true,
        CreatedAt: time.Now(),
    }

    err = client.UserRepo().Create(ctx, user)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created user with ID: %s\n", user.ID)

    // READ
    retrieved, err := client.UserRepo().Read(ctx, user.ID)
    if err != nil {
        log.Fatal(err)
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
        Filter(where.User.IsActive.Equal(true)).
        OrderBy(order.User.Name.Asc()).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d active users\n", len(activeUsers))

    // DELETE
    err = client.UserRepo().Delete(ctx, user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Deleted user")
}
```

## Running the Example

1. Start SurrealDB:
   ```bash
   surreal start --user root --pass root
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

## Expected Output

```
Created user with ID: user:abc123
Retrieved user: Alice
Updated user
Found 1 active users
Deleted user
```
