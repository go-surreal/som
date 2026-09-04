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

SOM is tested against SurrealDB [3.2.0](https://surrealdb.com/releases#v3-2-0). Compatibility with older versions is not guaranteed.

## Technical

### Why does SOM use code generation?

Code generation provides:

- **Compile-time type safety** - Catch errors before runtime
- **Zero runtime reflection** - Better performance
- **Full IDE support** - Autocompletion, refactoring, go-to-definition
- **Readable generated code** - Debug and understand what's happening

### What's the difference between som.Node[T] and som.Edge?

- **som.Node[T]** - A database record/table. The type parameter `T` determines the ID format (`som.ULID`, `som.UUID`, `som.Rand`, `som.String`, or a custom struct).
- **som.Edge** - A relationship between two nodes. Provides the edge's own ID; you declare the
  connected nodes yourself as two fields tagged `som:"in"` and `som:"out"`. Edges have no
  repository — they are created via `Relate()` on the source node's repository.
- **som.View** - A read-only table computed by the database (`DEFINE TABLE ... AS SELECT`).
- **som.Sink** - A write-only ingestion table whose rows are discarded after the write.

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

Query with `Nil(true)` and `Nil(false)`:

```go
filter.User.Nickname.Nil(false)  // Has a nickname
filter.User.Age.Nil(true)          // Age not set
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

Yes. `client.Raw(ctx, statement, som.Params{...})` executes arbitrary SurrealQL and returns a
result that can be scanned into your own types. See
[Raw Queries](../querying/06_raw_queries.md).

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

Ensure your structs embed `som.Node[T]`, `som.Edge`, `som.View` or `som.Sink`:

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

The channel closes when the context is cancelled or the server terminates the subscription. The
latter arrives as a `LiveKilled` event:

```go
for update := range updates {
    switch res := update.(type) {
    case query.LiveKilled[*model.User]:
        // server terminated the live query; re-subscribe if needed
        return
    case query.LiveUpdate[*model.User]:
        user, err := res.Get()
        ...
    }
}
```

## More Questions?

- Check [GitHub Issues](https://github.com/go-surreal/som/issues)
- Open a [GitHub Discussion](https://github.com/go-surreal/som/discussions)
- Read the [API Reference](../api_reference/README.md)
