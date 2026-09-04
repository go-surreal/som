# Core Concepts

This page introduces the fundamental concepts you need to understand when working with SOM.

## Nodes

A **Node** represents a database record (similar to a row in SQL or a document in NoSQL). Any Go struct that embeds `som.Node[T]` becomes a SurrealDB table.

```go
type User struct {
    som.Node[som.ULID]  // Required - makes this a database record

    Name  string
    Email string
}
```

The embedded `som.Node[T]` provides:
- An `ID()` method returning the node's ID (type depends on `T`)
- Table name derived from struct name (User -> `user`)

### ID Types

The type parameter `T` determines the ID format:

| Type | Description | Example ID |
|------|-------------|------------|
| `som.ULID` | ULID-based IDs (default choice) | `01HQMV8K2P...` |
| `som.UUID` | UUID-based IDs | `550e8400-e29b-...` |
| `som.Rand` | Random string IDs | `abc123def` |
| `som.String` | Application-supplied string IDs | `john` |
| Custom struct | Complex array/object IDs | `[city, date]` |

> **Note**: `ID()` returns the raw ID value. The repository automatically prefixes it with the table name (e.g. `user:01HQMV8K2P...`) when storing and querying records.

## Timestamps

Embed `som.Timestamps` for automatic time tracking:

```go
type User struct {
    som.Node[som.ULID]
    som.Timestamps  // Adds CreatedAt, UpdatedAt

    Name string
}
```

This adds two database fields, read via accessor methods:
- `CreatedAt() time.Time` - set automatically on creation (readonly)
- `UpdatedAt() time.Time` - updated automatically on every modification

These fields are managed by SurrealDB and cannot be manually set.

## Edges

An **Edge** represents a graph relationship between two nodes. SurrealDB has first-class graph support, and SOM leverages this via structs that embed `som.Edge`.

```go
type Follows struct {
    som.Edge        // Required - makes this a relationship
    som.Timestamps  // Optional

    From  User `som:"in"`   // Required - source node
    To    User `som:"out"`  // Required - target node

    Since time.Time // Edge metadata
}
```

The embedded `som.Edge` provides the edge's own `ID()`. The connected nodes are declared by you,
tagged `som:"in"` and `som:"out"` — the field names are free.

For the edge to be usable it must also be declared as a field on the source node:

```go
type User struct {
    som.Node[som.ULID]

    Name    string
    Follows []Follows
}
```

Edges are created through that node's repository, which issues a RELATE statement:

```go
err := client.UserRepo().Relate().Follows().Create(ctx, follows)
// RELATE user:alice->follows->user:bob
```

## ID Handling

```go
// After Create(), the ID is populated
user := &model.User{Name: "Alice"}
client.UserRepo().Create(ctx, user)
fmt.Println(user.ID())  // 01HQMV... (raw ULID)

// Create with specific ID
client.UserRepo().CreateWithID(ctx, "alice", user)
// Creates: user:alice
```

## Repositories

For each Node, SOM generates a **Repository** with standard operations:

```go
// Repository interface (generated)
type UserRepo interface {
    // Create new record (ID auto-generated)
    Create(ctx context.Context, user *model.User) error

    // Create with specific ID
    CreateWithID(ctx context.Context, id string, user *model.User) error

    // Bulk insert multiple records
    Insert(ctx context.Context, users []*model.User) error

    // Read by ID - returns (model, exists, error)
    Read(ctx context.Context, id string) (*model.User, bool, error)

    // Update existing record
    Update(ctx context.Context, user *model.User) error

    // Delete record
    Delete(ctx context.Context, user *model.User) error

    // Refresh record from database
    Refresh(ctx context.Context, user *model.User) error

    // Edge creation for edges declared on this node
    Relate() *relate.User

    // Index access (e.g. per-index Rebuild)
    Index() *index.User

    // Get query builder
    Query() query.Builder[model.User]

    // Lifecycle hooks
    OnBeforeCreate(fn func(ctx context.Context, node *model.User) error) func()
    OnAfterCreate(fn func(ctx context.Context, node *model.User) error) func()
    OnBeforeUpdate(fn func(ctx context.Context, node *model.User) error) func()
    OnAfterUpdate(fn func(ctx context.Context, node *model.User) error) func()
    OnBeforeDelete(fn func(ctx context.Context, node *model.User) error) func()
    OnAfterDelete(fn func(ctx context.Context, node *model.User) error) func()
}
```

Access repositories through the client:

```go
client.UserRepo()  // User repository
client.PostRepo()  // Post repository
```

Repositories exist for nodes, [views](../models/07_views.md) and
[sinks](../models/08_sinks.md). Edges have none — they are created via `Relate()` on the source
node's repository.

## Query Builder

The **Query Builder** provides a fluent, type-safe API for constructing database queries:

```go
users, err := client.UserRepo().Query().
    Where(filter.User.Email.Contains("@example.com")).
    Order(by.User.Name.Asc()).
    Limit(10).
    All(ctx)
```

### Execution Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `All(ctx)` | `([]*Model, error)` | All matching records |
| `First(ctx)` | `(*Model, error)` | First match (returns `ErrNotFound` if none) |
| `Count(ctx)` | `(int, error)` | Count of matches |
| `Exists(ctx)` | `(bool, error)` | Whether any exist |
| `Live(ctx)` | `(<-chan LiveResult[*Model], error)` | Real-time stream |

### Async Variants

Most execution methods have an async version (`Live` does not — it already returns a channel):

```go
result := client.UserRepo().Query().AllAsync(ctx)
// ... do other work ...
users := <-result.Val()
err := <-result.Err()
```

## Filters (filter)

Type-safe conditions for queries. Import from `gen/som/filter`:

```go
import "yourproject/gen/som/filter"

// Equality
filter.User.Name.Equal("Alice")

// Comparison
filter.User.Age.GreaterThan(18)

// String operations
filter.User.Email.Contains("@gmail.com").True()

// Multiple conditions (AND)
client.UserRepo().Query().
    Where(
        filter.User.IsActive.True(),
        filter.User.Age.GreaterThan(18),
    )
```

## Sorting (by)

Order results by fields. Import from `gen/som/by`:

```go
import "yourproject/gen/som/by"

// Ascending
by.User.Name.Asc()

// Descending
by.User.CreatedAt.Desc()

// Multiple sort criteria
client.UserRepo().Query().
    Order(
        by.User.LastName.Asc(),
        by.User.FirstName.Asc(),
    )
```

## Code Generation

SOM works through **code generation**. The workflow:

1. **Define models** as Go structs with `som.Node[T]` or `som.Edge`
2. **Run generator**: `som -i ./model`
   (optionally with a `//go:build som` definition file for analyzers and views —
   see [Schema Definitions](../code_generation/03_definitions.md))
3. **Import and use** the generated packages

Benefits:
- **Compile-time safety** - Typos caught by compiler
- **Zero reflection** - No runtime type inspection
- **IDE support** - Full autocompletion
- **Performance** - No overhead from ORM magic

Regenerate whenever you:
- Add or remove model fields
- Create new Node or Edge types
- Change field types
