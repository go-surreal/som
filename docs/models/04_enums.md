# Enums

SOM supports enumerated types via `som.Enum`. Enums are stored as strings in SurrealDB but keep
type safety in Go.

## Defining an Enum

Declare a type whose underlying type is `som.Enum`, inside your model package:

```go
package model

import "yourproject/gen/som"

type Status som.Enum

const (
    StatusPending  Status = "pending"
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
)
```

Any type declared as `som.Enum` in the model package is picked up by the generator as an
enumeration. No marker method is required; the constants define the allowed values.

## Using Enums in Models

Use your enum type as a field in any Node or Edge:

```go
type User struct {
    som.Node[som.ULID]

    Name   string
    Status Status
}

type Order struct {
    som.Node[som.ULID]

    OrderNumber string
    Status      OrderStatus
    Priority    Priority
}
```

## Enum Values

Define constants for all valid values:

```go
type Priority som.Enum

const (
    PriorityLow    Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh   Priority = "high"
    PriorityCritical Priority = "critical"
)

```

## Type-Safe Queries

Enum fields get type-safe filter operations:

```go
// Filter by enum value
users, err := client.UserRepo().Query().
    Where(filter.User.Status.Equal(model.StatusActive)).
    All(ctx)

// Filter with In (multiple values)
users, err := client.UserRepo().Query().
    Where(filter.User.Status.In([]model.Status{
        model.StatusActive,
        model.StatusPending,
    })).
    All(ctx)

// Filter with NotEqual
users, err := client.UserRepo().Query().
    Where(filter.User.Status.NotEqual(model.StatusInactive)).
    All(ctx)
```

## Available Filter Operations

Enum fields support these operations:

| Operation | Description |
|-----------|-------------|
| `Equal(val)` | Equals value |
| `NotEqual(val)` | Not equals value |
| `In(vals []T)` | Value is one of |
| `NotIn(vals []T)` | Value is not one of |

## Enum Slices

Use slices of enums for multiple values:

```go
type Role som.Enum

const (
    RoleAdmin  Role = "admin"
    RoleEditor Role = "editor"
    RoleViewer Role = "viewer"
)


type User struct {
    som.Node[som.ULID]

    Name  string
    Roles []Role  // User can have multiple roles
}
```

Query with slice operations:

```go
// Find users with admin role
users, err := client.UserRepo().Query().
    Where(filter.User.Roles.Contains(model.RoleAdmin)).
    All(ctx)

// Find users with any of these roles
users, err := client.UserRepo().Query().
    Where(filter.User.Roles.ContainsAny([]model.Role{
        model.RoleAdmin,
        model.RoleEditor,
    })).
    All(ctx)
```

## Optional Enums

Use pointers for optional enum fields:

```go
type User struct {
    som.Node[som.ULID]

    Name     string
    Status   Status   // Required
    Priority *Priority // Optional, can be nil
}
```

Filter optional enums:

```go
// Find users with priority set
users, err := client.UserRepo().Query().
    Where(filter.User.Priority.Nil(false)).
    All(ctx)

// Find high priority users
users, err := client.UserRepo().Query().
    Where(filter.User.Priority.Equal(model.PriorityHigh)).
    All(ctx)
```

## Sorting by Enum

Enums can be used in ORDER BY (sorts alphabetically by string value):

```go
users, err := client.UserRepo().Query().
    Order(by.User.Status.Asc()).
    All(ctx)
```

## Best Practices

### Use Descriptive Constants

```go
// Good
const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusConfirmed OrderStatus = "confirmed"
    OrderStatusShipped   OrderStatus = "shipped"
)

// Avoid generic names
const (
    Status1 Status = "1"  // Unclear
    Status2 Status = "2"
)
```

### Group Related Enums

```go
// order_status.go
type OrderStatus som.Enum

const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusConfirmed OrderStatus = "confirmed"
    // ...
)

```

### Document Valid Values

```go
// PaymentMethod represents accepted payment types.
// Valid values: credit_card, debit_card, bank_transfer, crypto
type PaymentMethod som.Enum

const (
    PaymentCreditCard   PaymentMethod = "credit_card"
    PaymentDebitCard    PaymentMethod = "debit_card"
    PaymentBankTransfer PaymentMethod = "bank_transfer"
    PaymentCrypto       PaymentMethod = "crypto"
)

```

## Benefits

- **Type safety**: Only defined enum values can be used
- **Compile-time checks**: Invalid values caught by compiler
- **IDE support**: Autocompletion shows valid options
- **Refactoring**: Rename values safely with IDE tools
- **Self-documenting**: Code clearly shows valid options
