# Filters

Filters narrow down query results using type-safe conditions. SOM generates comprehensive filter operations for each field type.

## Basic Filtering

Use the generated `filter` package:

```go
import "yourproject/gen/som/filter"

users, err := client.UserRepo().Query().
    Where(filter.User.Email.Equal("john@example.com")).
    All(ctx)
```

## Multiple Filters (AND)

Multiple filters in a single `Where()` call are combined with AND:

```go
users, err := client.UserRepo().Query().
    Where(
        filter.User.IsActive.IsTrue(),
        filter.User.Age.GreaterThan(18),
        filter.User.Email.Contains("@company.com"),
    ).
    All(ctx)
```

## Combining Filters (OR)

Use `filter.Any()` for OR conditions:

```go
users, err := client.UserRepo().Query().
    Where(
        filter.Any(
            filter.User.Role.Equal("admin"),
            filter.User.Role.Equal("moderator"),
        ),
    ).
    All(ctx)
```

Use `filter.All()` explicitly for AND:

```go
users, err := client.UserRepo().Query().
    Where(
        filter.All(
            filter.User.IsActive.IsTrue(),
            filter.User.Age.GreaterThan(18),
        ),
    ).
    All(ctx)
```

## Base Operations (All Types)

Available on all comparable field types:

| Operation | Description | Example |
|-----------|-------------|---------|
| `Equal(val)` | Equals | `filter.User.Name.Equal("John")` |
| `NotEqual(val)` | Not equals | `filter.User.Name.NotEqual("John")` |
| `In(vals...)` | In list | `filter.User.Status.In("active", "pending")` |
| `NotIn(vals...)` | Not in list | `filter.User.Status.NotIn("deleted")` |

## Comparison Operations (Numeric, Time, String)

| Operation | Description |
|-----------|-------------|
| `LessThan(val)` | Less than |
| `LessThanOrEqual(val)` | Less than or equal |
| `GreaterThan(val)` | Greater than |
| `GreaterThanOrEqual(val)` | Greater than or equal |

```go
filter.User.Age.GreaterThan(18)
filter.User.Age.LessThanOrEqual(65)
filter.User.CreatedAt.GreaterThan(lastWeek)
```

## String Operations

Strings have the most extensive filter operations:

### Pattern Matching

| Operation | Description | SurrealQL |
|-----------|-------------|-----------|
| `Contains(s)` | Contains substring | `CONTAINS` |
| `StartsWith(s)` | Starts with | `string::startsWith()` |
| `EndsWith(s)` | Ends with | `string::endsWith()` |
| `FuzzyMatch(s)` | Fuzzy match | `~` |
| `FuzzyNotMatch(s)` | Fuzzy not match | `!~` |

```go
filter.User.Email.Contains("@gmail")
filter.User.Name.StartsWith("John")
filter.User.Email.EndsWith(".com")
```

### Validation

| Operation | Description |
|-----------|-------------|
| `IsAlpha()` | Only letters |
| `IsAlphaNum()` | Letters and numbers |
| `IsAscii()` | ASCII characters |
| `IsEmail()` | Valid email format |
| `IsDomain()` | Valid domain |
| `IsURL()` | Valid URL |
| `IsIP()` | Valid IP address |
| `IsIPv4()` | Valid IPv4 |
| `IsIPv6()` | Valid IPv6 |
| `IsLatitude()` | Valid latitude |
| `IsLongitude()` | Valid longitude |
| `IsNumeric()` | Numeric string |
| `IsSemVer()` | Semantic version |
| `IsUUID()` | Valid UUID |
| `IsDateTime(format)` | Valid datetime |

```go
filter.User.Email.IsEmail()
filter.User.Website.IsURL()
filter.User.ExternalID.IsUUID()
```

### Transformation (for comparison)

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Lowercase()` | Convert to lowercase | String filter |
| `Uppercase()` | Convert to uppercase | String filter |
| `Trim()` | Remove whitespace | String filter |
| `Slug()` | Convert to slug | String filter |
| `Reverse()` | Reverse string | String filter |

```go
// Compare lowercase version
filter.User.Email.Lowercase().Equal("john@example.com")
```

### String Functions

| Operation | Description |
|-----------|-------------|
| `Len()` | String length |
| `Split(sep)` | Split into array |
| `Words()` | Split into words |
| `Slice(start, end)` | Substring |
| `Replace(old, new)` | Replace substring |
| `Repeat(n)` | Repeat n times |
| `Concat(strings...)` | Concatenate |
| `Join(strings...)` | Join with separator |

```go
// Filter by string length
filter.User.Name.Len().GreaterThan(3)
```

## Numeric Operations

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Add(n)` | Add value | Numeric filter |
| `Sub(n)` | Subtract value | Numeric filter |
| `Mul(n)` | Multiply | Numeric filter |
| `Div(n)` | Divide | Numeric filter |
| `Raise(n)` | Power | Numeric filter |
| `Abs()` | Absolute value | Numeric filter |

