# String Type

The string type is the most feature-rich type in SOM, with over 50 filter operations for pattern matching, transformation, and validation.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `string` / `*string` |
| Database Schema | `string` / `option<string>` |
| CBOR Encoding | Direct |
| Sortable | Yes |

## Definition

```go
type User struct {
    som.Node

    Name        string   // Required
    Bio         string   // Required
    Nickname    *string  // Optional (nullable)
}
```

## Schema

Generated SurrealDB schema:

```surql
DEFINE FIELD name ON user TYPE string;
DEFINE FIELD bio ON user TYPE string;
DEFINE FIELD nickname ON user TYPE option<string>;
```

## Filter Operations

### Equality Operations

```go
// Exact match
where.User.Name.Equal("Alice")

// Not equal
where.User.Name.NotEqual("Bob")

// Fuzzy match (pattern matching)
where.User.Name.FuzzyMatch("Ali*")

// Negated fuzzy match
where.User.Name.FuzzyNotMatch("*bot*")
```

### Set Membership

```go
// Value in set
where.User.Role.In("admin", "moderator", "user")

// Value not in set
where.User.Status.NotIn("banned", "deleted")
```

### Comparison Operations

```go
// Less than (lexicographic)
where.User.Name.LessThan("M")

// Less than or equal
where.User.Name.LessThanEqual("Alice")

// Greater than
where.User.Name.GreaterThan("A")

// Greater than or equal
where.User.Name.GreaterThanEqual("Alice")
```

### Substring Operations

```go
// Contains substring
where.User.Email.Contains("@gmail")

// Starts with prefix
where.User.Name.StartsWith("Dr.")

// Ends with suffix
where.User.Email.EndsWith(".com")
```

### Case Transformation

```go
// Convert to lowercase then compare
where.User.Email.Lowercase().Equal("alice@example.com")

// Convert to uppercase then compare
where.User.Name.Uppercase().Equal("ALICE")
```

### String Manipulation

```go
// Trim whitespace
where.User.Name.Trim().Equal("Alice")

// Reverse string
where.User.Code.Reverse().Equal("cba")

// Repeat string
where.User.Pattern.Repeat(3).Equal("abcabcabc")

// Replace substring
where.User.Name.Replace("old", "new").Equal("new value")

// Get slug version
where.Article.Title.Slug().Equal("hello-world")
```

### String Parts

```go
// Slice substring (start, length)
where.User.Name.Slice(0, 3).Equal("Ali")

// Split into array
where.User.Tags.Split(",").Contains("golang")

// Get words as array
where.Article.Title.Words().Len().GreaterThan(5)

// Concatenate strings
where.User.FullName.Concat(" Jr.").Equal("John Smith Jr.")

// Join array elements
where.User.Tags.Join(", ").Contains("go")
```

### Length Operations

```go
// Get string length
where.User.Name.Len().GreaterThan(3)

// Check minimum length
where.User.Password.Len().GreaterThanEqual(8)

// Check maximum length
where.User.Bio.Len().LessThanEqual(500)
```

### Validation Operations

These return boolean filters:

```go
// Email validation
where.User.Email.IsEmail()

// URL validation
where.User.Website.IsURL()

// Domain validation
where.User.Domain.IsDomain()

// UUID validation
where.User.ExternalID.IsUUID()

// Semantic version validation
where.Package.Version.IsSemVer()

// DateTime format validation
where.User.BirthDate.IsDateTime("%Y-%m-%d")

// IP address validation (any)
where.Server.Address.IsIP()

// IPv4 validation
where.Server.IPv4.IsIPv4()

// IPv6 validation
where.Server.IPv6.IsIPv6()

// Latitude validation
where.Location.Lat.IsLatitude()

// Longitude validation
where.Location.Lng.IsLongitude()

// Alphabetic only
where.User.Name.IsAlpha()

// Alphanumeric
where.User.Username.IsAlphaNum()

// ASCII only
where.User.Code.IsAscii()

// Hexadecimal
where.User.Color.IsHexadecimal()

// Numeric string
where.User.Phone.IsNumeric()
```

### Encoding Operations

```go
// Decode base64
where.User.EncodedData.Base64Decode().Equal("decoded value")
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
where.User.Nickname.IsNil()

// Check if not nil
where.User.Nickname.IsNotNil()
```

### Zero Value Check

```go
// Is empty string
where.User.Name.Zero(true)

// Is not empty string
where.User.Name.Zero(false)
```

### Truth Conversion

```go
// Convert to boolean (non-empty = true)
where.User.Name.Truth()
```

## Sorting

```go
// Ascending (A-Z)
query.Order(by.User.Name.Asc())

// Descending (Z-A)
query.Order(by.User.Name.Desc())

// Case-insensitive sort (transform first)
query.Order(by.User.Name.Lowercase().Asc())
```

