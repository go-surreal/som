# Time Type

The time type handles timestamps using Go's `time.Time` with special CBOR encoding for nanosecond precision.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `time.Time` / `*time.Time` |
| Database Schema | `datetime` / `option<datetime>` |
| CBOR Encoding | Tag 12 with `[unix_seconds, nanoseconds]` |
| Sortable | Yes |

## CBOR Encoding

Time values are encoded with CBOR Tag 12 as a two-element array:

```
Tag 12: [unix_seconds (int64), nanoseconds (int64)]
```

This provides:
- Full nanosecond precision
- Proper round-tripping with SurrealDB
- Handling of empty arrays for optional fields

## Definition

```go
type Event struct {
    som.Node

    StartTime   time.Time   // Required
    EndTime     time.Time   // Required
    CancelledAt *time.Time  // Optional
}
```

### Automatic Timestamps

Use `som.Timestamps` for automatic tracking:

```go
type User struct {
    som.Node
    som.Timestamps  // Adds CreatedAt and UpdatedAt

    Name string
}
```

- `CreatedAt` - Set automatically on create, readonly
- `UpdatedAt` - Updated automatically on every save

## Schema

Generated SurrealDB schema:

```surql
DEFINE FIELD start_time ON event TYPE datetime;
DEFINE FIELD end_time ON event TYPE datetime;
DEFINE FIELD cancelled_at ON event TYPE option<datetime>;

-- With som.Timestamps:
DEFINE FIELD created_at ON user TYPE datetime
    VALUE $before OR time::now() READONLY;
DEFINE FIELD updated_at ON user TYPE datetime
    VALUE time::now();
```

## Filter Operations

### Equality Operations

```go
// Exact match
where.Event.StartTime.Equal(targetTime)

// Not equal
where.Event.StartTime.NotEqual(excludeTime)
```

### Set Membership

```go
// Value in set
where.Event.StartTime.In(time1, time2, time3)

// Value not in set
where.Event.StartTime.NotIn(excludedTimes...)
```

### Comparison Operations

Time provides intuitive comparison aliases:

```go
// Before (less than)
where.Event.StartTime.Before(deadline)

// Before or equal
where.Event.StartTime.BeforeOrEqual(deadline)

// After (greater than)
where.Event.StartTime.After(startDate)

// After or equal
where.Event.StartTime.AfterOrEqual(startDate)

// Standard comparison also works:
where.Event.StartTime.LessThan(deadline)
where.Event.StartTime.GreaterThan(startDate)
```

### Arithmetic Operations

```go
// Add duration
where.Event.StartTime.Add(2 * time.Hour).Before(deadline)

// Subtract duration
where.Event.EndTime.Sub(30 * time.Minute).After(now)
```

### Component Extraction

Extract individual time components as numeric filters:

```go
// Year
where.Event.StartTime.Year().Equal(2024)

// Month (1-12)
where.Event.StartTime.Month().Equal(12)

// Day of month (1-31)
where.Event.StartTime.Day().Equal(25)

// Hour (0-23)
where.Event.StartTime.Hour().GreaterThanEqual(9)

// Minute (0-59)
where.Event.StartTime.Minute().Equal(0)

// Second (0-59)
where.Event.StartTime.Second().Equal(0)

// Nanosecond
where.Event.StartTime.Nano().Equal(0)

// Microseconds since epoch
where.Event.StartTime.Micros().GreaterThan(0)

// Milliseconds since epoch
where.Event.StartTime.Millis().GreaterThan(0)

// Unix timestamp (seconds)
where.Event.StartTime.Unix().GreaterThan(0)

// Day of week (0=Sunday, 6=Saturday)
where.Event.StartTime.Weekday().Equal(1)  // Monday

// Week of year (1-53)
where.Event.StartTime.Week().Equal(52)

// Day of year (1-366)
where.Event.StartTime.YearDay().LessThan(100)
```

### Grouping

Group times by unit:

```go
// Group by year
where.Event.StartTime.Group(time.Year).Equal(groupedTime)

// Group by month
where.Event.StartTime.Group(time.Month)

// Group by day
where.Event.StartTime.Group(time.Day)

// Group by hour
where.Event.StartTime.Group(time.Hour)

// Group by minute
where.Event.StartTime.Group(time.Minute)

// Group by second
where.Event.StartTime.Group(time.Second)
```

### Rounding Operations

```go
// Floor to hour boundary
where.Event.StartTime.Floor(time.Hour).Equal(hourStart)

// Round to nearest hour
where.Event.StartTime.Round(time.Hour).Equal(nearestHour)

// Floor to day
where.Event.StartTime.Floor(24 * time.Hour).Equal(dayStart)
```

### Formatting

Format time as string for comparison:

```go
// Format and compare
where.Event.StartTime.Format("%Y-%m-%d").Equal("2024-12-25")

// SurrealDB format specifiers:
// %Y - 4-digit year
// %m - 2-digit month
// %d - 2-digit day
// %H - 24-hour
// %M - minute
// %S - second
```

### Validation

```go
// Check if leap year
where.Event.StartTime.IsLeapYear().True()
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
where.Event.CancelledAt.IsNil()

// Check if not nil
where.Event.CancelledAt.IsNotNil()
```

