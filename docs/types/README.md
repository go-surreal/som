# Types Reference

This section provides comprehensive documentation for all field types supported by SOM. Each type includes its Go definition, SurrealDB schema, available filter operations, sorting capabilities, and practical examples.

## Quick Reference

| Type | Go Type | DB Schema | CBOR | Sortable | Filters |
|------|---------|-----------|------|----------|---------|
| [String](string.md) | `string` | `string` | Direct | Yes | 50+ |
| [Numeric](numeric.md) | `int`, `float64`, etc. | `int` / `float` | Direct | Yes | 20+ |
| [Bool](bool.md) | `bool` | `bool` | Direct | Yes | 4 |
| [Time](time.md) | `time.Time` | `datetime` | Tag 12 | Yes | 25+ |
| [Duration](duration.md) | `time.Duration` | `duration` | Tag 14 | Yes | 15+ |
| [UUID](uuid.md) | `uuid.UUID` | `uuid` | Tag 37 | Yes | 8 |
| [URL](url.md) | `url.URL` | `string` | Direct | Yes | 15+ |
| [Email](email.md) | `som.Email` | `string` | Direct | Yes | 10+ |
| [Enum](enum.md) | Custom type | String union | Direct | Yes | 8 |
| [Slice](slice.md) | `[]T` | `array<T>` | Direct | No | 30+ |
| [Node](node.md) | `*OtherNode` | `record<table>` | Direct | No | Via fields |
| [Struct](struct.md) | Embedded struct | `object` | Direct | No | Via fields |

## Type Categories

### Primitive Types

Basic Go types with direct database mapping:

- **[String](string.md)** - Text data with extensive string operations
- **[Numeric](numeric.md)** - Integers and floating-point numbers
- **[Bool](bool.md)** - Boolean true/false values

### Time Types

Temporal data with special CBOR encoding:

- **[Time](time.md)** - Timestamps using `time.Time`
- **[Duration](duration.md)** - Time intervals using `time.Duration`

### Special Types

Types with validation or special handling:

- **[UUID](uuid.md)** - Universally unique identifiers
- **[URL](url.md)** - Web addresses with component parsing
- **[Email](email.md)** - Email addresses with validation
- **[Enum](enum.md)** - Constrained string values

### Complex Types

Composite and reference types:

- **[Slice](slice.md)** - Arrays of any type
- **[Node](node.md)** - References to other records
- **[Struct](struct.md)** - Nested object structures

## Universal Operations

All types support these base filter operations:

```go
// Equality
where.User.Name.Equal("Alice")
where.User.Name.NotEqual("Bob")

// Set membership
where.User.Role.In("admin", "moderator")
where.User.Status.NotIn("banned", "suspended")

// Zero value check
where.User.Name.Zero(true)   // Is empty string
where.User.Age.Zero(false)   // Is not zero

// Truthiness
where.User.Name.Truth()      // Convert to bool filter
```

### Pointer Types (Optional Fields)

All types support pointer variants for optional fields:

```go
type User struct {
    som.Node

    Name     string   // Required
    Nickname *string  // Optional
}
```

Pointer types add nil-checking operations:

```go
where.User.Nickname.IsNil()     // Field is NULL/NONE
where.User.Nickname.IsNotNil()  // Field has a value
```

## CBOR Encoding

SOM uses CBOR (Concise Binary Object Representation) for efficient database communication. Most types use direct CBOR encoding, but some require special handling:

| Type | CBOR Tag | Format |
|------|----------|--------|
| DateTime | 12 | `[unix_seconds, nanoseconds]` |
| Duration | 14 | `[seconds, nanoseconds]` |
| UUID | 37 | Binary (16 bytes) |

## Sorting

Most types support ascending and descending sort:

```go
// Ascending
query.Order(by.User.Name.Asc())

// Descending
query.Order(by.User.CreatedAt.Desc())

// Multiple fields
query.Order(
    by.User.IsActive.Desc(),
    by.User.Name.Asc(),
)
```

**Not sortable**: Slice, Node, Struct (use nested field sorting instead)

## Method Chaining

Many filter operations return new filters, enabling powerful chains:

```go
// String transformations
where.User.Email.Lowercase().Contains("@gmail")

// Numeric math
where.Product.Price.Mul(1.1).LessThan(100)

// Time extraction
where.Event.StartTime.Year().Equal(2024)

// Slice operations
where.Post.Tags.Distinct().Len().GreaterThan(3)
```

## Schema Generation

SOM generates SurrealDB schema definitions for each field:

```surql
DEFINE FIELD name ON user TYPE string;
DEFINE FIELD age ON user TYPE option<int>;
DEFINE FIELD email ON user TYPE string ASSERT string::is::email($value);
DEFINE FIELD role ON user TYPE "admin" | "user" | "guest";
DEFINE FIELD created_at ON user TYPE datetime VALUE $before OR time::now() READONLY;
```
