# Special Types

SOM supports several special types from Go's standard library and popular packages.

## UUID

Use `github.com/google/uuid.UUID` for unique identifiers:

```go
import "github.com/google/uuid"

type Document struct {
    som.Node[som.ULID]

    ExternalID uuid.UUID
    TrackingID uuid.UUID
}
```

### CBOR Encoding

UUIDs are encoded using:

- **CBOR Tag 37** for UUID values
- Binary representation of the UUID bytes

This ensures efficient storage and proper round-tripping with SurrealDB.

### Creating UUIDs

```go
doc := &model.Document{
    ExternalID: uuid.New(),
    TrackingID: uuid.New(),
}
```

### Querying by UUID

```go
doc, exists, err := client.DocumentRepo().Query().
    Where(filter.Document.ExternalID.Equal(targetUUID)).
    First(ctx)
```

### UUID Filter Operations

| Operation | Description |
|-----------|-------------|
| `Equal(uuid)` | Equals value |
| `NotEqual(uuid)` | Not equals value |
| `In(uuids...)` | In list |
| `NotIn(uuids...)` | Not in list |

```go
// Find specific document
filter.Document.ExternalID.Equal(targetUUID)

// Find multiple documents
filter.Document.TrackingID.In(uuid1, uuid2, uuid3)
```

## URL

Use `net/url.URL` for web addresses:

```go
import "net/url"

type Bookmark struct {
    som.Node[som.ULID]

    Title string
    Link  url.URL
}
```

### Working with URLs

```go
link, _ := url.Parse("https://example.com/page")
bookmark := &model.Bookmark{
    Title: "Example",
    Link:  *link,
}
```

### URL Filter Operations

URLs support string-like filter operations:

```go
// Find bookmarks with specific host
filter.Bookmark.Link.Host.Equal("example.com")

// Find HTTPS links
filter.Bookmark.Link.Scheme.Equal("https")
```

## Optional Special Types

Use pointers for optional values:

```go
type User struct {
    som.Node[som.ULID]

    ProfileID *uuid.UUID  // Optional UUID
    Website   *url.URL    // Optional URL
}
```

Query optional fields:

```go
// Find users with a profile
filter.User.ProfileID.IsNotNil()

// Find users without a website
filter.User.Website.IsNil()
```

## UUID (gofrs)

SOM also supports `github.com/gofrs/uuid.UUID` as an alternative UUID implementation:

```go
import "github.com/gofrs/uuid"

type Resource struct {
    som.Node[som.ULID]

    ExternalID uuid.UUID
}
```

Both `google/uuid` and `gofrs/uuid` are encoded identically using CBOR Tag 37.

## Built-in Special Types

### som.Email

Type-safe email addresses:

```go
type User struct {
    som.Node[som.ULID]
    Email som.Email
}
```

Filter operations include `Equal`, `In`, `User()` (extract user part), and `Host()` (extract host part).

### som.Password

Secure password handling with automatic hashing. Supports multiple algorithms:

```go
type User struct {
    som.Node[som.ULID]
    Password som.Password[som.Bcrypt]   // Bcrypt algorithm
}

type Admin struct {
    som.Node[som.ULID]
    Password som.Password[som.Argon2]   // Argon2 (most secure)
}
```

**Supported algorithms:** `som.Bcrypt`, `som.Argon2`, `som.Pbkdf2`, `som.Scrypt`

**Key features:**
- Passwords are automatically hashed when stored
- `PERMISSIONS FOR SELECT NONE` - never returned in queries
- Only re-hashes when the value changes

See [Password Type Reference](../types/password.md) for complete documentation.

### som.SemVer

Semantic version strings with comparison and component extraction:

```go
type Release struct {
    som.Node[som.ULID]
    Version som.SemVer
}
```

Filter operations include `Equal`, `Compare`, `Major()`, `Minor()`, and `Patch()`.

## Type Reference

| Type | Package | CBOR Tag | Description |
|------|---------|----------|-------------|
| `time.Time` | `time` | 12 | Datetime with nanosecond precision |
| `time.Duration` | `time` | 14 | Duration with nanosecond precision |
| `time.Month` | `time` | - | Month of the year |
| `time.Weekday` | `time` | - | Day of the week |
| `uuid.UUID` | `github.com/google/uuid` | 37 | Universally unique identifier |
| `uuid.UUID` | `github.com/gofrs/uuid` | 37 | Universally unique identifier |
| `url.URL` | `net/url` | - | Web address |
| `som.Email` | generated | - | Email address string |
| `som.Password[A]` | generated | - | Auto-hashed password |
| `som.SemVer` | generated | - | Semantic version string |
