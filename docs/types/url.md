# URL Type

The URL type handles web addresses using Go's `net/url.URL` with component parsing capabilities.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `url.URL` / `*url.URL` |
| Database Schema | `string` / `option<string>` (with validation) |
| CBOR Encoding | Direct (as string) |
| Sortable | Yes |

## Definition

```go
import "net/url"

type Bookmark struct {
    som.Node

    Title   string
    Link    url.URL   // Required
    Favicon *url.URL  // Optional
}
```

## Schema

Generated SurrealDB schema with URL validation:

```surql
DEFINE FIELD link ON bookmark TYPE string
    ASSERT string::is::url($value);
DEFINE FIELD favicon ON bookmark TYPE option<string>
    ASSERT $value == NONE OR $value == NULL OR string::is::url($value);
```

## Creating URLs

```go
import "net/url"

// Parse from string
link, _ := url.Parse("https://example.com/page?query=value#section")

bookmark := &model.Bookmark{
    Title: "Example",
    Link:  *link,
}

// Optional URL
faviconURL, _ := url.Parse("https://example.com/favicon.ico")
bookmark.Favicon = faviconURL
```

## Filter Operations

### Equality Operations

```go
// Exact match
targetURL, _ := url.Parse("https://example.com")
filter.Bookmark.Link.Equal(*targetURL)

// Not equal
filter.Bookmark.Link.NotEqual(*excludeURL)
```

### Set Membership

```go
// Value in set
filter.Bookmark.Link.In(url1, url2, url3)

// Value not in set
filter.Bookmark.Link.NotIn(blockedURLs...)
```

### Comparison Operations

```go
// Lexicographic comparison (as strings)
filter.Bookmark.Link.LessThan(referenceURL)
filter.Bookmark.Link.GreaterThan(referenceURL)
```

### Component Extraction

Extract URL components as string filters:

```go
// Domain (without port)
filter.Bookmark.Link.Domain().Equal("example.com")

// Host (with port if present)
filter.Bookmark.Link.Host().Equal("example.com:8080")

// Port number
filter.Bookmark.Link.Port().Equal("443")

// Path
filter.Bookmark.Link.Path().StartsWith("/api/")

// Query string
filter.Bookmark.Link.Query().Contains("utm_source")

// Fragment (after #)
filter.Bookmark.Link.Fragment().Equal("section1")
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
filter.Bookmark.Favicon.IsNil()

// Check if not nil
filter.Bookmark.Favicon.IsNotNil()
```

### Zero Value Check

```go
// Is empty URL
filter.Bookmark.Link.Zero(true)

// Is not empty URL
filter.Bookmark.Link.Zero(false)
```

## Sorting

```go
// Ascending (alphabetic by full URL)
query.Order(by.Bookmark.Link.Asc())

// Descending
query.Order(by.Bookmark.Link.Desc())
```

## Method Chaining

URL filters support component extraction and string operations:

```go
// Find bookmarks on specific domain
filter.Bookmark.Link.Domain().Equal("github.com")

// Find API endpoints
filter.Bookmark.Link.Path().StartsWith("/api/v2/")

// Find URLs with specific query parameter
filter.Bookmark.Link.Query().Contains("page=")

// Find HTTPS URLs (via string contains)
filter.Bookmark.Link.Host().StartsWith("https://")
```

## Common Patterns

### Filter by Domain

```go
// All GitHub bookmarks
githubLinks, _ := client.BookmarkRepo().Query().
    Where(filter.Bookmark.Link.Domain().Equal("github.com")).
    All(ctx)
```

### Find API Endpoints

```go
// All API endpoints
apiEndpoints, _ := client.BookmarkRepo().Query().
    Where(filter.Bookmark.Link.Path().StartsWith("/api/")).
    All(ctx)
```

### URLs with Specific Port

```go
// Development servers
devServers, _ := client.BookmarkRepo().Query().
    Where(filter.Bookmark.Link.Port().Equal("3000")).
    All(ctx)
```

### Bookmarks with Favicons

```go
// Bookmarks with favicon configured
withFavicon, _ := client.BookmarkRepo().Query().
    Where(filter.Bookmark.Favicon.IsNotNil()).
    All(ctx)
```

## Complete Example

```go
package main

import (
    "context"
    "net/url"
    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/filter"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create bookmark
    link, _ := url.Parse("https://github.com/go-surreal/som")
    bookmark := &model.Bookmark{
        Title: "SOM Repository",
        Link:  *link,
    }
    client.BookmarkRepo().Create(ctx, bookmark)

    // Find GitHub links
    githubLinks, _ := client.BookmarkRepo().Query().
        Where(filter.Bookmark.Link.Domain().Equal("github.com")).
        All(ctx)

    // Find documentation pages
    docs, _ := client.BookmarkRepo().Query().
        Where(filter.Bookmark.Link.Path().Contains("/docs/")).
        All(ctx)

    // Find bookmarks with anchors
    withAnchors, _ := client.BookmarkRepo().Query().
        Where(filter.Bookmark.Link.Fragment().Zero(false)).
        All(ctx)

    // Find bookmarks with query strings
    withQuery, _ := client.BookmarkRepo().Query().
        Where(filter.Bookmark.Link.Query().Zero(false)).
        All(ctx)

    // Sort by domain
    sorted, _ := client.BookmarkRepo().Query().
        Order(by.Bookmark.Link.Domain().Asc()).
        All(ctx)

    // Bookmarks without favicons
    noFavicon, _ := client.BookmarkRepo().Query().
        Where(filter.Bookmark.Favicon.IsNil()).
        All(ctx)
}
```

## Filter Reference Table

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Equal(val)` | Exact match | Bool filter |
| `NotEqual(val)` | Not equal | Bool filter |
| `In(vals...)` | Value in set | Bool filter |
| `NotIn(vals...)` | Value not in set | Bool filter |
| `LessThan(val)` | Lexicographic < | Bool filter |
| `LessThanEqual(val)` | Lexicographic <= | Bool filter |
| `GreaterThan(val)` | Lexicographic > | Bool filter |
| `GreaterThanEqual(val)` | Lexicographic >= | Bool filter |
| `Domain()` | Extract domain | String filter |
| `Host()` | Extract host with port | String filter |
| `Port()` | Extract port | String filter |
| `Path()` | Extract path | String filter |
| `Query()` | Extract query string | String filter |
| `Fragment()` | Extract fragment | String filter |
| `Zero(bool)` | Check empty URL | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `IsNil()` | Is null (ptr) | Bool filter |
| `IsNotNil()` | Not null (ptr) | Bool filter |
