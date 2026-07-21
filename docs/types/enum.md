# Enum Type

The enum type provides type-safe enumerated values with a constrained set of allowed string values.

## Overview

| Property | Value |
|----------|-------|
| Go Type | Custom type implementing enum pattern |
| Database Schema | String union (e.g., `"active" \| "inactive"`) |
| CBOR Encoding | Direct (as string) |
| Sortable | Yes |

## Definition

Define an enum by creating a string type with constants:

```go
package model

type Status string

const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
    StatusPending  Status = "pending"
)

// Enum marker method (required)
func (s Status) Enum() {}

type User struct {
    som.Node

    Name   string
    Status Status   // Required enum
    Role   *Role    // Optional enum
}
```

### Multiple Enums

```go
type Role string

const (
    RoleAdmin     Role = "admin"
    RoleModerator Role = "moderator"
    RoleUser      Role = "user"
    RoleGuest     Role = "guest"
)

func (r Role) Enum() {}
```

## Schema

Generated SurrealDB schema with union type:

```surql
-- Values are sorted alphabetically in the schema
DEFINE FIELD status ON user TYPE "active" | "inactive" | "pending";
DEFINE FIELD role ON user TYPE option<"admin" | "guest" | "moderator" | "user">;
```

## Filter Operations

### Equality Operations

```go
// Exact match
filter.User.Status.Equal(model.StatusActive)

// Not equal
filter.User.Status.NotEqual(model.StatusPending)
```

### Set Membership

```go
// Value in set
filter.User.Role.In(model.RoleAdmin, model.RoleModerator)

// Value not in set
filter.User.Role.NotIn(model.RoleGuest)
```

### Comparison Operations

Enums can be compared lexicographically:

```go
// Lexicographic comparison (alphabetic)
filter.User.Status.LessThan(model.StatusPending)
filter.User.Status.GreaterThan(model.StatusActive)
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
filter.User.Role.IsNil()

// Check if not nil
filter.User.Role.IsNotNil()
```

### Zero Value Check

```go
// Is empty string (zero value)
filter.User.Status.Zero(true)

// Is not empty
filter.User.Status.Zero(false)
```

## Sorting

```go
// Ascending (alphabetic by value)
query.Order(by.User.Status.Asc())

// Descending
query.Order(by.User.Status.Desc())
```

## Common Patterns

### Filter by Status

```go
// Active users only
activeUsers, _ := client.UserRepo().Query().
    Where(filter.User.Status.Equal(model.StatusActive)).
    All(ctx)
```

### Multiple Allowed Values

```go
// Admins or moderators
privilegedUsers, _ := client.UserRepo().Query().
    Where(
        filter.User.Role.In(model.RoleAdmin, model.RoleModerator),
    ).
    All(ctx)
```

### Exclude Values

```go
// Everyone except guests
nonGuests, _ := client.UserRepo().Query().
    Where(filter.User.Role.NotEqual(model.RoleGuest)).
    All(ctx)
```

### Optional Enum with Default

```go
// Users with role set, or default to user
usersWithRole, _ := client.UserRepo().Query().
    Where(filter.User.Role.IsNotNil()).
    All(ctx)
```

### Count by Status

```go
activeCount, _ := client.UserRepo().Query().
    Where(filter.User.Status.Equal(model.StatusActive)).
    Count(ctx)

pendingCount, _ := client.UserRepo().Query().
    Where(filter.User.Status.Equal(model.StatusPending)).
    Count(ctx)
```

## Enum Slices

Enums can be used in slices:

```go
type User struct {
    som.Node

    Name  string
    Roles []Role  // Multiple roles
}
```

Query enum slices:

```go
// Users with admin role
admins, _ := client.UserRepo().Query().
    Where(filter.User.Roles.Contains(model.RoleAdmin)).
    All(ctx)

// Users with any privileged role
privileged, _ := client.UserRepo().Query().
    Where(
        filter.User.Roles.ContainsAny(model.RoleAdmin, model.RoleModerator),
    ).
    All(ctx)
```

## Complete Example

```go
package main

import (
    "context"
    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/filter"
    "yourproject/model"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create user with enum
    user := &model.User{
        Name:   "Alice",
        Status: model.StatusActive,
    }
    role := model.RoleAdmin
    user.Role = &role
    client.UserRepo().Create(ctx, user)

    // Find active users
    active, _ := client.UserRepo().Query().
        Where(filter.User.Status.Equal(model.StatusActive)).
        All(ctx)

    // Find admins and moderators
    privileged, _ := client.UserRepo().Query().
        Where(
            filter.User.Role.In(model.RoleAdmin, model.RoleModerator),
        ).
        All(ctx)

    // Find users without role assigned
    noRole, _ := client.UserRepo().Query().
        Where(filter.User.Role.IsNil()).
        All(ctx)

    // Exclude pending users
    notPending, _ := client.UserRepo().Query().
        Where(filter.User.Status.NotEqual(model.StatusPending)).
        All(ctx)

    // Count by status
    activeCount, _ := client.UserRepo().Query().
        Where(filter.User.Status.Equal(model.StatusActive)).
        Count(ctx)

    inactiveCount, _ := client.UserRepo().Query().
        Where(filter.User.Status.Equal(model.StatusInactive)).
        Count(ctx)

    // Sort by status, then name
    sorted, _ := client.UserRepo().Query().
        Order(
            by.User.Status.Asc(),
            by.User.Name.Asc(),
        ).
        All(ctx)
}
```

## Best Practices

### Use Constants

Always define constants for enum values:

```go
// Good
const StatusActive Status = "active"
user.Status = StatusActive

// Avoid
user.Status = Status("active")  // Works but bypasses type safety
```

### Validation

The database schema ensures only valid values are stored. Invalid values will cause database errors:

```go
// This will fail at database level
user.Status = Status("invalid_status")
client.UserRepo().Create(ctx, user)  // Error: invalid enum value
```

### Empty Values

Be careful with zero values:

```go
var user model.User
// user.Status is "" (empty string) - may not be valid

// Check for empty
if user.Status == "" {
    user.Status = model.StatusPending  // Set default
}
```

## Filter Reference Table

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Equal(val)` | Exact match | Bool filter |
| `NotEqual(val)` | Not equal | Bool filter |
| `In(vals...)` | Value in set | Bool filter |
| `NotIn(vals...)` | Value not in set | Bool filter |
| `LessThan(val)` | Lexicographic < | Bool filter |
| `LessThanEqual(val)` | Lexicographic <= | Bool filter |
| `GreaterThan(val)` | Lexicographic > | Bool filter |
| `GreaterThanEqual(val)` | Lexicographic >= | Bool filter |
| `Zero(bool)` | Check empty | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `IsNil()` | Is null (ptr) | Bool filter |
| `IsNotNil()` | Not null (ptr) | Bool filter |
