# FAQ

Frequently asked questions about SOM.

## General

### What is SOM?

SOM (SurrealDB Object Mapper) is a code generation tool that creates type-safe database access code for SurrealDB from Go struct definitions. It generates repositories, query builders, and filter helpers.

### Is SOM production-ready?

SOM is currently in **early development** and considered experimental. It works for many use cases but may have unknown bugs. Use with caution in production and pin to specific versions.

### What Go version is required?

Go 1.25 or later is required due to heavy use of generics and iterators.

### What SurrealDB version is supported?

SOM is tested against SurrealDB [3.0.0](https://surrealdb.com/releases). Compatibility with older versions is not guaranteed.

## Technical

### Why does SOM use code generation?

Code generation provides:

- **Compile-time type safety** - Catch errors before runtime
- **Zero runtime reflection** - Better performance
- **Full IDE support** - Autocompletion, refactoring, go-to-definition
- **Readable generated code** - Debug and understand what's happening

### What's the difference between som.Node[T] and som.Edge?

- **som.Node[T]** - A database record/table. The type parameter `T` determines the ID format (`som.ULID`, `som.UUID`, `som.Rand`, or a custom struct).
- **som.Edge** - A relationship between two nodes. Has an ID plus `In` and `Out` fields with `som:"in"` and `som:"out"` tags.

### Why does Read return (record, bool, error)?

The three-value return distinguishes between:

- Record found: `(record, true, nil)`
- Record not found: `(nil, false, nil)` - Not an error, just doesn't exist
- Database error: `(nil, false, err)` - Actual error occurred

This avoids conflating "not found" with errors.

### What does First return?

`First(ctx)` returns `(*Model, error)`. If no record matches, it returns `ErrNotFound`:

```go
user, err := client.UserRepo().Query().
    Where(filter.User.Email.Equal("john@example.com")).
    First(ctx)

if errors.Is(err, som.ErrNotFound) {
    // No matching record
}
```

### How do I use optional fields?

Use pointers for optional fields:

```go
type User struct {
    som.Node[som.ULID]

    Name     string   // Required
    Nickname *string  // Optional, can be nil
    Age      *int     // Optional
}
```

Query with `IsNil()` and `IsNotNil()`:

```go
filter.User.Nickname.IsNotNil()  // Has a nickname
filter.User.Age.IsNil()          // Age not set
```

### How do automatic timestamps work?

Embed `som.Timestamps` for automatic tracking:

```go
type User struct {
    som.Node[som.ULID]
    som.Timestamps  // Adds CreatedAt and UpdatedAt

    Name string
}
```

- `CreatedAt` - Set on create, readonly
- `UpdatedAt` - Updated on every save

Both are managed by SurrealDB and read-only in your code.

### Can I use raw SurrealQL queries?

Currently, SOM focuses on the type-safe query builder. Raw query support may be added in the future.

### How do I handle database migrations?

Migration support is planned but not yet implemented. Currently, schema changes require manual handling in SurrealDB.

### What types are NOT supported?

Currently not supported:

- `uint`, `uint64` (SurrealDB integer limitations)
- `complex64`, `complex128`
- Channels, functions
- Maps (except specific patterns)
- Recursive types

## Troubleshooting

### The generator isn't finding my models

Ensure your structs embed `som.Node[T]` or `som.Edge`:

```go
type User struct {
    som.Node[som.ULID]  // Required!
    Name string
}
```

### I get import errors in generated code

Regenerate the code after any model changes:

```bash
rm -rf ./gen/som
go run github.com/go-surreal/som@latest -i ./model
go mod tidy
```

### ID field shows as nil after Create

Make sure you're using a pointer and checking the `ID()` method:

```go
user := &model.User{Name: "Alice"}
err := client.UserRepo().Create(ctx, user)
fmt.Println(user.ID())  // Use ID() method, not a field
```

### Live query channel closes unexpectedly

Check context cancellation and error handling:

```go
for update := range updates {
    if update.Error != nil {
        log.Printf("Error: %v", update.Error)
        // Connection lost, consider reconnecting
        break
    }
    // Process update
}
```

## More Questions?

- Check [GitHub Issues](https://github.com/go-surreal/som/issues)
- Open a [GitHub Discussion](https://github.com/go-surreal/som/discussions)
- Read the [API Reference](../api_reference/README.md)
