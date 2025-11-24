# Enums

SOM provides support for enumerated types through the `som.Enum` interface. Enums are stored as strings in SurrealDB but provide type safety in Go.

## Defining an Enum

Create an enum by defining a string type and implementing the `som.Enum` interface:

```go
package model

type Status string

const (
    StatusPending  Status = "pending"
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
)

// Implement som.Enum interface
func (s Status) Enum() {}
```

The `Enum()` method is a marker that tells SOM to treat this type as an enumeration.

## Using Enums in Models

Use your enum type as a field in any Node or Edge:

```go
type User struct {
    som.Node

    Name   string
    Status Status
}

type Order struct {
    som.Node

    OrderNumber string
    Status      OrderStatus
    Priority    Priority
}
```

## Enum Values

Define constants for all valid values:

```go
type Priority string

const (
    PriorityLow    Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh   Priority = "high"
    PriorityCritical Priority = "critical"
)

func (p Priority) Enum() {}
```

## Type-Safe Queries

Enum fields get type-safe filter operations:

```go
// Filter by enum value
users, err := client.UserRepo().Query().
    Filter(where.User.Status.Equal(model.StatusActive)).
    All(ctx)

// Filter with In (multiple values)
users, err := client.UserRepo().Query().
    Filter(where.User.Status.In(
        model.StatusActive,
        model.StatusPending,
    )).
    All(ctx)

// Filter with NotEqual
users, err := client.UserRepo().Query().
    Filter(where.User.Status.NotEqual(model.StatusInactive)).
    All(ctx)
```

## Available Filter Operations

Enum fields support these operations:

| Operation | Description |
|-----------|-------------|
| `Equal(val)` | Equals value |
| `NotEqual(val)` | Not equals value |
| `In(vals...)` | Value is one of |
| `NotIn(vals...)` | Value is not one of |

## Enum Slices

Use slices of enums for multiple values:

```go
type Role string

const (
    RoleAdmin  Role = "admin"
    RoleEditor Role = "editor"
    RoleViewer Role = "viewer"
)

func (r Role) Enum() {}

type User struct {
    som.Node

    Name  string
    Roles []Role  // User can have multiple roles
}
```

Query with slice operations:

```go
// Find users with admin role
users, err := client.UserRepo().Query().
    Filter(where.User.Roles.Contains(model.RoleAdmin)).
    All(ctx)

// Find users with any of these roles
users, err := client.UserRepo().Query().
    Filter(where.User.Roles.ContainsAny(
        model.RoleAdmin,
        model.RoleEditor,
    )).
    All(ctx)
```

## Optional Enums

Use pointers for optional enum fields:

```go
type User struct {
    som.Node

    Name     string
    Status   Status   // Required
    Priority *Priority // Optional, can be nil
}
```

Filter optional enums:

```go
// Find users with priority set
users, err := client.UserRepo().Query().
    Filter(where.User.Priority.IsNotNil()).
    All(ctx)

// Find high priority users
users, err := client.UserRepo().Query().
    Filter(where.User.Priority.Equal(model.PriorityHigh)).
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
type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusConfirmed OrderStatus = "confirmed"
    // ...
)

func (s OrderStatus) Enum() {}
```

### Document Valid Values

```go
// PaymentMethod represents accepted payment types.
// Valid values: credit_card, debit_card, bank_transfer, crypto
type PaymentMethod string

const (
    PaymentCreditCard   PaymentMethod = "credit_card"
    PaymentDebitCard    PaymentMethod = "debit_card"
    PaymentBankTransfer PaymentMethod = "bank_transfer"
    PaymentCrypto       PaymentMethod = "crypto"
)

func (p PaymentMethod) Enum() {}
```

## Benefits

- **Type safety**: Only defined enum values can be used
- **Compile-time checks**: Invalid values caught by compiler
- **IDE support**: Autocompletion shows valid options
- **Refactoring**: Rename values safely with IDE tools
- **Self-documenting**: Code clearly shows valid options
