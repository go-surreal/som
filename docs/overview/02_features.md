# Features

## Core Features

### Code Generation

Generate type-safe database access code from Go struct models. The generator produces:

- **Repositories** with Create, Read, Update, Delete, Insert, and Query methods
- **Query builders** with fluent, chainable API
- **Filter builders** with 50+ operations per field type
- **Sort builders** for ordering results
- **Converters** for model-to-database transformation

### Query Builder

Fluent API for building complex queries with compile-time type checking:

```go
users, _ := client.UserRepo().Query().
    Where(filter.User.Email.Contains("@example.com")).
    Where(filter.User.CreatedAt.After(lastMonth)).
    Order(by.User.Name.Asc()).
    Limit(10).
    All(ctx)
```

### Real-Time Queries

Subscribe to database changes with type-safe live queries:

```go
updates, _ := client.UserRepo().Query().
    Where(filter.User.Status.Equal("active")).
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
    som.Node[som.ULID]
    som.Timestamps  // Adds CreatedAt, UpdatedAt
    Name string
}
```

### Optimistic Locking

Embed `som.OptimisticLock` to prevent concurrent update conflicts:

```go
type Document struct {
    som.Node[som.ULID]
    som.OptimisticLock  // Adds version tracking

    Title string
}

// Concurrent updates are detected
err := client.DocumentRepo().Update(ctx, staleDoc)
if errors.Is(err, som.ErrOptimisticLock) {
    // Handle conflict
}
```

### Soft Delete

Embed `som.SoftDelete` for non-destructive deletion with automatic query filtering:

```go
type User struct {
    som.Node[som.ULID]
    som.SoftDelete  // Adds deleted_at, WithDeleted(), Restore(), Erase()
    Name string
}
```

### Iterator Methods

Process large result sets efficiently with Go range-over-func:

```go
for user, err := range client.UserRepo().Query().Iterate(ctx, 100) {
    if err != nil {
        break
    }
    process(user)
}
```

### Bulk Insert

Insert multiple records efficiently in a single operation:

```go
users := []*model.User{
    {Name: "Alice"},
    {Name: "Bob"},
    {Name: "Charlie"},
}
err := client.UserRepo().Insert(ctx, users)
```

### Lifecycle Hooks

Register hooks for pre/post CRUD operations:

```go
unregister := client.UserRepo().OnBeforeCreate(func(ctx context.Context, user *model.User) error {
    // Validate or transform before creation
    return nil
})
defer unregister()
```

### Full-Text Search

BM25-based relevance searching with highlighting and score sorting:

```go
results, _ := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang tutorial")).
    AllMatches(ctx)
```

### Context Cache

Optional request-scoped caching for Read operations:

```go
ctx, cleanup := som.WithCache[model.User](ctx, som.Lazy())
defer cleanup()
```

### Complex ID Types

Support for array and object-based record IDs for range-efficient queries:

```go
type WeatherKey struct {
    som.ArrayID
    City string
    Date time.Time
}

type Weather struct {
    som.Node[WeatherKey]
    Temperature float64
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
| `time.Month` | Comparable operations | Yes | - |
| `time.Weekday` | Comparable operations | Yes | - |

### Special Types

| Type | Filter Operations | Sort | Notes |
|------|-------------------|------|-------|
| `uuid.UUID` (google) | Equal, NotEqual, In, NotIn | Yes | CBOR Tag 37 |
| `uuid.UUID` (gofrs) | Equal, NotEqual, In, NotIn | Yes | CBOR Tag 37 |
| `url.URL` | Equal, NotEqual | Yes | - |
| `som.Email` | Equal, In, User(), Host() | Yes | - |
| `som.Password[A]` | Zero, IsNil (auto-hashed) | No | - |
| `som.SemVer` | Equal, Compare, Major, Minor, Patch | Yes | - |

### Geometry Types

SOM supports geometry types from three popular Go libraries:

- `github.com/paulmach/orb` - Point, LineString, Polygon, MultiPoint, etc.
- `github.com/peterstace/simplefeatures/geom` - Same geometry types
- `github.com/twpayne/go-geom` - Same geometry types

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
    som.Node[som.ULID]
    Address Address
}

// Filter on nested fields
filter.User.Address.City.Equal("Berlin")
```

## Query Features

### Execution Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `All(ctx)` | `([]*Model, error)` | All matching records |
| `First(ctx)` | `(*Model, error)` | First match or `ErrNotFound` |
| `Count(ctx)` | `(int, error)` | Count of matches |
| `Exists(ctx)` | `(bool, error)` | Whether any match exists |
| `Live(ctx)` | `(<-chan LiveResult, error)` | Stream of changes |

### Query Modifiers

| Method | Description |
|--------|-------------|
| `Where(...)` | Add WHERE conditions (AND) |
| `Order(...)` | Sort results |
| `OrderRandom()` | Random ordering |
| `Limit(n)` | Max results |
| `Start(n)` | Skip results |
| `Fetch(...)` | Eager load relations |
| `Timeout(d)` | Query timeout |
| `Parallel(bool)` | Parallel execution |
| `TempFiles(bool)` | Disk-based processing for large result sets |
| `WithDeleted()` | Include soft-deleted records |
| `Range(from, to)` | Range query for complex IDs |

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
filter.All(
    filter.User.Status.Equal("active"),
    filter.User.Age.GreaterThan(18),
)

// Any condition (OR)
filter.Any(
    filter.User.Role.Equal("admin"),
    filter.User.Premium.IsTrue(),
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
