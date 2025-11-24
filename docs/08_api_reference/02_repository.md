# Repository API

Repositories provide CRUD operations for each model type.

## Accessing a Repository

```go
userRepo := client.UserRepo()
```

## Create

Insert a new record:

```go
user := &model.User{
    Name:  "John",
    Email: "john@example.com",
}

err := client.UserRepo().Create(ctx, user)
// user.ID is populated after successful creation
```

## Read

Fetch a record by ID:

```go
user, err := client.UserRepo().Read(ctx, id)
if err != nil {
    return err
}
```

### Handling Not Found

```go
user, err := client.UserRepo().Read(ctx, id)
if err != nil {
    // Check if record doesn't exist vs other errors
    return err
}
if user == nil {
    // Record not found
}
```

## Update

Modify an existing record:

```go
user.Name = "Jane"
err := client.UserRepo().Update(ctx, user)
```

The record must have a valid ID from a previous Create or Read operation.

## Delete

Remove a record by ID:

```go
err := client.UserRepo().Delete(ctx, id)
```

## Query

Access the query builder:

```go
query := client.UserRepo().Query()
```

See [Query Basics](../06_querying/01_basics.md) for query builder documentation.

## Batch Operations

For bulk operations, use queries with the appropriate methods:

```go
// Delete all inactive users
_, err := client.UserRepo().Query().
    Filter(where.User.IsActive.Equal(false)).
    Delete(ctx)
```

## Transaction Support

Wrap operations in transactions for atomicity:

```go
// Transaction support varies - check generated code
```
