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
| Display page of results to user | `Limit()` + `Start()` |
| Process all records in background job | `Iterate()` |
| Export all data | `Iterate()` |
| Build list of IDs for batch operation | `IterateID()` |
| Random access to results | `All()` |

Iterators automatically handle batching internally and support early termination via `break`.
