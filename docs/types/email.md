# Email Type

The email type provides validated email address storage with user and host extraction.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `som.Email` / `*som.Email` |
| Database Schema | `string` / `option<string>` (with validation) |
| CBOR Encoding | Direct (as string) |
| Sortable | Yes |

## Definition

```go
type User struct {
    som.Node

    PrimaryEmail   som.Email   // Required
    SecondaryEmail *som.Email  // Optional
}
```

## Schema

Generated SurrealDB schema with email validation:

```surql
DEFINE FIELD primary_email ON user TYPE string
    ASSERT string::is::email($value);
DEFINE FIELD secondary_email ON user TYPE option<string>
    ASSERT $value == NONE OR $value == NULL OR string::is::email($value);
```

## Creating Email Values

```go
user := &model.User{
    PrimaryEmail: som.Email("alice@example.com"),
}

// Optional email
secondary := som.Email("alice.work@company.com")
user.SecondaryEmail = &secondary
```

## Filter Operations

### Equality Operations

```go
// Exact match
where.User.PrimaryEmail.Equal(som.Email("alice@example.com"))

// Not equal
where.User.PrimaryEmail.NotEqual(som.Email("test@example.com"))
```

### Set Membership

```go
// Value in set
where.User.PrimaryEmail.In(
    som.Email("alice@example.com"),
    som.Email("bob@example.com"),
)

// Value not in set
where.User.PrimaryEmail.NotIn(blockedEmails...)
```

### Comparison Operations

```go
// Lexicographic comparison
where.User.PrimaryEmail.LessThan(som.Email("m@example.com"))
where.User.PrimaryEmail.GreaterThan(som.Email("a@example.com"))
```

### Component Extraction

Extract email parts:

```go
// User part (before @)
where.User.PrimaryEmail.User().Equal("alice")
where.User.PrimaryEmail.User().StartsWith("admin")

// Host part (after @)
where.User.PrimaryEmail.Host().Equal("example.com")
where.User.PrimaryEmail.Host().EndsWith(".edu")
```

### String Operations on Components

After extracting email components, you can use string operations:

```go
// String operations on user part
where.User.PrimaryEmail.User().Lowercase().Equal("alice")
where.User.PrimaryEmail.User().StartsWith("support")

// String operations on host part
where.User.PrimaryEmail.Host().Lowercase().Equal("example.com")
where.User.PrimaryEmail.Host().EndsWith(".edu")
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
where.User.SecondaryEmail.IsNil()

// Check if not nil
where.User.SecondaryEmail.IsNotNil()
```

### Zero Value Check

```go
// Is empty email
where.User.PrimaryEmail.Zero(true)

// Is not empty email
where.User.PrimaryEmail.Zero(false)
```

## Sorting

```go
// Ascending (alphabetic)
query.Order(by.User.PrimaryEmail.Asc())

// Descending
query.Order(by.User.PrimaryEmail.Desc())

// Sort by domain, then user
query.Order(
    by.User.PrimaryEmail.Host().Asc(),
    by.User.PrimaryEmail.User().Asc(),
)
```

## Method Chaining

Email filters support component extraction, which returns String filters for further chaining:

```go
// Find company emails
where.User.PrimaryEmail.Host().Equal("company.com")

// Find admin users
where.User.PrimaryEmail.User().StartsWith("admin")

// Case-insensitive domain match (Host() returns String filter)
where.User.PrimaryEmail.Host().Lowercase().Equal("company.com")

// Complex user filtering (User() returns String filter)
where.User.PrimaryEmail.User().Lowercase().Contains("support")
```

## Common Patterns

### Filter by Domain

```go
// All company emails
companyUsers, _ := client.UserRepo().Query().
    Filter(where.User.PrimaryEmail.Host().Equal("company.com")).
    All(ctx)
```

### Find Gmail Users

```go
gmailUsers, _ := client.UserRepo().Query().
    Filter(where.User.PrimaryEmail.EndsWith("@gmail.com")).
    All(ctx)
```

### Find Admin Emails

```go
admins, _ := client.UserRepo().Query().
    Filter(where.User.PrimaryEmail.User().StartsWith("admin")).
    All(ctx)
```

### Users with Secondary Email

```go
withSecondary, _ := client.UserRepo().Query().
    Filter(where.User.SecondaryEmail.IsNotNil()).
    All(ctx)
```

### Educational Institutions

```go
eduUsers, _ := client.UserRepo().Query().
    Filter(where.User.PrimaryEmail.Host().EndsWith(".edu")).
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
    "yourproject/model"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create user with email
    user := &model.User{
        PrimaryEmail: som.Email("alice@company.com"),
    }
    client.UserRepo().Create(ctx, user)

    // Find by exact email
    found, exists, _ := client.UserRepo().Query().
        Filter(where.User.PrimaryEmail.Equal(som.Email("alice@company.com"))).
        First(ctx)

    // Find all company employees
    employees, _ := client.UserRepo().Query().
        Filter(where.User.PrimaryEmail.Host().Equal("company.com")).
        Order(by.User.PrimaryEmail.Asc()).
        All(ctx)

    // Find users by email pattern
    supportTeam, _ := client.UserRepo().Query().
        Filter(where.User.PrimaryEmail.User().StartsWith("support")).
        All(ctx)

    // Find Gmail users
    gmailUsers, _ := client.UserRepo().Query().
        Filter(where.User.PrimaryEmail.EndsWith("@gmail.com")).
        All(ctx)

    // Users with backup email configured
    withBackup, _ := client.UserRepo().Query().
        Filter(where.User.SecondaryEmail.IsNotNil()).
        All(ctx)

    // Case-insensitive search
    caseInsensitive, _ := client.UserRepo().Query().
        Filter(
            where.User.PrimaryEmail.Lowercase().Equal("alice@company.com"),
        ).
        All(ctx)
}
```

## Filter Reference Table

### Base Operations

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Equal(val)` | Exact match | Bool filter |
| `NotEqual(val)` | Not equal | Bool filter |
| `In(vals...)` | Value in set | Bool filter |
| `NotIn(vals...)` | Value not in set | Bool filter |
| `Zero(bool)` | Check empty | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `IsNil()` | Is null (ptr only) | Bool filter |
| `IsNotNil()` | Not null (ptr only) | Bool filter |

### Component Extraction

| Operation | Description | Returns |
|-----------|-------------|---------|
| `User()` | Extract user part (before @) | String filter |
| `Host()` | Extract host part (after @) | String filter |

After extraction, all [String filter operations](string.md) are available on the component.
