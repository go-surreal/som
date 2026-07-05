# Nodes

Nodes are the primary way to define database records in SOM. A Node corresponds to a table in SurrealDB.

## Defining a Node

Embed `som.Node[T]` in any struct to make it a database record. The type parameter `T` determines the ID format:

```go
package model

import "yourproject/gen/som"

type User struct {
    som.Node[som.ULID]

    Username string
    Email    string
    Age      int
    IsActive bool
}
```

This creates a `user` table in SurrealDB with columns for each field.

## The ID Field

The `som.Node[T]` embedding provides an `ID()` method. You don't define it yourself:

```go
user := &model.User{Username: "john"}

// Create generates a ULID automatically
client.UserRepo().Create(ctx, user)
fmt.Println(user.ID())  // 01HQMV8K2P...

// Or specify your own ID
client.UserRepo().CreateWithID(ctx, "john", user)
```

### ID Types

| Type | Description | Example |
|------|-------------|---------|
| `som.ULID` | ULID-based IDs (recommended) | `01HQMV8K2P...` |
| `som.UUID` | UUID-based IDs | `550e8400-e29b-...` |
| `som.Rand` | Random string IDs | `abc123def` |
| Custom struct | Complex array or object IDs | See [Complex IDs](#complex-id-types) |

### Complex ID Types

For range-efficient queries, you can define complex IDs using `som.ArrayID` or `som.ObjectID`:

```go
// Array-based ID: stored as [city, date]
type WeatherKey struct {
    som.ArrayID
    City string
    Date time.Time
}

type Weather struct {
    som.Node[WeatherKey]
    Temperature float64
}

// Object-based ID: stored as {name: "...", age: ...}
type PersonKey struct {
    som.ObjectID
    Name string
    Age  int
}

type PersonObj struct {
    som.Node[PersonKey]
    Email string
}
```

Complex IDs enable efficient range queries:

```go
// Query weather for Berlin in a date range
results, _ := client.WeatherRepo().Query().
    Range(
        som.From(WeatherKey{City: "Berlin", Date: startDate}),
        som.To(WeatherKey{City: "Berlin", Date: endDate}),
    ).
    All(ctx)
```

## Timestamps

Add automatic time tracking by embedding `som.Timestamps`:

```go
type User struct {
    som.Node[som.ULID]
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
    som.Node[som.ULID]
    som.OptimisticLock  // Adds version tracking

    Title   string
    Content string
}
```

See [Optimistic Locking](05_optimistic_locking.md) for detailed documentation.

## Soft Delete

Enable non-destructive deletion by embedding `som.SoftDelete`:

```go
type User struct {
    som.Node[som.ULID]
    som.SoftDelete  // Adds soft delete functionality

    Name  string
    Email string
}
```

See [Soft Delete](06_soft_delete.md) for detailed documentation.

## Field Types

Nodes support all [Data Types](../data_types/README.md):

```go
type AllFields struct {
    som.Node[som.ULID]

    // Primitives
    Name     string
    Age      int
    Score    float64
    Active   bool

    // Time
    Birthday  time.Time
    Duration  time.Duration
    Month     time.Month
    Weekday   time.Weekday

    // Special
    ProfileID uuid.UUID
    Website   url.URL
    Contact   som.Email
    Version   som.SemVer

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

## Field Name Override

By default, Go field names are converted to snake_case for the database. Use the `som` tag to override:

```go
type User struct {
    som.Node[som.ULID]

    FullName string `som:"name"`           // Stored as "name" in DB
    EMail    string `som:"email_address"`  // Stored as "email_address"
}
```

## Full-Text Search Index

Mark string fields for full-text search indexing:

```go
type Article struct {
    som.Node[som.ULID]

    Title   string `som:"fulltext=english_search"`
    Content string `som:"fulltext=english_search"`
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
    // Create with auto-generated ID
    Create(ctx context.Context, user *model.User) error

    // Create with specific ID
    CreateWithID(ctx context.Context, id string, user *model.User) error

    // Bulk insert multiple records
    Insert(ctx context.Context, users []*model.User) error

    // Read by ID - returns (model, exists, error)
    Read(ctx context.Context, id string) (*model.User, bool, error)

    // Update existing record (must have ID)
    Update(ctx context.Context, user *model.User) error

    // Delete record
    Delete(ctx context.Context, user *model.User) error

    // Refresh from database
    Refresh(ctx context.Context, user *model.User) error

    // Rebuild all indexes
    // Index access (e.g. per-index Rebuild)
    Index() *index.User

    // Query builder
    Query() query.Builder[model.User]

    // Lifecycle hooks
    OnBeforeCreate(fn func(ctx context.Context, node *model.User) error) func()
    OnAfterCreate(fn func(ctx context.Context, node *model.User) error) func()
    OnBeforeUpdate(fn func(ctx context.Context, node *model.User) error) func()
    OnAfterUpdate(fn func(ctx context.Context, node *model.User) error) func()
    OnBeforeDelete(fn func(ctx context.Context, node *model.User) error) func()
    OnAfterDelete(fn func(ctx context.Context, node *model.User) error) func()
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
    "yourproject/gen/som"
    "github.com/google/uuid"
)

type User struct {
    som.Node[som.ULID]
    som.Timestamps

    // Required fields
    Username string
    Email    som.Email

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
