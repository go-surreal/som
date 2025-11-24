# Filters

Filters narrow down query results using type-safe conditions. SOM generates comprehensive filter operations for each field type.

## Basic Filtering

Use the generated `where` package:

```go
import "yourproject/gen/som/where"

users, err := client.UserRepo().Query().
    Filter(where.User.Email.Equal("john@example.com")).
    All(ctx)
```

## Multiple Filters (AND)

Multiple filters in a single `Filter()` call are combined with AND:

```go
users, err := client.UserRepo().Query().
    Filter(
        where.User.IsActive.IsTrue(),
        where.User.Age.GreaterThan(18),
        where.User.Email.Contains("@company.com"),
    ).
    All(ctx)
```

## Combining Filters (OR)

Use `where.Any()` for OR conditions:

```go
users, err := client.UserRepo().Query().
    Filter(
        where.Any(
            where.User.Role.Equal("admin"),
            where.User.Role.Equal("moderator"),
        ),
    ).
    All(ctx)
```

Use `where.All()` explicitly for AND:

```go
users, err := client.UserRepo().Query().
    Filter(
        where.All(
            where.User.IsActive.IsTrue(),
            where.User.Age.GreaterThan(18),
        ),
    ).
    All(ctx)
```

## Base Operations (All Types)

Available on all comparable field types:

| Operation | Description | Example |
|-----------|-------------|---------|
| `Equal(val)` | Equals | `where.User.Name.Equal("John")` |
| `NotEqual(val)` | Not equals | `where.User.Name.NotEqual("John")` |
| `In(vals...)` | In list | `where.User.Status.In("active", "pending")` |
| `NotIn(vals...)` | Not in list | `where.User.Status.NotIn("deleted")` |

## Comparison Operations (Numeric, Time, String)

| Operation | Description |
|-----------|-------------|
| `LessThan(val)` | Less than |
| `LessThanOrEqual(val)` | Less than or equal |
| `GreaterThan(val)` | Greater than |
| `GreaterThanOrEqual(val)` | Greater than or equal |

```go
where.User.Age.GreaterThan(18)
where.User.Age.LessThanOrEqual(65)
where.User.CreatedAt.GreaterThan(lastWeek)
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
where.User.Email.Contains("@gmail")
where.User.Name.StartsWith("John")
where.User.Email.EndsWith(".com")
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
where.User.Email.IsEmail()
where.User.Website.IsURL()
where.User.ExternalID.IsUUID()
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
where.User.Email.Lowercase().Equal("john@example.com")
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
where.User.Name.Len().GreaterThan(3)
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
where.User.Age.Add(5).GreaterThan(25)

// Absolute value of balance > 100
where.Account.Balance.Abs().GreaterThan(100)
```

## Boolean Operations

| Operation | Description |
|-----------|-------------|
| `Equal(bool)` | Equals value |
| `IsTrue()` | Is true |
| `IsFalse()` | Is false |

```go
where.User.IsActive.IsTrue()
where.User.IsDeleted.IsFalse()
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
where.User.CreatedAt.After(time.Now().AddDate(0, 0, -7))

// Created this year
where.User.CreatedAt.After(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
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
where.Session.Duration.After(time.Hour)
```

## Pointer/Optional Operations

| Operation | Description |
|-----------|-------------|
| `IsNil()` | Is null |
| `IsNotNil()` | Is not null |

```go
// Find soft-deleted users
where.User.DeletedAt.IsNotNil()

// Find users without avatar
where.User.AvatarURL.IsNil()
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
where.Post.Tags.NotEmpty()

// Has specific tag
where.Post.Tags.Contains("golang")

// Has any of these tags
where.Post.Tags.ContainsAny("golang", "rust", "python")

// Has all required tags
where.Post.Tags.ContainsAll("featured", "published")

// More than 5 tags
where.Post.Tags.Length().GreaterThan(5)
```

## Nested Field Filters

Filter on embedded struct fields:

```go
// Filter by nested city
where.User.Address.City.Equal("Berlin")

// Deeply nested
where.User.Address.Coordinates.Lat.GreaterThan(52.0)
```

## Enum Filters

```go
where.User.Status.Equal(model.StatusActive)
where.User.Status.In(model.StatusActive, model.StatusPending)
where.User.Status.NotEqual(model.StatusDeleted)
```

## Complex Example

```go
users, err := client.UserRepo().Query().
    Filter(
        // Active users
        where.User.IsActive.IsTrue(),

        // Created this month
        where.User.CreatedAt.After(startOfMonth),

        // Has email from allowed domains
        where.Any(
            where.User.Email.EndsWith("@company.com"),
            where.User.Email.EndsWith("@partner.com"),
        ),

        // Age between 18 and 65
        where.User.Age.GreaterThanOrEqual(18),
        where.User.Age.LessThanOrEqual(65),

        // Has at least one role
        where.User.Roles.NotEmpty(),

        // In Berlin
        where.User.Address.City.Equal("Berlin"),
    ).
    Order(by.User.CreatedAt.Desc()).
    Limit(100).
    All(ctx)
```
