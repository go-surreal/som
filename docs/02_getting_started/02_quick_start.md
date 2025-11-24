# Quick Start

This guide walks you through creating your first SOM-powered application.

> **Note**: This package is currently tested against SurrealDB version [2.3.10](https://surrealdb.com/releases#v2-3-10).

## Disclaimer

This library is currently considered **HIGHLY EXPERIMENTAL** and under heavy development.

SOM is in the stage of early development. While the basic functionality should be working as expected, there could be unknown and critical bugs. This could theoretically lead to your database and the data within it to get broken. It is not yet recommended for production use, so please be careful!

If you still want to use it for production applications, we do not accept any liability for data loss or other kind of issues. You have been warned!

## Step 1: Define Your Model

Create a model file at `model/user.go`:

```go
package model

import "github.com/go-surreal/som"

type User struct {
    som.Node

    Username string `som:"username"`
    Password string `som:"password"`
    Email    string `som:"email"`
}
```

The `som.Node` embedding tells SOM this struct represents a database record. It automatically provides an `ID` field.

## Step 2: Generate the Client

Run the code generator:

```bash
go run github.com/go-surreal/som/cmd/som@latest gen ./model ./gen/som
```

This creates type-safe repository and query builder code in `./gen/som`.

## Step 3: Use the Generated Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "yourproject/gen/som"
    "yourproject/gen/som/where"
    "yourproject/model"
)

func main() {
    ctx := context.Background()

    // Create a new client
    client, err := som.NewClient(ctx, som.Config{
        Address:   "ws://localhost:8000",
        Username:  "root",
        Password:  "root",
        Namespace: "test",
        Database:  "test",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a user
    user := &model.User{
        Username: "johndoe",
        Password: "secret",
        Email:    "john@example.com",
    }

    err = client.UserRepo().Create(ctx, user)
    if err != nil {
        log.Fatal(err)
    }

    // Query users by email
    found, err := client.UserRepo().Query().
        Filter(
            where.User.Email.Equal("john@example.com"),
        ).
        First(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found user: %s\n", found.Username)
}
```

## Next Steps

- Learn about [Models](../03_models/README.md) in depth
- Explore [Data Types](../04_data_types/README.md) supported by SOM
- Master the [Query Builder](../06_querying/README.md)
