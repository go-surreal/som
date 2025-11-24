# Features

## Core Features

### Code Generation

Generate type-safe database access code from Go struct models. The generator produces:

- **Repositories** with Create, Read, Update, Delete, and Query methods
- **Query builders** with fluent, chainable API
- **Filter builders** with 50+ operations per field type
- **Sort builders** for ordering results
- **Converters** for model-to-database transformation

### Query Builder

Fluent API for building complex queries with compile-time type checking:

```go
users, _ := client.UserRepo().Query().
    Filter(where.User.Email.Contains("@example.com")).
    Filter(where.User.CreatedAt.After(lastMonth)).
    Order(by.User.Name.Asc()).
    Limit(10).
    All(ctx)
```

### Real-Time Queries

Subscribe to database changes with type-safe live queries:

```go
updates, _ := client.UserRepo().Query().
    Filter(where.User.Status.Equal("active")).
    Live(ctx)

for update := range updates {
    // update.Action: "CREATE", "UPDATE", "DELETE"
    // update.Data: *User (the changed record)
}
```

### Async Operations

All query methods have async variants for concurrent execution:

```go
result := client.UserRepo().Query().AllAsync(ctx)
// Do other work...
users := <-result.Val()
err := <-result.Err()
```

### Eager Loading (Fetch)

Load related records in a single query:

```go
users, _ := client.UserRepo().Query().
    Fetch(with.User.Groups...).
    All(ctx)
```

### Automatic Timestamps

Embed `som.Timestamps` for auto-managed `CreatedAt` and `UpdatedAt` fields:

```go
type User struct {
    som.Node
    som.Timestamps  // Adds CreatedAt, UpdatedAt
    Name string
}
```

## Supported Data Types

### Primitive Types

| Type | Filter Operations | Sort | Notes |
|------|-------------------|------|-------|
| `string` | 30+ operations | Yes | Contains, StartsWith, IsEmail, Slug, etc. |
| `int`, `int8`, `int16`, `int32`, `int64` | Numeric + comparison | Yes | Add, Sub, Mul, Div, Abs |
| `uint8`, `uint16`, `uint32` | Numeric + comparison | Yes | |
| `float32`, `float64` | Numeric + comparison | Yes | |
| `bool` | Equal, IsTrue, IsFalse | Yes | |
| `rune` | Numeric operations | Yes | Treated as int32 |
| `byte`, `[]byte` | Basic comparison | No | Binary data |

### Time Types

| Type | Filter Operations | Sort | CBOR Tag |
|------|-------------------|------|----------|
| `time.Time` | Before, After, Add, Sub, Floor, Round, Format | Yes | 12 |
| `time.Duration` | Before, After, Add, Sub | Yes | 14 |

### Special Types

| Type | Filter Operations | Sort | CBOR Tag |
|------|-------------------|------|----------|
| `uuid.UUID` | Equal, NotEqual, In, NotIn | No | 37 |
| `url.URL` | Equal, NotEqual | No | - |

### Custom Types

- `som.Enum` - String-based enumerations with type-safe constants

### Collections

- **Slices** of any supported type with special operations:
  - `Length()`, `Contains()`, `ContainsAll()`, `ContainsAny()`
  - `Empty()`, `NotEmpty()`, `Intersects()`, `Inside()`
- **Pointers** to any type with `IsNil()` / `IsNotNil()` checks

### Embedded Structs

Nest structs within models with full filter support on nested fields:

```go
type Address struct {
    City    string
    Country string
}

type User struct {
    som.Node
    Address Address
}

// Filter on nested fields
where.User.Address.City.Equal("Berlin")
```

## Query Features

### Execution Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `All(ctx)` | `([]*Model, error)` | All matching records |
| `First(ctx)` | `(*Model, bool, error)` | First match or nil |
| `One(ctx)` | `(*Model, bool, error)` | Exactly one match (errors if multiple) |
| `Count(ctx)` | `(int, error)` | Count of matches |
| `Exists(ctx)` | `(bool, error)` | Whether any match exists |
| `Live(ctx)` | `(<-chan LiveResult, error)` | Stream of changes |

### Query Modifiers

| Method | Description |
|--------|-------------|
| `Filter(...)` | Add WHERE conditions (AND) |
| `Order(...)` | Sort results |
| `OrderRandom()` | Random ordering |
| `Limit(n)` | Max results |
| `Offset(n)` | Skip results |
| `Fetch(...)` | Eager load relations |
| `Timeout(d)` | Query timeout |
| `Parallel(bool)` | Parallel execution |

## Filter Operations

### String Operations (30+)

- **Comparison**: Equal, NotEqual, In, NotIn
- **Pattern**: Contains, StartsWith, EndsWith, FuzzyMatch
- **Validation**: IsEmail, IsURL, IsDomain, IsIP, IsUUID, IsSemVer
- **Transform**: Lowercase, Uppercase, Trim, Slug, Reverse
- **Utility**: Len, Split, Words, Join, Concat, Replace, Slice

### Numeric Operations

- **Comparison**: Equal, NotEqual, LessThan, GreaterThan, etc.
- **Arithmetic**: Add, Sub, Mul, Div, Raise, Abs

### Combining Filters

```go
// All conditions (AND)
where.All(
    where.User.Status.Equal("active"),
    where.User.Age.GreaterThan(18),
)

// Any condition (OR)
where.Any(
    where.User.Role.Equal("admin"),
    where.User.Premium.IsTrue(),
)
```

## Current Limitations

### Unsupported Go Types

- `uint`, `uint64`, `uintptr` - SurrealDB limitations with large integers
- `complex64`, `complex128` - No SurrealDB equivalent
- `map` types - Planned for future release

### Other Limitations

- No automatic migrations (schema changes require manual handling)
- No transaction support across multiple operations
- Official SurrealDB Go client not supported (uses sdbc)
