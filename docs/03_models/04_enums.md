# Enums

SOM provides support for enumerated types through the `som.Enum` interface.

## Defining an Enum

Create an enum by defining a type and implementing the `som.Enum` interface:

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

## Using Enums in Models

Use your enum type as a field in any Node:

```go
type User struct {
    som.Node

    Name   string
    Status Status
}
```

## Type-Safe Queries

Enum fields get type-safe filter operations:

```go
// Filter by enum value
users, err := client.UserRepo().Query().
    Filter(where.User.Status.Equal(model.StatusActive)).
    All(ctx)
```

## Enum Slices

You can also use slices of enums:

```go
type User struct {
    som.Node

    Name  string
    Roles []Role  // Multiple enum values
}
```

## Benefits

- **Type safety**: Only valid enum values can be used
- **Refactoring support**: Rename values with IDE support
- **Documentation**: Enum values are self-documenting in code
