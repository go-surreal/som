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

### Creating UUIDs

```go
doc := &model.Document{
    ExternalID: uuid.New(),
    TrackingID: uuid.New(),
}
```

### Querying by UUID

```go
doc, err := client.DocumentRepo().Query().
    Filter(where.Document.ExternalID.Equal(targetUUID)).
    First(ctx)
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

## Optional Special Types

Use pointers for optional values:

```go
type User struct {
    som.Node

    ProfileID *uuid.UUID  // Optional UUID
    Website   *url.URL    // Optional URL
}
```

## Future Types

The following types are planned for future releases:

- `net.IP` - IP addresses
- `regexp.Regexp` - Regular expressions
- `big.Int` / `big.Float` - Arbitrary precision numbers
