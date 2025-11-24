# Custom Types

SOM provides custom types for common patterns that aren't covered by Go's standard library.

## som.Enum

Enums provide type-safe enumerated values. See [Enums](../03_models/04_enums.md) for detailed documentation.

```go
type Status string

const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
)

func (s Status) Enum() {}
```

## Planned Custom Types

The following custom types are planned for future releases:

### som.Email

Type-safe email addresses with validation:

```go
type User struct {
    som.Node
    Email som.Email
}
```

### som.Password

Secure password handling with automatic hashing:

```go
type User struct {
    som.Node
    Password som.Password
}
```

### som.Slug

URL-friendly slugs with automatic generation:

```go
type Article struct {
    som.Node
    Title string
    Slug  som.Slug
}
```

### Geometry Types

For spatial data:

- `som.GeometryPoint`
- `som.GeometryLine`
- `som.GeometryPolygon`
- `som.GeometryMultiPoint`
- `som.GeometryMultiLine`
- `som.GeometryMultiPolygon`
- `som.GeometryCollection`

## Slices

All custom types support slice usage:

```go
type User struct {
    som.Node
    Roles []Role  // Slice of enum
}
```
