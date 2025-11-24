# Query Basics

SOM generates a type-safe query builder for each model. This page covers the fundamentals.

## Starting a Query

Access the query builder through the repository:

```go
query := client.UserRepo().Query()
```

## Executing Queries

### Get All Results

```go
users, err := client.UserRepo().Query().All(ctx)
```

### Get First Result

```go
user, err := client.UserRepo().Query().First(ctx)
```

### Check Existence

```go
exists, err := client.UserRepo().Query().Exists(ctx)
```

### Count Results

```go
count, err := client.UserRepo().Query().Count(ctx)
```

## Chaining Methods

Query methods return the query builder for chaining:

```go
users, err := client.UserRepo().Query().
    Filter(where.User.IsActive.Equal(true)).
    OrderBy(order.User.CreatedAt.Desc()).
    Limit(10).
    All(ctx)
```

## Basic CRUD Operations

While the query builder handles complex queries, repositories provide simple CRUD:

```go
// Create
err := client.UserRepo().Create(ctx, user)

// Read by ID
user, err := client.UserRepo().Read(ctx, id)

// Update
err := client.UserRepo().Update(ctx, user)

// Delete
err := client.UserRepo().Delete(ctx, id)
```

## Query Reuse

Build queries incrementally:

```go
baseQuery := client.UserRepo().Query().
    Filter(where.User.IsActive.Equal(true))

// Reuse for different operations
count, _ := baseQuery.Count(ctx)
first, _ := baseQuery.First(ctx)
all, _ := baseQuery.All(ctx)
```
