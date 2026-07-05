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
filter.User.Name.Equal("Alice")

// Not equal
filter.User.Name.NotEqual("Bob")

// Fuzzy match (pattern matching)
filter.User.Name.FuzzyMatch("Ali*")

// Negated fuzzy match
filter.User.Name.FuzzyNotMatch("*bot*")
```

### Set Membership

```go
// Value in set
filter.User.Role.In("admin", "moderator", "user")

// Value not in set
filter.User.Status.NotIn("banned", "deleted")
```

### Comparison Operations

```go
// Less than (lexicographic)
filter.User.Name.LessThan("M")

// Less than or equal
filter.User.Name.LessThanEqual("Alice")

// Greater than
filter.User.Name.GreaterThan("A")

// Greater than or equal
filter.User.Name.GreaterThanEqual("Alice")
```

### Substring Operations

```go
// Contains substring
filter.User.Email.Contains("@gmail")

// Starts with prefix
filter.User.Name.StartsWith("Dr.")

// Ends with suffix
filter.User.Email.EndsWith(".com")
```

### Case Transformation

```go
// Convert to lowercase then compare
filter.User.Email.Lowercase().Equal("alice@example.com")

// Convert to uppercase then compare
filter.User.Name.Uppercase().Equal("ALICE")
```

### String Manipulation

```go
// Trim whitespace
filter.User.Name.Trim().Equal("Alice")

// Reverse string
filter.User.Code.Reverse().Equal("cba")

// Repeat string
filter.User.Pattern.Repeat(3).Equal("abcabcabc")

// Replace substring
filter.User.Name.Replace("old", "new").Equal("new value")

// Get slug version
filter.Article.Title.Slug().Equal("hello-world")
```

### String Parts

```go
// Slice substring (start, length)
filter.User.Name.Slice(0, 3).Equal("Ali")

// Split into array
filter.User.Tags.Split(",").Contains("golang")

// Get words as array
filter.Article.Title.Words().Len().GreaterThan(5)

// Concatenate strings
filter.User.FullName.Concat(" Jr.").Equal("John Smith Jr.")

// Join array elements
filter.User.Tags.Join(", ").Contains("go")
```

### Length Operations

```go
// Get string length
filter.User.Name.Len().GreaterThan(3)

// Check minimum length
filter.User.Password.Len().GreaterThanEqual(8)

// Check maximum length
filter.User.Bio.Len().LessThanEqual(500)
```

### Validation Operations

These return boolean filters:

```go
// Email validation
filter.User.Email.IsEmail()

// URL validation
filter.User.Website.IsURL()

// Domain validation
filter.User.Domain.IsDomain()

// UUID validation
filter.User.ExternalID.IsUUID()

// Semantic version validation
filter.Package.Version.IsSemVer()

// DateTime format validation
filter.User.BirthDate.IsDateTime("%Y-%m-%d")

// IP address validation (any)
filter.Server.Address.IsIP()

// IPv4 validation
filter.Server.IPv4.IsIPv4()

// IPv6 validation
filter.Server.IPv6.IsIPv6()

// Latitude validation
filter.Location.Lat.IsLatitude()

// Longitude validation
filter.Location.Lng.IsLongitude()

// Alphabetic only
filter.User.Name.IsAlpha()

// Alphanumeric
filter.User.Username.IsAlphaNum()

// ASCII only
filter.User.Code.IsAscii()

// Hexadecimal
filter.User.Color.IsHexadecimal()

// Numeric string
filter.User.Phone.IsNumeric()
```

### Encoding Operations

```go
// Decode base64
filter.User.EncodedData.Base64Decode().Equal("decoded value")
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
filter.User.Nickname.IsNil()

// Check if not nil
filter.User.Nickname.IsNotNil()
```

### Zero Value Check

```go
// Is empty string
filter.User.Name.Zero(true)

// Is not empty string
filter.User.Name.Zero(false)
```

### Truth Conversion

```go
// Convert to boolean (non-empty = true)
filter.User.Name.Truth()
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
filter.User.Bio.Matches("software engineer")

// With highlighting
filter.User.Bio.Matches("golang").WithHighlights("<mark>", "</mark>")

// With explicit ref and offsets
filter.User.Bio.Matches("developer").Ref(0).WithOffsets()
```

Search is used with the query builder's `Search()` method:

```go
results, err := client.UserRepo().Query().
    Search(filter.User.Bio.Matches("golang developer")).
    AllMatches(ctx)
```

See [Full-Text Search](../querying/05_fulltext_search.md) for the complete guide.

## Method Chaining

String filters support extensive chaining:

```go
// Lowercase email domain check
filter.User.Email.Lowercase().EndsWith("@company.com")

// Trim and check length
filter.User.Bio.Trim().Len().LessThan(100)

// Slug and compare
filter.Article.Title.Slug().Equal("my-article-title")

// Complex chain
filter.User.Name.Trim().Lowercase().StartsWith("admin")
```

## Complete Example

```go
package main

import (
    "context"
    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/filter"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Find users with gmail addresses
    gmailUsers, _ := client.UserRepo().Query().
        Where(
            filter.User.Email.Lowercase().EndsWith("@gmail.com"),
            filter.User.Email.IsEmail(),
        ).
        Order(by.User.Name.Asc()).
        All(ctx)

    // Search by name pattern
    matches, _ := client.UserRepo().Query().
        Where(filter.User.Name.FuzzyMatch("*smith*")).
        All(ctx)

    // Find users with long bios
    longBios, _ := client.UserRepo().Query().
        Where(filter.User.Bio.Len().GreaterThan(500)).
        All(ctx)

    // Validate data quality
    invalidEmails, _ := client.UserRepo().Query().
        Where(filter.User.Email.IsEmail().Invert()).
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
