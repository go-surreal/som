# Nodes

Nodes are the primary way to define database records in SOM. A Node corresponds to a table in SurrealDB.

## Defining a Node

Embed `som.Node` in any struct to make it a database record:

```go
package model

import "github.com/go-surreal/som"

type User struct {
    som.Node

    Username string
    Email    string
    Age      int
}
```

## The ID Field

The `som.Node` embedding provides an `ID` field automatically. You don't need to define it yourself:

```go
// The ID is automatically available
user := &model.User{Username: "john"}
client.UserRepo().Create(ctx, user)
fmt.Println(user.ID)  // Populated after creation
```

## Field Tags

Use the `som` tag to customize field behavior:

```go
type User struct {
    som.Node

    Username string `som:"username"`        // Custom database field name
    Email    string `som:"email,omitempty"` // Omit if empty
    Internal string `som:"-"`               // Ignore this field
}
```

## Timestamps

You can add automatic timestamps by embedding `som.Timestamps`:

```go
type User struct {
    som.Node
    som.Timestamps  // Adds CreatedAt and UpdatedAt

    Username string
}
```

## Table Naming

By default, the table name is the lowercase struct name. A `User` struct becomes a `user` table.
