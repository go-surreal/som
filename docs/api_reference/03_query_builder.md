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
Query().Where(conditions...)
```

```go
query.Where(
    filter.User.IsActive.True(),
    filter.User.Age.GreaterThan(18),
)
```

### Search

Add full-text search conditions. Multiple conditions are ORed (any match):

```go
Query().Search(searches...)
```

```go
query.Search(filter.Article.Content.Matches("golang tutorial"))
```

### SearchAll

Add full-text search conditions with AND semantics (all must match):

```go
Query().SearchAll(searches...)
```

```go
query.SearchAll(
    filter.Article.Content.Matches("golang"),
    filter.Article.Content.Matches("tutorial"),
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

### Start

Skip first n results (for pagination):

```go
Query().Start(n int)
```

### Fetch

Eager load related records:

```go
Query().Fetch(relations...)
```

```go
query.Fetch(with.User.Groups())

// Nested relations chain
query.Fetch(with.User.Organization().Owner())
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
    Where(filter.User.IsActive.True()).
    Limit(100000).
    TempFiles(true).
    All(ctx)
```

Note: TempFiles reduces memory usage at the cost of slower performance. Not available with Live queries.

### WithDeleted

Include soft-deleted records in results (only for models with `som.SoftDelete`):

```go
Query().WithDeleted()
```

### WithExpired

Include records past their `expires_at` that have not been purged yet (only for models with
`som.Expiry`):

```go
Query().WithExpired()
```

### Range

Range query over record IDs, for models with string IDs (ULID, UUID, Rand) or complex IDs
(ArrayID, ObjectID). It uses SurrealDB's native range syntax and avoids a table scan:

```go
Query().Range(from som.RangeFrom, to som.RangeTo)
```

Range boundary constructors:

| Constructor | Description |
|-------------|-------------|
| `som.From[T](val)` | Inclusive start bound |
| `som.FromExclusive[T](val)` | Exclusive start bound |
| `som.FromStart()` | Open-ended start (no lower bound) |
| `som.To[T](val)` | Exclusive end bound |
| `som.ToInclusive[T](val)` | Inclusive end bound |
| `som.ToEnd()` | Open-ended end (no upper bound) |

String ID example:

```go
results, err := client.UserRepo().Query().Range(
    som.From(som.ULID("01HY5E8ZQA1BCD2EF3GH4JK5MN")),
    som.To(som.ULID("01HY5E8ZQA9ZZZ9ZZ9ZZ9ZZ9ZZ")),
).All(ctx)
```

Open-ended range:

```go
results, err := client.UserRepo().Query().Range(
    som.From(som.ULID("01HY5E8ZQA1BCD2EF3GH4JK5MN")),
    som.ToEnd(),
).All(ctx)
```

Complex ID example:

```go
results, err := client.WeatherRepo().Query().Range(
    som.From(model.WeatherKey{City: "London", Date: start}),
    som.To(model.WeatherKey{City: "London", Date: end}),
).All(ctx)
```

Note: `Range()` returns a `BuilderNoLive` — live queries are not supported with range queries.

## Execution Methods

These methods execute the query and return results.

### All

Get all matching records:

```go
func (b Builder) All(ctx context.Context) ([]*Model, error)
```

```go
users, err := client.UserRepo().Query().
    Where(filter.User.IsActive.True()).
    All(ctx)
```

### First

Get the first matching record. Returns `ErrNotFound` if no match:

```go
func (b Builder) First(ctx context.Context) (*Model, error)
```

```go
user, err := client.UserRepo().Query().
    Where(filter.User.Email.Equal("john@example.com")).
    First(ctx)

if err != nil {
    if errors.Is(err, som.ErrNotFound) {
        // No matching record
    }
    return err
}
fmt.Println(user.Name)
```

### Count

Get count of matching records:

```go
func (b Builder) Count(ctx context.Context) (int, error)
```

```go
count, err := client.UserRepo().Query().
    Where(filter.User.IsActive.True()).
    Count(ctx)
```

### Exists

Check if any matching records exist:

```go
func (b Builder) Exists(ctx context.Context) (bool, error)
```

```go
exists, err := client.UserRepo().Query().
    Where(filter.User.Email.Equal("john@example.com")).
    Exists(ctx)
```

### Live

Subscribe to real-time updates:

```go
func (b Builder[M]) Live(ctx context.Context) (<-chan LiveResult[*M], error)
```

`LiveResult` is an interface; each event is one of `LiveCreate`, `LiveUpdate`, `LiveDelete` or
`LiveKilled`:

```go
updates, err := client.UserRepo().Query().
    Where(filter.User.IsActive.True()).
    Live(ctx)

for update := range updates {
    switch res := update.(type) {
    case query.LiveCreate[*model.User]:
        user, err := res.Get()
        ...
    case query.LiveKilled[*model.User]:
        return
    }
}
```

See [Live Queries](../querying/04_live_queries.md).

### LiveCount

Track the number of matching records in real time:

```go
func (b Builder[M]) LiveCount(ctx context.Context) (<-chan int, error)
```

### AllIDs / FirstID

Fetch only record IDs:

```go
ids, err := client.UserRepo().Query().AllIDs(ctx)
id, err := client.UserRepo().Query().FirstID(ctx)
```

### Paginate

Cursor-based (keyset) pagination:

```go
page, err := client.UserRepo().Query().
    Order(by.User.CreatedAt.Desc()).
    Paginate().First(20).Get(ctx)
```

See [Ordering & Pagination](../querying/03_ordering_pagination.md).

### Distinct

Distinct values of a single field (package-level function):

```go
values, err := query.Distinct(ctx, client.UserRepo().Query(), field.User.Country)
```

See [Distinct Values](../querying/09_distinct.md).

### Changes

For models with a changefeed, the repository exposes a change query builder:

```go
entries, err := client.UserRepo().Changes().SinceVersionstamp(0).Show(ctx)
```

See [Changefeed](../querying/08_changefeed.md).

### Describe / Debug

Inspect the generated statement:

```go
stmt := client.UserRepo().Query().Where(...).Describe()
stmt = client.UserRepo().Query().Where(...).DescribeWithVars()

// Debug prints the statement and keeps the builder chainable
users, err := client.UserRepo().Query().Debug("users").All(ctx)
```

### AllMatches

Get all search results with metadata (scores, highlights, offsets):

```go
func (b Builder) AllMatches(ctx context.Context) ([]SearchResult[Model], error)
```

```go
results, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang")).
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
    Search(filter.Article.Content.Matches("golang")).
    FirstMatch(ctx)

if found {
    fmt.Printf("Best match: %s\n", result.Model.Title)
}
```

## Iterator Methods

For processing large result sets efficiently, use the iterator methods. These leverage Go's range-over-func feature to stream results in batches.

### Iterate

Iterate over all matching records in batches:

```go
func (b Builder) Iterate(ctx context.Context, batchSize int) iter.Seq2[*Model, error]
```

```go
// Process all active users in batches of 100
for user, err := range client.UserRepo().Query().
    Where(filter.User.IsActive.True()).
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
| `AllIDs(ctx)` | `AllIDsAsync(ctx)` |
| `First(ctx)` | `FirstAsync(ctx)` |
| `FirstID(ctx)` | `FirstIDAsync(ctx)` |
| `Count(ctx)` | `CountAsync(ctx)` |
| `Exists(ctx)` | `ExistsAsync(ctx)` |
| `Paginate().Get(ctx)` | `Paginate().GetAsync(ctx)` |

Live queries have no async variant — they already return a channel.

### Using Async Methods

```go
// Start query in background
result := client.UserRepo().Query().
    Where(filter.User.IsActive.True()).
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

## LiveResult Types

```go
type LiveResult[M any] interface{ ... }

type LiveCreate[M any] interface{ Get() (M, error) }
type LiveUpdate[M any] interface{ Get() (M, error) }
type LiveDelete[M any] interface{ Get() (M, error) }
type LiveKilled[M any] interface{ ... }  // server terminated the subscription
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
        filter.Article.Title.Matches("golang").Ref(0),
        filter.Article.Content.Matches("golang").Ref(1),
    ).
    Order(query.Score(0, 1).Weighted(2.0, 1.0).Desc()).
    AllMatches(ctx)
```

## Complete Example

```go
// Complex query with all features
users, err := client.UserRepo().Query().
    // Filter conditions
    Where(
        filter.User.IsActive.True(),
        filter.User.Age.GreaterThanEqual(18),
        filter.Any(
            filter.User.Role.Equal("admin"),
            filter.User.Role.Equal("moderator"),
        ),
    ).
    // Sorting
    Order(
        by.User.CreatedAt.Desc(),
        by.User.Name.Asc(),
    ).
    // Pagination
    Limit(20).
    Start(0).
    // Eager loading
    Fetch(with.User.Posts()).
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
        Where(filter.User.IsActive.True()).
        Order(by.User.CreatedAt.Desc()).
        Limit(pageSize).
        Start((page - 1) * pageSize).
        All(ctx)
}

// Get total for pagination UI
func GetTotal(ctx context.Context) (int, error) {
    return client.UserRepo().Query().
        Where(filter.User.IsActive.True()).
        Count(ctx)
}
```

## Query Reuse

Queries can be built incrementally:

```go
// Base query
baseQuery := client.UserRepo().Query().
    Where(filter.User.IsActive.True())

// Different executions
count, _ := baseQuery.Count(ctx)
first, _ := baseQuery.First(ctx)
all, _ := baseQuery.Limit(10).All(ctx)
```

Note: Each execution creates a new query based on the builder state at that point.