### Zero Value Check

```go
// Is zero time
where.Event.StartTime.Zero(true)

// Is not zero time
where.Event.StartTime.Zero(false)
```

## Sorting

```go
// Ascending (oldest first)
query.Order(by.Event.StartTime.Asc())

// Descending (newest first)
query.Order(by.Event.StartTime.Desc())

// Multiple sort fields
query.Order(
    by.Event.StartTime.Desc(),
    by.Event.Name.Asc(),
)
```

## Method Chaining

Time filters support powerful chaining:

```go
// Events starting in morning hours
where.Event.StartTime.Hour().GreaterThanEqual(6)
where.Event.StartTime.Hour().LessThan(12)

// Events in Q4 2024
where.Event.StartTime.Year().Equal(2024)
where.Event.StartTime.Month().GreaterThanEqual(10)

// Events 2 hours from now
where.Event.StartTime.Add(2 * time.Hour).Before(time.Now())
```

## Common Patterns

### Date Range Queries

```go
startOfDay := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
endOfDay := startOfDay.Add(24 * time.Hour)

events, _ := client.EventRepo().Query().
    Filter(
        where.Event.StartTime.AfterOrEqual(startOfDay),
        where.Event.StartTime.Before(endOfDay),
    ).
    All(ctx)
```

### Recent Records

```go
oneWeekAgo := time.Now().AddDate(0, 0, -7)

recentEvents, _ := client.EventRepo().Query().
    Filter(where.Event.CreatedAt.After(oneWeekAgo)).
    Order(by.Event.CreatedAt.Desc()).
    All(ctx)
```

### Upcoming Events

```go
now := time.Now()
nextMonth := now.AddDate(0, 1, 0)

upcoming, _ := client.EventRepo().Query().
    Filter(
        where.Event.StartTime.After(now),
        where.Event.StartTime.Before(nextMonth),
    ).
    Order(by.Event.StartTime.Asc()).
    All(ctx)
```

## Complete Example

```go
package main

import (
    "context"
    "time"
    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/where"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})
    now := time.Now()

    // Events starting today
    startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
    todayEvents, _ := client.EventRepo().Query().
        Filter(
            where.Event.StartTime.AfterOrEqual(startOfToday),
            where.Event.StartTime.Before(startOfToday.Add(24*time.Hour)),
        ).
        Order(by.Event.StartTime.Asc()).
        All(ctx)

    // Cancelled events
    cancelled, _ := client.EventRepo().Query().
        Filter(where.Event.CancelledAt.IsNotNil()).
        All(ctx)

    // Morning events (9 AM - 12 PM)
    morningEvents, _ := client.EventRepo().Query().
        Filter(
            where.Event.StartTime.Hour().GreaterThanEqual(9),
            where.Event.StartTime.Hour().LessThan(12),
        ).
        All(ctx)

    // Events in December 2024
    december2024, _ := client.EventRepo().Query().
        Filter(
            where.Event.StartTime.Year().Equal(2024),
            where.Event.StartTime.Month().Equal(12),
        ).
        All(ctx)

    // Events ending within 2 hours
    endingSoon, _ := client.EventRepo().Query().
        Filter(
            where.Event.EndTime.After(now),
            where.Event.EndTime.Before(now.Add(2*time.Hour)),
        ).
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
| `Before(val)` | Before time | Bool filter |
| `BeforeOrEqual(val)` | Before or equal | Bool filter |
| `After(val)` | After time | Bool filter |
| `AfterOrEqual(val)` | After or equal | Bool filter |
| `LessThan(val)` | Same as Before | Bool filter |
| `LessThanEqual(val)` | Same as BeforeOrEqual | Bool filter |
| `GreaterThan(val)` | Same as After | Bool filter |
| `GreaterThanEqual(val)` | Same as AfterOrEqual | Bool filter |
| `Add(duration)` | Add duration | Time filter |
| `Sub(duration)` | Subtract duration | Time filter |
| `Year()` | Extract year | Numeric filter |
| `Month()` | Extract month (1-12) | Numeric filter |
| `Day()` | Extract day (1-31) | Numeric filter |
| `Hour()` | Extract hour (0-23) | Numeric filter |
| `Minute()` | Extract minute (0-59) | Numeric filter |
| `Second()` | Extract second (0-59) | Numeric filter |
| `Nano()` | Extract nanoseconds | Numeric filter |
| `Micros()` | Microseconds since epoch | Numeric filter |
| `Millis()` | Milliseconds since epoch | Numeric filter |
| `Unix()` | Unix timestamp | Numeric filter |
| `Weekday()` | Day of week (0-6) | Numeric filter |
| `Week()` | Week of year (1-53) | Numeric filter |
| `YearDay()` | Day of year (1-366) | Numeric filter |
| `Group(unit)` | Group by time unit | Time filter |
| `Floor(duration)` | Floor to duration | Time filter |
| `Round(duration)` | Round to duration | Time filter |
| `Format(fmt)` | Format as string | String filter |
| `IsLeapYear()` | Check leap year | Bool filter |
| `Zero(bool)` | Check zero time | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `IsNil()` | Is null (ptr) | Bool filter |
| `IsNotNil()` | Not null (ptr) | Bool filter |
