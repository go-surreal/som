# Raw Queries

For cases where the generated query builder doesn't cover your needs, SOM provides raw query execution with parameter binding.

## Executing a Raw Query

Use the `Raw` method on the generated client:

```go
result, err := client.Raw(ctx, "SELECT * FROM user WHERE age > $min_age", som.Params{
    "min_age": 18,
})
if err != nil {
    return err
}
```

## Scanning Results

### Multiple Rows

Use `Scan` to unmarshal the result set into a slice:

```go
var users []map[string]any
err := result.Scan(&users)
```

### Single Row

Use `ScanOne` to unmarshal the first row. Returns `som.ErrNotFound` if the result set is empty:

```go
var user map[string]any
err := result.ScanOne(&user)

if errors.Is(err, som.ErrNotFound) {
    // No matching record
}
```

## Parameter Binding

Always use parameterized queries to prevent injection:

```go
result, err := client.Raw(ctx,
    "SELECT * FROM user WHERE name = $name AND status = $status",
    som.Params{
        "name":   "Alice",
        "status": "active",
    },
)
```

The `som.Params` type is `map[string]any`.

## Multi-Statement Queries

When executing multiple statements, only the first statement's result set is returned:

```go
result, err := client.Raw(ctx,
    "LET $active = (SELECT * FROM user WHERE active = true); SELECT * FROM $active WHERE age > $min;",
    som.Params{"min": 25},
)
```

## Typed Results

You can scan into typed structs:

```go
type UserSummary struct {
    Name  string `json:"name"`
    Count int    `json:"count"`
}

var summaries []UserSummary
err := result.Scan(&summaries)
```

## When to Use Raw Queries

Raw queries are useful when you need:
- SurrealQL features not yet covered by the query builder
- Complex aggregations or subqueries
- Administrative commands
- Custom SurrealQL functions
