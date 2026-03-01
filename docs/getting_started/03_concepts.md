# Core Concepts

This page introduces the fundamental concepts you need to understand when working with SOM.

## Nodes

A **Node** represents a database record (similar to a row in SQL or a document in NoSQL). Any Go struct that embeds `som.Node` becomes a SurrealDB table.

```go
type User struct {
    som.Node  // Required - makes this a database record

    Name  string
    Email string
}
```

The embedded `som.Node` provides:
- An `ID` field of type `*som.ID` (auto-generated ULID on create)
- Table name derived from struct name (User → `user`)

## Timestamps

Embed `som.Timestamps` for automatic time tracking:

```go
type User struct {
    som.Node
    som.Timestamps  // Adds CreatedAt, UpdatedAt

    Name string
}
```

This adds:
- `CreatedAt time.Time` - Set automatically on creation (readonly)
- `UpdatedAt time.Time` - Updated automatically on every modification

These fields are managed by SurrealDB and cannot be manually set.

## Edges

An **Edge** represents a graph relationship between two nodes. SurrealDB has first-class graph support, and SOM leverages this via structs that embed `som.Edge`.

```go
type Follows struct {
    som.Edge        // Required - makes this a relationship
    som.Timestamps  // Optional

    Since time.Time // Edge metadata
}
```

Edges automatically have:
- `ID` - Unique identifier for the edge itself
- `In` - The source node (where the relationship starts)
- `Out` - The target node (where the relationship points)

Creating edges uses the RELATE statement:
```
RELATE user:alice->follows->user:bob
```

## ID Handling

SOM provides utilities for working with SurrealDB record IDs:

```go
// After Create(), the ID is populated
user := &model.User{Name: "Alice"}
client.UserRepo().Create(ctx, user)
fmt.Println(user.ID)  // user:01HQMV...

// Create with specific ID
client.UserRepo().CreateWithID(ctx, "alice", user)
// Creates: user:alice

// Create ID manually
id := som.NewRecordID("user", "alice")

// Create pointer to ID
idPtr := som.MakeID("user", "alice")

// Reference table for auto-generated IDs
table := som.Table("user")
```

## Repositories

For each Node and Edge, SOM generates a **Repository** with standard operations:

```go
// Repository interface (generated)
type UserRepo interface {
    // Create new record (ID auto-generated)
    Create(ctx context.Context, user *model.User) error

    // Create with specific ID
    CreateWithID(ctx context.Context, id string, user *model.User) error

    // Read by ID - returns (model, exists, error)
    Read(ctx context.Context, id *som.ID) (*model.User, bool, error)

    // Update existing record
    Update(ctx context.Context, user *model.User) error

    // Delete record
    Delete(ctx context.Context, user *model.User) error

    // Refresh record from database
    Refresh(ctx context.Context, user *model.User) error

    // Get query builder
    Query() query.Builder[model.User, conv.User]
}
```

Access repositories through the client:

```go
client.UserRepo()     // User repository
client.PostRepo()     // Post repository
client.FollowsRepo()  // Edge repository
```

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
| `First(ctx)` | `(*Model, bool, error)` | First match (bool = exists) |
| `One(ctx)` | `(*Model, bool, error)` | Exactly one (errors if multiple) |
| `Count(ctx)` | `(int, error)` | Count of matches |
| `Exists(ctx)` | `(bool, error)` | Whether any exist |
| `Live(ctx)` | `(<-chan LiveResult, error)` | Real-time stream |

### Async Variants

Every method has an async version:

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
filter.User.Email.Contains("@gmail.com")

// Multiple conditions (AND)
client.UserRepo().Query().
    Where(
        filter.User.IsActive.IsTrue(),
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

1. **Define models** as Go structs with `som.Node` or `som.Edge`
2. **Run generator**: `som gen ./model ./gen/som`
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
