# Nodes

Nodes are the primary way to define database records in SOM. A Node corresponds to a table in SurrealDB.

## Defining a Node

Embed `som.Node` in any struct to make it a database record:

```go
package model

import "github.com/go-surreal/som"

type User struct {
    som.Node

    Username string
    Email    string
    Age      int
    IsActive bool
}
```

This creates a `user` table in SurrealDB with columns for each field.

## The ID Field

The `som.Node` embedding provides an `ID` field of type `*som.ID`. You don't define it yourself:

```go
user := &model.User{Username: "john"}

// Create generates a ULID automatically
client.UserRepo().Create(ctx, user)
fmt.Println(user.ID)  // user:01HQMV8K2P...

// Or specify your own ID
client.UserRepo().CreateWithID(ctx, "john", user)
fmt.Println(user.ID)  // user:john
```

### ID Format

SurrealDB record IDs have the format `table:id`. SOM handles this automatically:
- `user:01HQMV8K2P...` - Auto-generated ULID
- `user:john` - Custom string ID
- `user:123` - Numeric ID

## Timestamps

Add automatic time tracking by embedding `som.Timestamps`:

```go
type User struct {
    som.Node
    som.Timestamps  // Adds CreatedAt and UpdatedAt

    Username string
}
```

This adds two fields:

| Field | Type | Behavior |
|-------|------|----------|
| `CreatedAt` | `time.Time` | Set on creation, **readonly** |
| `UpdatedAt` | `time.Time` | Updated on every modification |

These fields are managed by SurrealDB:

```surql
DEFINE FIELD created_at ON user VALUE $before OR time::now() READONLY;
DEFINE FIELD updated_at ON user VALUE time::now();
```

You cannot manually set these values - they're controlled by the database.

## Optimistic Locking

Prevent concurrent update conflicts by embedding `som.OptimisticLock`:

```go
type Document struct {
    som.Node
    som.OptimisticLock  // Adds version tracking

    Title   string
    Content string
}
```

This adds a hidden version field that:

- Starts at 1 for new records
- Increments on each successful update
- Throws an error if updating with a stale version

```go
// Detect conflicts
err := client.DocumentRepo().Update(ctx, staleDoc)
if errors.Is(err, som.ErrOptimisticLock) {
    // Another process updated this record
}
```

See [Optimistic Locking](05_optimistic_locking.md) for detailed documentation.

## Field Types

Nodes support all [Data Types](../data_types/README.md):

```go
type AllFields struct {
    som.Node

    // Primitives
    Name     string
    Age      int
    Score    float64
    Active   bool

    // Time
    Birthday  time.Time
    Duration  time.Duration

    // Special
    ProfileID uuid.UUID
    Website   url.URL

    // Collections
    Tags      []string
    Scores    []int

    // Optional (pointers)
    Nickname  *string
    DeletedAt *time.Time

    // Nested structs
    Address   Address
    Settings  *Settings
}
```

## Table Naming

By default, the table name is the lowercase struct name:

| Struct | Table |
|--------|-------|
| `User` | `user` |
| `BlogPost` | `blogpost` |
| `UserProfile` | `userprofile` |

## Repository Methods

For each Node, SOM generates a repository with these methods:

```go
type UserRepo interface {
    // Create with auto-generated ULID
    Create(ctx context.Context, user *model.User) error

    // Create with specific ID
    CreateWithID(ctx context.Context, id string, user *model.User) error

    // Read by ID - returns (model, exists, error)
    Read(ctx context.Context, id *som.ID) (*model.User, bool, error)

    // Update existing record (must have ID)
    Update(ctx context.Context, user *model.User) error

    // Delete record
    Delete(ctx context.Context, user *model.User) error

    // Refresh from database
    Refresh(ctx context.Context, user *model.User) error

    // Query builder
    Query() query.Builder[model.User, conv.User]
}
```

## Generated Filters

Each field generates type-safe filters in the `filter` package:

```go
// String fields
filter.User.Username.Equal("john")
filter.User.Username.Contains("oh")
filter.User.Email.EndsWith("@gmail.com")

// Numeric fields
filter.User.Age.GreaterThan(18)
filter.User.Age.LessThanOrEqual(65)

// Boolean fields
filter.User.Active.IsTrue()
filter.User.Active.IsFalse()

// Optional fields
filter.User.DeletedAt.IsNil()
filter.User.DeletedAt.IsNotNil()
```

## Generated Sorts

Each sortable field generates entries in the `by` package:

```go
// Sort ascending
by.User.Username.Asc()

// Sort descending
by.User.CreatedAt.Desc()
```

## Example: Complete Node

```go
package model

import (
    "time"
    "github.com/go-surreal/som"
    "github.com/google/uuid"
)

type User struct {
    som.Node
    som.Timestamps

    // Required fields
    Username string
    Email    string

    // Optional fields
    DisplayName *string
    AvatarURL   *string

    // Typed fields
    ExternalID  uuid.UUID
    LastLoginAt *time.Time

    // Nested data
    Profile     UserProfile
    Preferences *UserPreferences
}

type UserProfile struct {
    Bio      string
    Location string
    Website  *url.URL
}

type UserPreferences struct {
    Theme        string
    Notifications bool
}
```
