# Core Concepts

This page introduces the fundamental concepts you need to understand when working with SOM.

## Nodes

A **Node** represents a database record (similar to a row in SQL or a document in NoSQL). Any Go struct that embeds `som.Node` becomes a database table.

```go
type User struct {
    som.Node  // Makes this struct a database record

    Name  string
    Email string
}
```

The embedded `som.Node` provides:
- An `ID` field automatically managed by SurrealDB
- Timestamps (`CreatedAt`, `UpdatedAt`) if configured

## Edges

An **Edge** represents a relationship between two nodes. SurrealDB has first-class graph support, and SOM leverages this via structs that embed `som.Edge`.

```go
type Follows struct {
    som.Edge  // Makes this a relationship

    Since time.Time  // Edge properties
}
```

Edges have:
- `In` - The source node
- `Out` - The target node
- Custom properties defined in the struct

## Repositories

For each Node, SOM generates a **Repository** with standard CRUD operations:

```go
// Generated for the User model
client.UserRepo().Create(ctx, user)
client.UserRepo().Read(ctx, id)
client.UserRepo().Update(ctx, user)
client.UserRepo().Delete(ctx, id)
client.UserRepo().Query()  // Returns a query builder
```

## Query Builder

The **Query Builder** provides a fluent, type-safe API for constructing database queries:

```go
users, err := client.UserRepo().Query().
    Filter(where.User.Email.Contains("@example.com")).
    OrderBy(order.User.Name.Asc()).
    Limit(10).
    All(ctx)
```

All field names and operations are type-checked at compile time.

## Code Generation

SOM works through **code generation**. You define your models as Go structs, then run the generator to create:

- Repositories for each model
- Type-safe query builders
- Filter and ordering helpers
- Conversion utilities

This approach provides:
- Compile-time type safety
- Zero runtime reflection
- IDE autocompletion support
