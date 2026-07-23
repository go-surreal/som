# Ordering & Pagination

Control the order and size of query results.

## Ordering

Use the generated `by` package:

```go
import "yourproject/gen/som/by"

users, err := client.UserRepo().Query().
    Order(by.User.Name.Asc()).
    All(ctx)
```

### Ascending Order

```go
Order(by.User.CreatedAt.Asc())
```

### Descending Order

```go
Order(by.User.CreatedAt.Desc())
```

### Random Order

```go
client.UserRepo().Query().OrderRandom().All(ctx)
```

### Multiple Order Clauses

```go
users, err := client.UserRepo().Query().
    Order(
        by.User.LastName.Asc(),
        by.User.FirstName.Asc(),
    ).
    All(ctx)
```

## Pagination

### Limit

Restrict the number of results:

```go
users, err := client.UserRepo().Query().
    Limit(10).
    All(ctx)
```

### Start

Skip a number of results (for pagination):

```go
// Page 2 with 10 items per page
users, err := client.UserRepo().Query().
    Order(by.User.CreatedAt.Desc()).
    Limit(10).
    Start(10).
    All(ctx)
```

### Pagination Helper

Combine limit and start for page-based pagination:

```go
func GetUsersPage(ctx context.Context, page, pageSize int) ([]*model.User, error) {
    return client.UserRepo().Query().
        Order(by.User.CreatedAt.Desc()).
        Limit(pageSize).
        Start((page - 1) * pageSize).
        All(ctx)
}
```

### Cursor-based Pagination

`Start()`/`Limit()` is offset pagination: simple, but it re-scans skipped rows on
every page and can skip or duplicate rows when the underlying data changes between
page loads. For stable paging over large or live datasets, use `Paginate()`, which
uses keyset (cursor) pagination.

A cursor is an opaque string that encodes a position in the result set. You never
build one by hand — read it from the returned page and pass it back to fetch the
next or previous page.

```go
import "yourproject/gen/som/query"

page, err := client.UserRepo().Query().
    Order(by.User.CreatedAt.Desc()).
    Paginate(ctx, query.First(20))
if err != nil {
    return err
}

for _, u := range page.Items {
    fmt.Println(u.Name)
}
```

The returned `*query.Page[User]` holds:

```go
page.Items       // []*User — the records for this page
page.Entries     // []query.Entry[User] — {Node *User, Cursor string} per item
page.PageInfo    // navigation metadata (see below)
page.TotalCount  // int; -1 unless WithTotalCount() was passed
```

`PageInfo` follows the GraphQL Relay convention:

```go
page.PageInfo.HasNextPage      // more items after this page
page.PageInfo.HasPreviousPage  // more items before this page
page.PageInfo.StartCursor      // cursor of the first item
page.PageInfo.EndCursor        // cursor of the last item
```

#### Options

| Option | Meaning |
|--------|---------|
| `query.First(n)` | Forward pagination, first `n` items |
| `query.After(cursor)` | Continue forward, after the given cursor |
| `query.Last(n)` | Backward pagination, last `n` items |
| `query.Before(cursor)` | Continue backward, before the given cursor |
| `query.WithTotalCount()` | Also run a COUNT query, exposed as `page.TotalCount` |
| `query.WithAccuratePageInfo()` | Extra query for exact `HasPreviousPage`/`HasNextPage` on boundary pages |

Use `First` with optional `After` for forward paging, or `Last` with optional
`Before` for backward paging. `First`+`Last` (or `After`+`Before`) together is an
error, as is omitting both `First` and `Last`.

#### Forward and backward

```go
// Next page: pass the previous page's end cursor.
next, err := client.UserRepo().Query().
    Order(by.User.CreatedAt.Desc()).
    Paginate(ctx, query.First(20), query.After(page.PageInfo.EndCursor))

// Previous page: pass the current page's start cursor.
prev, err := client.UserRepo().Query().
    Order(by.User.CreatedAt.Desc()).
    Paginate(ctx, query.Last(20), query.Before(page.PageInfo.StartCursor))
```

