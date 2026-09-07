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

Filters live in the generated `filter` package, one accessor per model and field.

## Multiple Filters (AND)

Multiple filters in a single `Where()` call are combined with AND:

```go
users, err := client.UserRepo().Query().
    Where(
        filter.User.IsActive.True(),
        filter.User.Age.GreaterThan(18),
        filter.User.Email.Contains("@company.com").True(),
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

### Negation

`filter.Not()` inverts a filter. The model type has to be given explicitly:

```go
users, err := client.UserRepo().Query().
    Where(
        filter.Not[model.User](
            filter.User.Status.Equal("archived"),
        ),
    ).
    All(ctx)
```

> `NotEqual` and `NotIn` deliberately exclude records where the field is unset (`NONE`).
> Use `filter.Not(...)` combined with `Nil(true)` if you want those included.

Use `filter.All()` explicitly for AND:

```go
users, err := client.UserRepo().Query().
    Where(
        filter.All(
            filter.User.IsActive.True(),
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
| `In(vals []T)` | In list | `filter.User.Status.In([]string{"active", "pending"})` |
| `NotIn(vals []T)` | Not in list | `filter.User.Status.NotIn([]string{"deleted"})` |

## Comparison Operations (Numeric, Time, String)

| Operation | Description |
|-----------|-------------|
| `LessThan(val)` | Less than |
| `LessThanEqual(val)` | Less than or equal |
| `GreaterThan(val)` | Greater than |
| `GreaterThanEqual(val)` | Greater than or equal |
| `Between(from, to)` | Value within range (inclusive by default) |

```go
filter.User.Age.GreaterThan(18)
filter.User.Age.LessThanEqual(65)
filter.User.CreatedAt.GreaterThan(lastWeek)
```

### Between Filter

The `Between` filter checks whether a field value falls within a range. By default, both bounds are inclusive:

```go
// Age between 18 and 65 (inclusive)
filter.User.Age.Between(18, 65)
```

Control bound inclusivity with chainable methods:

```go
// Exclusive lower bound: 18 < age <= 65
filter.User.Age.Between(18, 65).FromExclusive()

// Exclusive upper bound: 18 <= age < 65
filter.User.Age.Between(18, 65).ToExclusive()

// Both exclusive: 18 < age < 65
filter.User.Age.Between(18, 65).BothExclusive()
```

Works with any comparable type including time:

```go
filter.User.CreatedAt.Between(startDate, endDate)
```

## String Operations

Strings have the most extensive filter operations:

### Pattern Matching

| Operation | Returns | Description |
|-----------|---------|-------------|
| `Contains(s)` | bool expression | Contains substring |
| `StartsWith(s)` | bool expression | Starts with |
| `EndsWith(s)` | bool expression | Ends with |
| `Matches(regex)` | bool expression | Matches a regular expression |
| `FuzzyMatch(s)` | filter | Fuzzy match (`~`) |
| `FuzzyNotMatch(s)` | filter | Fuzzy not match (`!~`) |

Operations that return a **bool expression** are not filters yet — finish them with `.True()`,
`.False()` or `.Is(bool)`:

```go
filter.User.Email.Contains("@gmail").True()
filter.User.Name.StartsWith("John").True()
filter.User.Email.EndsWith(".com").False()   // does NOT end with ".com"

// FuzzyMatch is already a filter
filter.User.Name.FuzzyMatch("jon")
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
| `IsULID()` | Valid ULID |
| `IsHexadecimal()` | Hexadecimal string |
| `IsDateTime(format)` | Valid datetime |

All validation operations return a bool expression, so they also need `.True()` / `.False()`:

```go
filter.User.Email.IsEmail().True()
filter.User.Website.IsURL().True()
filter.User.ExternalID.IsUUID().False()
```

Additional checks: `IsHexadecimal()`, `IsULID()`.

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
| `Len()` | String length (numeric) |
| `Split(sep)` | Split into a slice filter |
| `Words()` | Split into words (slice filter) |
| `Slice(start, end)` | Substring |
| `Replace(old, new)` | Replace substring |
| `Repeat(n)` | Repeat n times |
| `Capitalize()` | Capitalize |
| `Concat(vals...)` | Concatenate (bool expression) |
| `Join(vals...)` | Join with separator |

String similarity and distance helpers are available as well: `SimilarityFuzzy`,
`SimilarityJaro`, `SimilarityJaroWinkler`, `DistanceLevenshtein`, `DistanceHamming`,
`DistanceDamerauLevenshtein`, `DistanceOsa` and their normalized variants.

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
| `Is(val)` | Equals the given value |
| `True()` | Is true |
| `False()` | Is false |
| `Invert()` | Negate the expression (returns a bool expression) |

```go
filter.User.IsActive.True()
filter.User.IsDeleted.False()
filter.User.IsActive.Is(false)
filter.User.IsActive.Invert().True()   // same as Is(false)
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
| `Nil(true)` | Is `NONE` or `NULL` |
| `Nil(false)` | Is set |

Pointer fields additionally expose all operations of their underlying type.

```go
// Find soft-deleted users
filter.User.DeletedAt.Nil(false)

// Find users without avatar
filter.User.AvatarURL.Nil(true)
```

## Slice Operations

| Operation | Description |
|-----------|-------------|
| `Len()` | Array length (numeric) |
| `Contains(val)` | Contains element |
| `ContainsNot(val)` | Does not contain element |
| `ContainsAll(vals []E)` | Contains all elements |
| `ContainsAny(vals []E)` | Contains any element |
| `ContainsNone(vals []E)` | Contains no elements |
| `Empty(is bool)` | Is / is not empty |
| `IsEmpty()` / `NotEmpty()` | Convenience wrappers around `Empty` |
| `AnyIn(vals []E)` / `AllIn(vals []E)` / `NoneIn(vals []E)` | Membership of the slice values |
| `AnyEqual(val)` / `AllEqual(val)` | Element comparison |
| `At(i)` / `First()` / `Last()` / `Min()` / `Max()` | Element access, returns the element filter |
| `Distinct()`, `Reverse()`, `SortAsc()`, `SortDesc()`, `Union(vals)`, `Intersect(vals)`, `Diff(vals)` | Slice transformations |

```go
// Has at least one tag
filter.Post.Tags.NotEmpty()

// Has specific tag
filter.Post.Tags.Contains("golang").True()

// Has any of these tags
filter.Post.Tags.ContainsAny([]string{"golang", "rust", "python"})

// Has all required tags
filter.Post.Tags.ContainsAll([]string{"featured", "published"})

// More than 5 tags
filter.Post.Tags.Len().GreaterThan(5)
```

## Nested Field Filters

Filter on embedded struct fields:

```go
// Filter by nested city
filter.User.Address().City.Equal("Berlin")

// Deeply nested
filter.User.Address().Coordinates.Lat.GreaterThan(52.0)
```

## Enum Filters

```go
filter.User.Status.Equal(model.StatusActive)
filter.User.Status.In([]model.Status{model.StatusActive, model.StatusPending})
filter.User.Status.NotEqual(model.StatusDeleted)
```

## Complex Example

```go
users, err := client.UserRepo().Query().
    Where(
        // Active users
        filter.User.IsActive.True(),

        // Created this month
        filter.User.CreatedAt.After(startOfMonth),

        // Has email from allowed domains
        filter.Any(
            filter.User.Email.EndsWith("@company.com").True(),
            filter.User.Email.EndsWith("@partner.com").True(),
        ),

        // Age between 18 and 65
        filter.User.Age.Between(18, 65),

        // Has at least one role
        filter.User.Roles.NotEmpty(),

        // In Berlin
        filter.User.Address().City.Equal("Berlin"),
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
        filter.Article.Published.True(),
        filter.Article.Category.Equal("programming"),
    ).
    AllMatches(ctx)
```

The search and filter conditions are combined with AND in the WHERE clause. See [Full-Text Search](05_fulltext_search.md) for the complete search guide.
