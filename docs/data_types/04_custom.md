# Custom Types

SOM provides custom types for common patterns that aren't covered by Go's standard library.

## som.Enum

Enums provide type-safe enumerated values. See [Enums](../models/04_enums.md) for detailed documentation.

```go
type Status string

const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
)

func (s Status) Enum() {}
```

## som.Email

Type-safe email addresses with filter support:

```go
type User struct {
    som.Node[som.ULID]
    Email som.Email
}
```

Filter operations:

```go
filter.User.Email.Equal("john@example.com")
filter.User.Email.In("a@example.com", "b@example.com")
filter.User.Email.User()    // Extract user part
filter.User.Email.Host()    // Extract host part
```

## som.Password

Secure password handling with automatic hashing. Passwords are hashed using the specified algorithm before being stored in SurrealDB, and are never returned in query results.

```go
type User struct {
    som.Node[som.ULID]
    Password som.Password[som.Bcrypt]
}
```

### Supported Algorithms

| Algorithm | Type | Notes |
|-----------|------|-------|
| `som.Bcrypt` | Bcrypt | Widely used, good default |
| `som.Argon2` | Argon2id | Most secure, memory-hard |
| `som.Pbkdf2` | PBKDF2 | NIST recommended |
| `som.Scrypt` | Scrypt | Memory-hard alternative |

### Generated Schema

```surql
DEFINE FIELD password ON TABLE user TYPE string
    VALUE IF $input AND crypto::bcrypt::compare($before, $input) = false {
        crypto::bcrypt::generate($input)
    } ELSE IF $before {
        $before
    } ELSE {
        crypto::bcrypt::generate($value)
    }
    PERMISSIONS FOR SELECT NONE;
```

## som.SemVer

Semantic version strings with comparison and component extraction:

```go
type Release struct {
    som.Node[som.ULID]
    Version som.SemVer
}
```

Filter operations:

```go
filter.Release.Version.Equal("1.2.3")
filter.Release.Version.Major().GreaterThan(1)
filter.Release.Version.Minor().Equal(2)
filter.Release.Version.Patch().LessThan(5)
```

## Geometry Types

SOM supports geometry types from three popular Go libraries:

### paulmach/orb

```go
import "github.com/paulmach/orb"

type Location struct {
    som.Node[som.ULID]
    Position orb.Point
    Area     orb.Polygon
}
```

### peterstace/simplefeatures

```go
import "github.com/peterstace/simplefeatures/geom"

type Location struct {
    som.Node[som.ULID]
    Position geom.Point
    Area     geom.Polygon
}
```

### twpayne/go-geom

```go
import "github.com/twpayne/go-geom"

type Location struct {
    som.Node[som.ULID]
    Position *gogeom.Point
    Area     *gogeom.Polygon
}
```

### Supported Geometry Types

All three libraries support:

| Type | Description |
|------|-------------|
| Point | Single coordinate |
| LineString | Series of connected points |
| Polygon | Closed shape |
| MultiPoint | Collection of points |
| MultiLineString | Collection of line strings |
| MultiPolygon | Collection of polygons |
| GeometryCollection | Mixed geometry collection |

## Slices

All custom types support slice usage:

```go
type User struct {
    som.Node[som.ULID]
    Roles []Role  // Slice of enum
    Tags  []string
}
```
