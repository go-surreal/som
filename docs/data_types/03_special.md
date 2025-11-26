# Special Types

SOM supports several special types from Go's standard library and popular packages.

## UUID

Use `github.com/google/uuid.UUID` for unique identifiers:

```go
import "github.com/google/uuid"

type Document struct {
    som.Node

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
    Filter(where.Document.ExternalID.Equal(targetUUID)).
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
where.Document.ExternalID.Equal(targetUUID)

// Find multiple documents
where.Document.TrackingID.In(uuid1, uuid2, uuid3)
```

## URL

Use `net/url.URL` for web addresses:

```go
import "net/url"

type Bookmark struct {
    som.Node

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
where.Bookmark.Link.Host.Equal("example.com")

// Find HTTPS links
where.Bookmark.Link.Scheme.Equal("https")
```

## Optional Special Types

Use pointers for optional values:

```go
type User struct {
    som.Node

    ProfileID *uuid.UUID  // Optional UUID
    Website   *url.URL    // Optional URL
}
```

Query optional fields:

```go
// Find users with a profile
where.User.ProfileID.IsNotNil()

// Find users without a website
where.User.Website.IsNil()
```

## Record IDs

SOM uses SurrealDB's native Record ID format:

```go
type ID = models.RecordID
```

Record IDs are automatically managed:

```go
user := &model.User{Name: "Alice"}
err := client.UserRepo().Create(ctx, user)
// user.ID() returns *som.ID like "user:01HQ..."
```

### Creating Custom IDs

Use `CreateWithID` for specific IDs:

```go
err := client.UserRepo().CreateWithID(ctx, "alice", user)
// Creates user:alice
```

### ID Utilities

```go
// Create a new RecordID
id := som.NewRecordID("user", "alice")

// Create a pointer to RecordID
idPtr := som.MakeID("user", "alice")

// Reference a table (for server-generated IDs)
table := som.Table("user")
```

## Built-in Special Types

SOM provides marker types for special handling:

### som.Email

Type-safe email addresses:

```go
type User struct {
    som.Node
    Email som.Email
}
```

### som.Password

Secure password handling with automatic hashing. Supports multiple algorithms:

```go
type User struct {
    som.Node
    Password som.Password[som.Bcrypt]   // Bcrypt algorithm
}

type Admin struct {
    som.Node
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

Semantic version strings:

```go
type Release struct {
    som.Node
    Version som.SemVer
}
```

## Type Reference

| Type | Package | CBOR Tag | Description |
|------|---------|----------|-------------|
| `time.Time` | `time` | 12 | Datetime with nanosecond precision |
| `time.Duration` | `time` | 14 | Duration with nanosecond precision |
| `uuid.UUID` | `github.com/google/uuid` | 37 | Universally unique identifier |
| `url.URL` | `net/url` | - | Web address |
| `som.Email` | generated | - | Email address string |
| `som.Password` | generated | - | Password string |
| `som.SemVer` | generated | - | Semantic version string |
