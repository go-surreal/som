# Bool Type

The boolean type represents true/false values with simple, focused filter operations.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `bool` / `*bool` |
| Database Schema | `bool` / `option<bool>` |
| CBOR Encoding | Direct |
| Sortable | Yes |

## Definition

```go
type User struct {
    som.Node

    IsActive    bool   // Required
    IsAdmin     bool   // Required
    IsVerified  *bool  // Optional (nullable)
}
```

## Schema

Generated SurrealDB schema:

```surql
DEFINE FIELD is_active ON user TYPE bool;
DEFINE FIELD is_admin ON user TYPE bool;
DEFINE FIELD is_verified ON user TYPE option<bool>;
```

## Filter Operations

### Value Check

```go
// Check specific value
where.User.IsActive.Is(true)
where.User.IsAdmin.Is(false)
```

### Convenience Methods

```go
// Check if true
where.User.IsActive.True()

// Check if false
where.User.IsAdmin.False()
```

### Logical Inversion

```go
// Invert the boolean
where.User.IsActive.Invert().True()  // Same as Is(false)
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
where.User.IsVerified.IsNil()

// Check if not nil
where.User.IsVerified.IsNotNil()
```

## Sorting

Boolean fields sort with `false` before `true` in ascending order:

```go
// Ascending (false, then true)
query.Order(by.User.IsActive.Asc())

// Descending (true, then false)
query.Order(by.User.IsActive.Desc())

// Common pattern: active users first
query.Order(
    by.User.IsActive.Desc(),  // Active first
    by.User.Name.Asc(),       // Then alphabetical
)
```

## Method Chaining

Boolean filters can be inverted:

```go
// Find inactive users
where.User.IsActive.Invert().True()

// Equivalent to
where.User.IsActive.False()
where.User.IsActive.Is(false)
```

## Common Patterns

### Combining Boolean Filters

```go
// Active admin users
users, _ := client.UserRepo().Query().
    Filter(
        where.User.IsActive.True(),
        where.User.IsAdmin.True(),
    ).
    All(ctx)
```

### OR Logic with Booleans

```go
// Active OR admin users
users, _ := client.UserRepo().Query().
    Filter(
        where.Any(
            where.User.IsActive.True(),
            where.User.IsAdmin.True(),
        ),
    ).
    All(ctx)
```

### Optional Boolean Handling

```go
// Users with verified status set
verifiedSet, _ := client.UserRepo().Query().
    Filter(where.User.IsVerified.IsNotNil()).
    All(ctx)

// Users explicitly verified
verified, _ := client.UserRepo().Query().
    Filter(
        where.User.IsVerified.IsNotNil(),
        where.User.IsVerified.True(),
    ).
    All(ctx)

// Users not verified (either false or null)
notVerified, _ := client.UserRepo().Query().
    Filter(
        where.Any(
            where.User.IsVerified.IsNil(),
            where.User.IsVerified.False(),
        ),
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
    "yourproject/gen/som/where"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Find active users
    activeUsers, _ := client.UserRepo().Query().
        Filter(where.User.IsActive.True()).
        All(ctx)

    // Find inactive non-admin users
    inactiveNonAdmins, _ := client.UserRepo().Query().
        Filter(
            where.User.IsActive.False(),
            where.User.IsAdmin.False(),
        ).
        All(ctx)

    // Count admins
    adminCount, _ := client.UserRepo().Query().
        Filter(where.User.IsAdmin.True()).
        Count(ctx)

    // Check if any unverified users exist
    hasUnverified, _ := client.UserRepo().Query().
        Filter(
            where.Any(
                where.User.IsVerified.IsNil(),
                where.User.IsVerified.False(),
            ),
        ).
        Exists(ctx)

    // Get users sorted by status
    sorted, _ := client.UserRepo().Query().
        Order(
            by.User.IsAdmin.Desc(),   // Admins first
            by.User.IsActive.Desc(),  // Then active
            by.User.Name.Asc(),       // Then alphabetical
        ).
        All(ctx)
}
```

## Filter Reference Table

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Is(val)` | Check boolean value | Bool filter |
| `True()` | Check if true | Bool filter |
| `False()` | Check if false | Bool filter |
| `Invert()` | Logical NOT | Bool filter |
| `IsNil()` | Is null (ptr only) | Bool filter |
| `IsNotNil()` | Not null (ptr only) | Bool filter |