> The same `Order(...)` and `Where(...)` must be repeated on every page call.
> The cursor only encodes the position, not the query shape — changing the sort
> keys between pages produces meaningless results.

#### Stateless (frontend-driven) pagination

Because the cursor is an opaque string, it round-trips through the client. Each
HTTP request is independent — no server-side state:

```go
func handleUsers(w http.ResponseWriter, r *http.Request) {
    opts := []query.PaginateOption{query.First(20)}
    if after := r.URL.Query().Get("after"); after != "" {
        opts = append(opts, query.After(after))
    }

    page, err := client.UserRepo().Query().
        Order(by.User.CreatedAt.Desc()).
        Paginate(r.Context(), opts...)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]any{
        "items":     page.Items,
        "endCursor": page.PageInfo.EndCursor, // client echoes this back as ?after=
        "hasNext":   page.PageInfo.HasNextPage,
    })
}
```

#### In-process loop

For batch jobs that walk the whole set in one call, hold the query in a variable
and thread the cursor:

```go
q := client.UserRepo().Query().Order(by.User.CreatedAt.Desc())

var after string
for {
    opts := []query.PaginateOption{query.First(500)}
    if after != "" {
        opts = append(opts, query.After(after))
    }
    page, err := q.Paginate(ctx, opts...)
    if err != nil {
        return err
    }
    process(page.Items)
    if !page.PageInfo.HasNextPage {
        break
    }
    after = page.PageInfo.EndCursor
}
```

> For pure "process everything" jobs, `Iterate()` (below) is simpler — reach for
> cursor pagination when you also need the cursors, `PageInfo`, or stable paging
> exposed to a client.

#### Multi-field sorting

Cursor pagination handles compound sort keys automatically. An `id` tiebreaker is
always appended so the order is stable even when sort values tie:

```go
page, err := client.UserRepo().Query().
    Where(filter.User.Age.GreaterThan(18)).
    Order(by.User.LastName.Asc(), by.User.Age.Desc()).
    Paginate(ctx, query.First(10))
```

#### Limitations

- **String-ID models only.** Models with a complex ID (`ArrayID`/`ObjectID`) have
  no single `id` tiebreaker; `Paginate()` returns an error for them. Use
  `Range()` instead.
- **Top-level scalar sort fields only.** Sorting on nested fields
  (`group.name`) is not supported for cursor pagination.
- **Not compatible with `OrderRandom()`.**

## Combining With Filters

```go
users, err := client.UserRepo().Query().
    Where(filter.User.IsActive.Equal(true)).
    Order(by.User.Name.Asc()).
    Limit(20).
    All(ctx)
```

## Getting Total Count

For pagination UI, get the total count alongside results:

```go
query := client.UserRepo().Query().
    Where(filter.User.IsActive.Equal(true))

total, _ := query.Count(ctx)
users, _ := query.Limit(10).Start(0).All(ctx)
```

## Iterating Large Result Sets

For processing large datasets without loading everything into memory, use the iterator methods:

### Iterate

Stream records in batches:

```go
for user, err := range client.UserRepo().Query().
    Where(filter.User.IsActive.Equal(true)).
    Order(by.User.CreatedAt.Desc()).
    Iterate(ctx, 100) {  // Batch size of 100

    if err != nil {
        log.Fatal(err)
    }
    processUser(user)
}
```

### IterateID

Stream only record IDs (more efficient when you just need IDs):

```go
for id, err := range client.UserRepo().Query().IterateID(ctx, 500) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(id)
}
```

### When to Use Iterators vs Pagination

| Use Case | Recommended Approach |
|----------|---------------------|
| Stable / infinite-scroll paging for a client | `Paginate()` (cursor) |
| Simple numbered pages over small, stable data | `Limit()` + `Start()` |
| Process all records in background job | `Iterate()` |
| Export all data | `Iterate()` |
| Build list of IDs for batch operation | `IterateID()` |
| Random access to results | `All()` |

Iterators automatically handle batching internally and support early termination via `break`.
