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
    Insert(ctx context.Context, users []*model.User) error
    Read(ctx context.Context, id string) (*model.User, bool, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, user *model.User) error
    Refresh(ctx context.Context, user *model.User) error
    Query() query.Builder[model.User]

    // Index access (e.g. per-index Rebuild)
    Index() *index.User

    // Lifecycle hooks
    OnBeforeCreate(fn func(ctx context.Context, node *model.User) error) func()
    OnAfterCreate(fn func(ctx context.Context, node *model.User) error) func()
    OnBeforeUpdate(fn func(ctx context.Context, node *model.User) error) func()
    OnAfterUpdate(fn func(ctx context.Context, node *model.User) error) func()
    OnBeforeDelete(fn func(ctx context.Context, node *model.User) error) func()
    OnAfterDelete(fn func(ctx context.Context, node *model.User) error) func()
}
```

## Create

Insert a new record with auto-generated ID:

```go
user := &model.User{
    Name:  "John",
    Email: "john@example.com",
}

err := client.UserRepo().Create(ctx, user)
if err != nil {
    return err
}

// user.ID() is populated after successful creation
fmt.Println("Created:", user.ID())  // 01HQMV8K2P...
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

## Insert

Bulk insert multiple records in a single operation:

```go
users := []*model.User{
    {Name: "Alice", Email: "alice@example.com"},
    {Name: "Bob", Email: "bob@example.com"},
    {Name: "Charlie", Email: "charlie@example.com"},
}

err := client.UserRepo().Insert(ctx, users)
```

## Read

Fetch a record by ID:

```go
func (r *UserRepo) Read(ctx context.Context, id string) (*model.User, bool, error)
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

## Index

Access the index manager for this table. Each index exposes a `Rebuild(ctx)` method:

```go
err := client.UserRepo().Index().Count().Rebuild(ctx)
```

## Query

Access the query builder for complex queries:

```go
query := client.UserRepo().Query()

// Chain methods
users, err := query.
    Where(filter.User.IsActive.IsTrue()).
    Order(by.User.Name.Asc()).
    Limit(10).
    All(ctx)
```

See [Query Builder API](03_query_builder.md) for full documentation.

## Lifecycle Hooks

Register callbacks that execute before or after CRUD operations:

```go
// Register a hook
unregister := client.UserRepo().OnBeforeCreate(func(ctx context.Context, user *model.User) error {
    // Validate or transform before creation
    if user.Email == "" {
        return errors.New("email is required")
    }
    return nil
})

// Unregister the hook when no longer needed
defer unregister()
```

Available hooks:

| Hook | When |
|------|------|
| `OnBeforeCreate` | Before a record is created |
| `OnAfterCreate` | After a record is created |
| `OnBeforeUpdate` | Before a record is updated |
| `OnAfterUpdate` | After a record is updated |
| `OnBeforeDelete` | Before a record is deleted |
| `OnAfterDelete` | After a record is deleted |

Each hook returns an unregister function. Call it to remove the hook.

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
    log.Printf("Created user: %s", user.ID())

    // Read
    found, exists, err := repo.Read(ctx, string(user.ID()))
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
        Where(filter.User.IsActive.IsTrue()).
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
