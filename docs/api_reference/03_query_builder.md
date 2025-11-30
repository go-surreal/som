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

### Search

Add full-text search conditions. Multiple conditions are ORed (any match):

```go
Query().Search(searches...)
```

```go
query.Search(where.Article.Content.Matches("golang tutorial"))
```

### SearchAll

Add full-text search conditions with AND semantics (all must match):

```go
Query().SearchAll(searches...)
```

```go
query.SearchAll(
    where.Article.Content.Matches("golang"),
    where.Article.Content.Matches("tutorial"),
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

### TempFiles

Enable temporary file-based query processing for large result sets:

```go
Query().TempFiles(enabled bool)
```

```go
// Process large result sets using temporary files instead of memory
users, err := client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    Limit(100000).
    TempFiles(true).
    All(ctx)
```

Note: TempFiles reduces memory usage at the cost of slower performance. Not available with Live queries.

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

### AllMatches

Get all search results with metadata (scores, highlights, offsets):

```go
func (b Builder) AllMatches(ctx context.Context) ([]SearchResult[Model], error)
```

```go
results, err := client.ArticleRepo().Query().
    Search(where.Article.Content.Matches("golang")).
    AllMatches(ctx)

for _, result := range results {
    fmt.Printf("Score: %f, Title: %s\n", result.Score(), result.Model.Title)
}
```

### FirstMatch

Get the first search result with metadata:

```go
func (b Builder) FirstMatch(ctx context.Context) (*SearchResult[Model], bool, error)
```

```go
result, found, err := client.ArticleRepo().Query().
    Search(where.Article.Content.Matches("golang")).
    FirstMatch(ctx)

if found {
    fmt.Printf("Best match: %s\n", result.Model.Title)
}
```

## Iterator Methods

For processing large result sets efficiently, use the iterator methods. These leverage Go 1.22+'s range-over-func feature to stream results in batches.

### Iterate

Iterate over all matching records in batches:

```go
func (b Builder) Iterate(ctx context.Context, batchSize int) iter.Seq2[*Model, error]
```

```go
// Process all active users in batches of 100
for user, err := range client.UserRepo().Query().
    Filter(where.User.IsActive.IsTrue()).
    Iterate(ctx, 100) {

    if err != nil {
        log.Fatal(err)
    }
    processUser(user)
}
```

### IterateID

Iterate over record IDs only (more efficient when you only need IDs):

```go
func (b Builder) IterateID(ctx context.Context, batchSize int) iter.Seq2[string, error]
```

```go
// Collect all user IDs
var ids []string
for id, err := range client.UserRepo().Query().IterateID(ctx, 500) {
    if err != nil {
        log.Fatal(err)
    }
    ids = append(ids, id)
}
```

### Early Termination

Iterators support breaking out early:

```go
// Find first 10 matching a condition
count := 0
for user, err := range client.UserRepo().Query().Iterate(ctx, 50) {
    if err != nil {
        break
    }
    if someCondition(user) {
        count++
        if count >= 10 {
            break  // Stop iteration
        }
    }
}
```

### When to Use Iterators

| Scenario | Method |
|----------|--------|
| Process all records without loading into memory | `Iterate` |
| Need only record IDs | `IterateID` |
| Fixed number of results | `Limit().All()` |
| Need random access to results | `All()` |

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

## SearchResult Type

Returned by `AllMatches()` and `FirstMatch()`:

```go
type SearchResult[M any] struct {
    Model      M                    // The matched model
    Scores     []float64            // BM25 relevance scores
    Highlights map[int]string       // Highlighted text by ref
    Offsets    map[int][]Offset     // Match positions by ref
}
```

### Helper Methods

```go
// Get the primary score (first in slice)
func (r SearchResult[M]) Score() float64

// Get highlighted text for a specific ref (defaults to 0)
func (r SearchResult[M]) Highlighted(ref ...int) string

// Get match offsets for a specific ref (defaults to 0)
func (r SearchResult[M]) Offset(ref ...int) []Offset
```

## Score Sorting

The `query` package provides score-based sorting for search queries:

```go
import "yourproject/gen/som/query"
```

### Basic Score Sort

```go
query.Score(0).Desc()      // Sort by score descending
query.Score(0).Asc()       // Sort by score ascending
```

### Multiple Refs

```go
query.Score(0, 1).Desc()   // Sort by combined scores
```

### Combination Modes

```go
query.Score(0, 1).Sum().Desc()              // Sum scores (default)
query.Score(0, 1).Max().Desc()              // Maximum score
query.Score(0, 1).Average().Desc()          // Average score
query.Score(0, 1).Weighted(2.0, 0.5).Desc() // Weighted combination
```

### Usage Example

```go
results, err := client.ArticleRepo().Query().
    Search(
        where.Article.Title.Matches("golang").Ref(0),
        where.Article.Content.Matches("golang").Ref(1),
    ).
    Order(query.Score(0, 1).Weighted(2.0, 1.0).Desc()).
    AllMatches(ctx)
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
