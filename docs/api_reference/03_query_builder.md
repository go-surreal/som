# Query Builder API

The query builder provides a fluent interface for constructing database queries. It's generated for each model and provides compile-time type safety.

## Getting a Query Builder

Access through the repository:

```go
query := client.UserRepo().Query()
```

## Builder Methods (Chainable)

All builder methods return the builder for chaining.

### Filter

Add WHERE conditions. Multiple filters are ANDed together:

```go
Query().Filter(conditions...)
```

```go
query.Filter(
    where.User.IsActive.IsTrue(),
    where.User.Age.GreaterThan(18),
)
```

### Order

Sort results by one or more fields:

```go
Query().Order(sorts...)
```

```go
query.Order(by.User.Name.Asc())
query.Order(by.User.CreatedAt.Desc(), by.User.Name.Asc())
```

### OrderRandom

Sort results randomly:

```go
Query().OrderRandom()
```

### Limit

Restrict maximum number of results:

```go
Query().Limit(n int)
```

### Offset

Skip first n results (for pagination):

```go
Query().Offset(n int)
```

### Fetch

Eager load related records:

```go
Query().Fetch(relations...)
```

```go
query.Fetch(with.User.Groups...)
```

### Timeout

Set query execution timeout:

```go
Query().Timeout(d time.Duration)
```

### Parallel

Enable parallel query execution:

```go
Query().Parallel(enabled bool)
```

## Execution Methods

These methods execute the query and return results.

### All

Get all matching records:

```go
func (b Builder) All(ctx context.Context) ([]*Model, error)
```

```go
users, err := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    All(ctx)
```

### First

Get the first matching record:

```go
func (b Builder) First(ctx context.Context) (*Model, bool, error)
```

Returns:
- `*Model` - The record (or nil if not found)
- `bool` - Whether a record was found
- `error` - Any error that occurred

```go
user, exists, err := client.UserRepo().Query().
    Filter(where.User.Email.Equal("john@example.com")).
    First(ctx)

if exists {
    fmt.Println(user.Name)
}
```

### One

Get exactly one matching record. Errors if multiple exist:

```go
func (b Builder) One(ctx context.Context) (*Model, bool, error)
```

```go
user, exists, err := client.UserRepo().Query().
    Filter(where.User.Email.Equal("john@example.com")).
    One(ctx)
```

### Count

Get count of matching records:

```go
func (b Builder) Count(ctx context.Context) (int, error)
```

```go
count, err := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    Count(ctx)
```

### Exists

Check if any matching records exist:

```go
func (b Builder) Exists(ctx context.Context) (bool, error)
```

```go
exists, err := client.UserRepo().Query().
    Filter(where.User.Email.Equal("john@example.com")).
    Exists(ctx)
```

### Live

Subscribe to real-time updates:

```go
func (b Builder) Live(ctx context.Context) (<-chan LiveResult[Model], error)
```

```go
updates, err := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    Live(ctx)

for update := range updates {
    if update.Error != nil {
        log.Println("Error:", update.Error)
        continue
    }
    fmt.Printf("Action: %s, Data: %+v\n", update.Action, update.Data)
}
```

## Async Methods

Every execution method has an async variant that returns immediately:

| Sync | Async |
|------|-------|
| `All(ctx)` | `AllAsync(ctx)` |
| `First(ctx)` | `FirstAsync(ctx)` |
| `One(ctx)` | `OneAsync(ctx)` |
| `Count(ctx)` | `CountAsync(ctx)` |
| `Exists(ctx)` | `ExistsAsync(ctx)` |
| `Live(ctx)` | `LiveAsync(ctx)` |

### Using Async Methods

```go
// Start query in background
result := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    AllAsync(ctx)

// Do other work...
doOtherWork()

// Get results when needed
users := <-result.Val()
err := <-result.Err()
```

### Async Result Type

```go
type asyncResult[T any] struct {
    val chan T
    err chan error
}

func (r *asyncResult[T]) Val() <-chan T
func (r *asyncResult[T]) Err() <-chan error
```

## LiveResult Type

```go
type LiveResult[M any] struct {
    Action string  // "CREATE", "UPDATE", or "DELETE"
    Data   *M      // The affected record
    Error  error   // Any error
}
```

## Complete Example

```go
// Complex query with all features
users, err := client.UserRepo().Query().
    // Filter conditions
    Filter(
        where.User.IsActive.IsTrue(),
        where.User.Age.GreaterThanOrEqual(18),
        where.Any(
            where.User.Role.Equal("admin"),
            where.User.Role.Equal("moderator"),
        ),
    ).
    // Sorting
    Order(
        by.User.CreatedAt.Desc(),
        by.User.Name.Asc(),
    ).
    // Pagination
    Limit(20).
    Offset(0).
    // Eager loading
    Fetch(with.User.Posts...).
    // Execution options
    Timeout(5 * time.Second).
    Parallel(true).
    // Execute
    All(ctx)
```

## Pagination Helper

```go
func GetPage(ctx context.Context, page, pageSize int) ([]*model.User, error) {
    return client.UserRepo().Query().
        Filter(where.User.IsActive.IsTrue()).
        Order(by.User.CreatedAt.Desc()).
        Limit(pageSize).
        Offset((page - 1) * pageSize).
        All(ctx)
}

// Get total for pagination UI
func GetTotal(ctx context.Context) (int, error) {
    return client.UserRepo().Query().
        Filter(where.User.IsActive.IsTrue()).
        Count(ctx)
}
```

## Query Reuse

Queries can be built incrementally:

```go
// Base query
baseQuery := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue())

// Different executions
count, _ := baseQuery.Count(ctx)
first, _, _ := baseQuery.First(ctx)
all, _ := baseQuery.Limit(10).All(ctx)
```

Note: Each execution creates a new query based on the builder state at that point.
