# Ordering & Pagination

Control the order and size of query results.

## Ordering

Use the generated `order` package:

```go
import "yourproject/gen/som/order"

users, err := client.UserRepo().Query().
    OrderBy(order.User.Name.Asc()).
    All(ctx)
```

### Ascending Order

```go
OrderBy(order.User.CreatedAt.Asc())
```

### Descending Order

```go
OrderBy(order.User.CreatedAt.Desc())
```

### Multiple Order Clauses

```go
users, err := client.UserRepo().Query().
    OrderBy(
        order.User.LastName.Asc(),
        order.User.FirstName.Asc(),
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

### Offset

Skip a number of results (for pagination):

```go
// Page 2 with 10 items per page
users, err := client.UserRepo().Query().
    OrderBy(order.User.CreatedAt.Desc()).
    Limit(10).
    Offset(10).
    All(ctx)
```

### Pagination Helper

Combine limit and offset for page-based pagination:

```go
func GetUsersPage(ctx context.Context, page, pageSize int) ([]*model.User, error) {
    return client.UserRepo().Query().
        OrderBy(order.User.CreatedAt.Desc()).
        Limit(pageSize).
        Offset((page - 1) * pageSize).
        All(ctx)
}
```

## Combining With Filters

```go
users, err := client.UserRepo().Query().
    Filter(where.User.IsActive.Equal(true)).
    OrderBy(order.User.Name.Asc()).
    Limit(20).
    All(ctx)
```

## Getting Total Count

For pagination UI, get the total count alongside results:

```go
query := client.UserRepo().Query().
    Filter(where.User.IsActive.Equal(true))

total, _ := query.Count(ctx)
users, _ := query.Limit(10).Offset(0).All(ctx)
```

## Iterating Large Result Sets

For processing large datasets without loading everything into memory, use the iterator methods:

### Iterate

Stream records in batches:

```go
for user, err := range client.UserRepo().Query().
    Filter(where.User.IsActive.Equal(true)).
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
| Display page of results to user | `Limit()` + `Offset()` |
| Process all records in background job | `Iterate()` |
| Export all data | `Iterate()` |
| Build list of IDs for batch operation | `IterateID()` |
| Random access to results | `All()` |

Iterators automatically handle batching internally and support early termination via `break`.
