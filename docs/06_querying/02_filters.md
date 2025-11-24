# Filters

Filters narrow down query results using type-safe conditions.

## Basic Filtering

Use the generated `where` package:

```go
import "yourproject/gen/som/where"

users, err := client.UserRepo().Query().
    Filter(where.User.Email.Equal("john@example.com")).
    All(ctx)
```

## Comparison Operators

### Equality

```go
where.User.Name.Equal("John")
where.User.Name.NotEqual("John")
```

### Numeric Comparisons

```go
where.User.Age.GreaterThan(18)
where.User.Age.GreaterThanOrEqual(18)
where.User.Age.LessThan(65)
where.User.Age.LessThanOrEqual(65)
```

### String Operations

```go
where.User.Name.Contains("john")
where.User.Email.StartsWith("admin")
where.User.Email.EndsWith("@example.com")
```

### Null Checks

```go
where.User.DeletedAt.IsNull()
where.User.DeletedAt.IsNotNull()
```

## Multiple Filters

Multiple filters are combined with AND:

```go
users, err := client.UserRepo().Query().
    Filter(
        where.User.IsActive.Equal(true),
        where.User.Age.GreaterThan(18),
        where.User.Email.Contains("@company.com"),
    ).
    All(ctx)
```

## Nested Field Filters

Filter on embedded struct fields:

```go
users, err := client.UserRepo().Query().
    Filter(where.User.Address.City.Equal("New York")).
    All(ctx)
```

## Enum Filters

Filter by enum values:

```go
users, err := client.UserRepo().Query().
    Filter(where.User.Status.Equal(model.StatusActive)).
    All(ctx)
```

## Time Filters

```go
users, err := client.UserRepo().Query().
    Filter(where.User.CreatedAt.GreaterThan(time.Now().AddDate(0, -1, 0))).
    All(ctx)
```
