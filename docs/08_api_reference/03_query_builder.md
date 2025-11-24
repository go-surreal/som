# Query Builder API

The query builder provides a fluent interface for constructing database queries.

## Method Reference

### Filter

Add conditions to narrow results:

```go
Query().Filter(conditions...)
```

Multiple conditions are combined with AND.

### OrderBy

Sort results:

```go
Query().OrderBy(order.Model.Field.Asc())
Query().OrderBy(order.Model.Field.Desc())
```

### Limit

Restrict result count:

```go
Query().Limit(10)
```

### Offset

Skip results for pagination:

```go
Query().Offset(20)
```

## Execution Methods

### All

Get all matching results:

```go
results, err := Query().All(ctx)
// Returns []*Model, error
```

### First

Get the first matching result:

```go
result, err := Query().First(ctx)
// Returns *Model, error
```

### Count

Get the count of matching results:

```go
count, err := Query().Count(ctx)
// Returns int, error
```

### Exists

Check if any results exist:

```go
exists, err := Query().Exists(ctx)
// Returns bool, error
```

### Live

Subscribe to real-time updates:

```go
channel, err := Query().Live(ctx)
// Returns <-chan LiveResult, error
```

## Chaining Example

```go
users, err := client.UserRepo().Query().
    Filter(
        where.User.IsActive.Equal(true),
        where.User.Age.GreaterThan(18),
    ).
    OrderBy(order.User.Name.Asc()).
    Limit(10).
    Offset(0).
    All(ctx)
```

## Filter Operations

| Operation | Description |
|-----------|-------------|
| `Equal(v)` | Equals value |
| `NotEqual(v)` | Not equals value |
| `GreaterThan(v)` | Greater than value |
| `GreaterThanOrEqual(v)` | Greater than or equal |
| `LessThan(v)` | Less than value |
| `LessThanOrEqual(v)` | Less than or equal |
| `Contains(v)` | String contains |
| `StartsWith(v)` | String starts with |
| `EndsWith(v)` | String ends with |
| `IsNull()` | Is null/nil |
| `IsNotNull()` | Is not null/nil |