```go
// Age + 5 > 25
filter.User.Age.Add(5).GreaterThan(25)

// Absolute value of balance > 100
filter.Account.Balance.Abs().GreaterThan(100)
```

## Boolean Operations

| Operation | Description |
|-----------|-------------|
| `Equal(bool)` | Equals value |
| `IsTrue()` | Is true |
| `IsFalse()` | Is false |

```go
filter.User.IsActive.IsTrue()
filter.User.IsDeleted.IsFalse()
```

## Time Operations

| Operation | Description |
|-----------|-------------|
| `Before(time)` | Before time |
| `BeforeOrEqual(time)` | Before or equal |
| `After(time)` | After time |
| `AfterOrEqual(time)` | After or equal |
| `Add(duration)` | Add duration |
| `Sub(duration)` | Subtract duration |
| `Floor(duration)` | Floor to duration |
| `Round(duration)` | Round to duration |
| `Format(format)` | Format as string |

```go
// Created in last 7 days
filter.User.CreatedAt.After(time.Now().AddDate(0, 0, -7))

// Created this year
filter.User.CreatedAt.After(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
```

## Duration Operations

| Operation | Description |
|-----------|-------------|
| `Before(duration)` | Shorter than |
| `After(duration)` | Longer than |
| `Add(duration)` | Add durations |
| `Sub(duration)` | Subtract durations |

```go
// Session longer than 1 hour
filter.Session.Duration.After(time.Hour)
```

## Pointer/Optional Operations

| Operation | Description |
|-----------|-------------|
| `IsNil()` | Is null |
| `IsNotNil()` | Is not null |

```go
// Find soft-deleted users
filter.User.DeletedAt.IsNotNil()

// Find users without avatar
filter.User.AvatarURL.IsNil()
```

## Slice Operations

| Operation | Description |
|-----------|-------------|
| `Length()` / `Len()` | Array length |
| `Contains(val)` | Contains element |
| `ContainsAll(vals...)` | Contains all elements |
| `ContainsAny(vals...)` | Contains any element |
| `ContainsNone(vals...)` | Contains no elements |
| `Empty()` | Is empty |
| `NotEmpty()` | Is not empty |
| `Intersects(vals...)` | Has common elements |
| `Inside(vals...)` | All elements in list |

```go
// Has at least one tag
filter.Post.Tags.NotEmpty()

// Has specific tag
filter.Post.Tags.Contains("golang")

// Has any of these tags
filter.Post.Tags.ContainsAny("golang", "rust", "python")

// Has all required tags
filter.Post.Tags.ContainsAll("featured", "published")

// More than 5 tags
filter.Post.Tags.Length().GreaterThan(5)
```

## Nested Field Filters

Filter on embedded struct fields:

```go
// Filter by nested city
filter.User.Address.City.Equal("Berlin")

// Deeply nested
filter.User.Address.Coordinates.Lat.GreaterThan(52.0)
```

## Enum Filters

```go
filter.User.Status.Equal(model.StatusActive)
filter.User.Status.In(model.StatusActive, model.StatusPending)
filter.User.Status.NotEqual(model.StatusDeleted)
```

## Complex Example

```go
users, err := client.UserRepo().Query().
    Where(
        // Active users
        filter.User.IsActive.IsTrue(),

        // Created this month
        filter.User.CreatedAt.After(startOfMonth),

        // Has email from allowed domains
        filter.Any(
            filter.User.Email.EndsWith("@company.com"),
            filter.User.Email.EndsWith("@partner.com"),
        ),

        // Age between 18 and 65
        filter.User.Age.GreaterThanOrEqual(18),
        filter.User.Age.LessThanOrEqual(65),

        // Has at least one role
        filter.User.Roles.NotEmpty(),

        // In Berlin
        filter.User.Address.City.Equal("Berlin"),
    ).
    Order(by.User.CreatedAt.Desc()).
    Limit(100).
    All(ctx)
```

## Combining with Full-Text Search

Filters can be combined with full-text search conditions:

```go
results, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang tutorial")).
    Where(
        filter.Article.Published.IsTrue(),
        filter.Article.Category.Equal("programming"),
    ).
    AllMatches(ctx)
```

The search and filter conditions are combined with AND in the WHERE clause. See [Full-Text Search](05_fulltext_search.md) for the complete search guide.
