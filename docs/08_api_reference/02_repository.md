# Repository API

Repositories provide CRUD operations and query access for each model type.

## Accessing a Repository

Each model gets a typed repository accessor on the client:

```go
userRepo := client.UserRepo()
postRepo := client.PostRepo()
followsRepo := client.FollowsRepo()  // Edge repository
```

## Repository Interface

Generated interface for each model:

```go
type UserRepo interface {
    Create(ctx context.Context, user *model.User) error
    CreateWithID(ctx context.Context, id string, user *model.User) error
    Read(ctx context.Context, id *som.ID) (*model.User, bool, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, user *model.User) error
    Refresh(ctx context.Context, user *model.User) error
    Query() query.Builder[model.User, conv.User]
}
```

## Create

Insert a new record with auto-generated ULID:

```go
user := &model.User{
    Name:  "John",
    Email: "john@example.com",
}

err := client.UserRepo().Create(ctx, user)
if err != nil {
    return err
}

// user.ID is populated after successful creation
fmt.Println("Created:", user.ID)  // user:01HQMV8K2P...
```

## CreateWithID

Insert with a specific ID:

```go
user := &model.User{
    Name:  "John",
    Email: "john@example.com",
}

err := client.UserRepo().CreateWithID(ctx, "john", user)
// Creates record with ID: user:john
```

## Read

Fetch a record by ID:

```go
func (r *UserRepo) Read(ctx context.Context, id *som.ID) (*model.User, bool, error)
```

Returns:
- `*model.User` - The record (nil if not found)
- `bool` - Whether the record exists
- `error` - Any error that occurred

```go
user, exists, err := client.UserRepo().Read(ctx, id)
if err != nil {
    return err  // Database error
}
if !exists {
    return errors.New("user not found")
}
fmt.Println("Found:", user.Name)
```

## Update

Modify an existing record:

```go
// The record must have a valid ID
user.Name = "Jane"
user.Email = "jane@example.com"

err := client.UserRepo().Update(ctx, user)
if err != nil {
    return err
}
```

The record must have a valid ID from a previous `Create`, `CreateWithID`, or `Read` operation.

## Delete

Remove a record:

```go
err := client.UserRepo().Delete(ctx, user)
if err != nil {
    return err
}
```

## Refresh

Reload a record from the database:

```go
// Refresh to get latest data
err := client.UserRepo().Refresh(ctx, user)
if err != nil {
    return err
}
// user now contains current database values
```

Useful when:
- Other processes may have modified the record
- You need to verify current state
- After timestamp fields update

## Query

Access the query builder for complex queries:

```go
query := client.UserRepo().Query()

// Chain methods
users, err := query.
    Filter(where.User.IsActive.IsTrue()).
    Order(by.User.Name.Asc()).
    Limit(10).
    All(ctx)
```

See [Query Builder API](03_query_builder.md) for full documentation.

## Edge Repository (Relate)

Edge repositories have an additional `Relate()` method:

```go
type FollowsRepo interface {
    // Standard CRUD methods...
    Create(ctx context.Context, follows *model.Follows) error
    // ...

    // Edge-specific: Relate builder
    Relate() *relate.Follows
}
```

Using Relate:

```go
follows := &model.Follows{
    Since: time.Now(),
}

err := client.FollowsRepo().Relate().
    From(alice).
    To(bob).
    Create(ctx, follows)
```

## Complete Example

```go
func UserService(ctx context.Context, client *som.Client) error {
    repo := client.UserRepo()

    // Create
    user := &model.User{
        Name:     "Alice",
        Email:    "alice@example.com",
        IsActive: true,
    }
    if err := repo.Create(ctx, user); err != nil {
        return fmt.Errorf("create: %w", err)
    }
    log.Printf("Created user: %s", user.ID)

    // Read
    found, exists, err := repo.Read(ctx, user.ID)
    if err != nil {
        return fmt.Errorf("read: %w", err)
    }
    if !exists {
        return errors.New("user not found after create")
    }
    log.Printf("Read user: %s", found.Name)

    // Update
    user.Name = "Alice Smith"
    if err := repo.Update(ctx, user); err != nil {
        return fmt.Errorf("update: %w", err)
    }
    log.Printf("Updated user")

    // Query
    activeUsers, err := repo.Query().
        Filter(where.User.IsActive.IsTrue()).
        All(ctx)
    if err != nil {
        return fmt.Errorf("query: %w", err)
    }
    log.Printf("Found %d active users", len(activeUsers))

    // Delete
    if err := repo.Delete(ctx, user); err != nil {
        return fmt.Errorf("delete: %w", err)
    }
    log.Printf("Deleted user")

    return nil
}
```

## Error Handling

Always check errors from repository operations:

```go
user, exists, err := client.UserRepo().Read(ctx, id)

// Check error first
if err != nil {
    // Database connection error, query error, etc.
    return fmt.Errorf("failed to read user: %w", err)
}

// Then check existence
if !exists {
    // Record simply doesn't exist (not an error)
    return ErrUserNotFound
}

// Use the record
fmt.Println(user.Name)
```