## Full-Text Search

String fields with fulltext indexes support search operations:

```go
// Basic search
where.User.Bio.Matches("software engineer")

// With highlighting
where.User.Bio.Matches("golang").WithHighlights("<mark>", "</mark>")

// With explicit ref and offsets
where.User.Bio.Matches("developer").Ref(0).WithOffsets()
```

Search is used with the query builder's `Search()` method:

```go
results, err := client.UserRepo().Query().
    Search(where.User.Bio.Matches("golang developer")).
    AllMatches(ctx)
```

See [Full-Text Search](../querying/05_fulltext_search.md) for the complete guide.

## Method Chaining

String filters support extensive chaining:

```go
// Lowercase email domain check
where.User.Email.Lowercase().EndsWith("@company.com")

// Trim and check length
where.User.Bio.Trim().Len().LessThan(100)

// Slug and compare
where.Article.Title.Slug().Equal("my-article-title")

// Complex chain
where.User.Name.Trim().Lowercase().StartsWith("admin")
```

## Complete Example

```go
package main

import (
    "context"
    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/where"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Find users with gmail addresses
    gmailUsers, _ := client.UserRepo().Query().
        Filter(
            where.User.Email.Lowercase().EndsWith("@gmail.com"),
            where.User.Email.IsEmail(),
        ).
        Order(by.User.Name.Asc()).
        All(ctx)

    // Search by name pattern
    matches, _ := client.UserRepo().Query().
        Filter(where.User.Name.FuzzyMatch("*smith*")).
        All(ctx)

    // Find users with long bios
    longBios, _ := client.UserRepo().Query().
        Filter(where.User.Bio.Len().GreaterThan(500)).
        All(ctx)

    // Validate data quality
    invalidEmails, _ := client.UserRepo().Query().
        Filter(where.User.Email.IsEmail().Invert()).
        All(ctx)
}
```

## Filter Reference Table

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Equal(val)` | Exact match | Bool filter |
| `NotEqual(val)` | Not equal | Bool filter |
| `FuzzyMatch(pattern)` | Pattern match with wildcards | Bool filter |
| `FuzzyNotMatch(pattern)` | Negated pattern match | Bool filter |
| `In(vals...)` | Value in set | Bool filter |
| `NotIn(vals...)` | Value not in set | Bool filter |
| `LessThan(val)` | Lexicographic < | Bool filter |
| `LessThanEqual(val)` | Lexicographic <= | Bool filter |
| `GreaterThan(val)` | Lexicographic > | Bool filter |
| `GreaterThanEqual(val)` | Lexicographic >= | Bool filter |
| `Contains(sub)` | Contains substring | Bool filter |
| `StartsWith(prefix)` | Starts with prefix | Bool filter |
| `EndsWith(suffix)` | Ends with suffix | Bool filter |
| `Lowercase()` | Convert to lowercase | String filter |
| `Uppercase()` | Convert to uppercase | String filter |
| `Trim()` | Remove whitespace | String filter |
| `Reverse()` | Reverse characters | String filter |
| `Repeat(n)` | Repeat n times | String filter |
| `Replace(old, new)` | Replace substring | String filter |
| `Slug()` | URL-friendly slug | String filter |
| `Slice(start, len)` | Substring | String filter |
| `Split(sep)` | Split to array | Slice filter |
| `Words()` | Split into words | Slice filter |
| `Concat(val)` | Append string | String filter |
| `Join(sep)` | Join array | String filter |
| `Len()` | String length | Numeric filter |
| `IsEmail()` | Validate email | Bool filter |
| `IsURL()` | Validate URL | Bool filter |
| `IsDomain()` | Validate domain | Bool filter |
| `IsUUID()` | Validate UUID | Bool filter |
| `IsSemVer()` | Validate semver | Bool filter |
| `IsDateTime(fmt)` | Validate datetime | Bool filter |
| `IsIP()` | Validate IP | Bool filter |
| `IsIPv4()` | Validate IPv4 | Bool filter |
| `IsIPv6()` | Validate IPv6 | Bool filter |
| `IsLatitude()` | Validate latitude | Bool filter |
| `IsLongitude()` | Validate longitude | Bool filter |
| `IsAlpha()` | Alphabetic only | Bool filter |
| `IsAlphaNum()` | Alphanumeric | Bool filter |
| `IsAscii()` | ASCII only | Bool filter |
| `IsHexadecimal()` | Hex string | Bool filter |
| `IsNumeric()` | Numeric string | Bool filter |
| `Base64Decode()` | Decode base64 | String filter |
| `Zero(bool)` | Check empty | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `IsNil()` | Is null (ptr) | Bool filter |
| `IsNotNil()` | Is not null (ptr) | Bool filter |
