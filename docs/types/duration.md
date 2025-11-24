# Duration Type

The duration type handles time intervals using Go's `time.Duration` with special CBOR encoding.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `time.Duration` / `*time.Duration` |
| Database Schema | `duration` / `option<duration>` |
| CBOR Encoding | Tag 14 with `[seconds, nanoseconds]` |
| Sortable | Yes |

## CBOR Encoding

Duration values are encoded with CBOR Tag 14 as a two-element array:

```
Tag 14: [total_seconds (int64), remaining_nanoseconds (int64)]
```

This provides:
- Full nanosecond precision
- Proper round-tripping with SurrealDB
- Handling of empty arrays for optional fields

## Definition

```go
type Task struct {
    som.Node

    Duration    time.Duration   // Required
    Timeout     time.Duration   // Required
    GracePeriod *time.Duration  // Optional
}
```

## Schema

Generated SurrealDB schema:

```surql
DEFINE FIELD duration ON task TYPE duration;
DEFINE FIELD timeout ON task TYPE duration;
DEFINE FIELD grace_period ON task TYPE option<duration>;
```

## Filter Operations

### Equality Operations

```go
// Exact match
where.Task.Duration.Equal(30 * time.Minute)

// Not equal
where.Task.Timeout.NotEqual(0)
```

### Set Membership

```go
// Value in set
where.Task.Duration.In(
    15 * time.Minute,
    30 * time.Minute,
    1 * time.Hour,
)

// Value not in set
where.Task.Duration.NotIn(0, -1)
```

### Comparison Operations

```go
// Less than
where.Task.Duration.LessThan(1 * time.Hour)

// Less than or equal
where.Task.Duration.LessThanEqual(30 * time.Minute)

// Greater than
where.Task.Timeout.GreaterThan(5 * time.Second)

// Greater than or equal
where.Task.Timeout.GreaterThanEqual(1 * time.Minute)
```

### Component Extraction

Extract duration components as numeric filters:

```go
// Days
where.Task.Duration.Days().GreaterThan(0)

// Hours
where.Task.Duration.Hours().LessThan(24)

// Minutes
where.Task.Duration.Mins().Equal(30)

// Seconds
where.Task.Duration.Secs().GreaterThan(0)

// Milliseconds
where.Task.Duration.Millis().LessThan(1000)

// Microseconds
where.Task.Duration.Micros().GreaterThan(0)

// Nanoseconds
where.Task.Duration.Nanos().Equal(0)

// Weeks
where.Task.Duration.Weeks().Equal(1)

// Years (approximate, 365 days)
where.Task.Duration.Years().LessThan(1)
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
where.Task.GracePeriod.IsNil()

// Check if not nil
where.Task.GracePeriod.IsNotNil()
```

### Zero Value Check

```go
// Is zero duration
where.Task.Duration.Zero(true)

// Is not zero duration
where.Task.Duration.Zero(false)
```

## Sorting

```go
// Ascending (shortest first)
query.Order(by.Task.Duration.Asc())

// Descending (longest first)
query.Order(by.Task.Duration.Desc())

// Multiple sort fields
query.Order(
    by.Task.Duration.Asc(),
    by.Task.Name.Asc(),
)
```

## Method Chaining

Duration filters support component extraction:

```go
// Tasks taking more than 2 hours
where.Task.Duration.Hours().GreaterThan(2)

// Timeouts under 100ms
where.Task.Timeout.Millis().LessThan(100)

// Tasks lasting multiple days
where.Task.Duration.Days().GreaterThanEqual(1)
```

## Common Patterns

### Duration Ranges

```go
// Tasks between 30 minutes and 2 hours
tasks, _ := client.TaskRepo().Query().
    Filter(
        where.Task.Duration.GreaterThanEqual(30 * time.Minute),
        where.Task.Duration.LessThan(2 * time.Hour),
    ).
    All(ctx)
```

### Short vs Long Tasks

```go
// Quick tasks (under 5 minutes)
quickTasks, _ := client.TaskRepo().Query().
    Filter(where.Task.Duration.LessThan(5 * time.Minute)).
    All(ctx)

// Long-running tasks (over 1 hour)
longTasks, _ := client.TaskRepo().Query().
    Filter(where.Task.Duration.GreaterThan(1 * time.Hour)).
    All(ctx)
```

### Timeout Validation

```go
// Tasks with reasonable timeouts
validTimeouts, _ := client.TaskRepo().Query().
    Filter(
        where.Task.Timeout.GreaterThan(0),
        where.Task.Timeout.LessThanEqual(30 * time.Minute),
    ).
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

    // Find quick tasks (under 5 minutes)
    quickTasks, _ := client.TaskRepo().Query().
        Filter(where.Task.Duration.LessThan(5 * time.Minute)).
        Order(by.Task.Duration.Asc()).
        All(ctx)

    // Find tasks with hour-long durations
    hourlyTasks, _ := client.TaskRepo().Query().
        Filter(where.Task.Duration.Hours().Equal(1)).
        All(ctx)

    // Tasks with grace period configured
    withGrace, _ := client.TaskRepo().Query().
        Filter(where.Task.GracePeriod.IsNotNil()).
        All(ctx)

    // Tasks with multi-day durations
    multiDay, _ := client.TaskRepo().Query().
        Filter(where.Task.Duration.Days().GreaterThanEqual(1)).
        All(ctx)

    // Tasks sorted by duration
    sorted, _ := client.TaskRepo().Query().
        Order(by.Task.Duration.Desc()).
        Limit(10).
        All(ctx)

    // Tasks with sub-second timeout
    fastTimeout, _ := client.TaskRepo().Query().
        Filter(where.Task.Timeout.Millis().LessThan(1000)).
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
| `LessThan(val)` | Strictly less | Bool filter |
| `LessThanEqual(val)` | Less or equal | Bool filter |
| `GreaterThan(val)` | Strictly greater | Bool filter |
| `GreaterThanEqual(val)` | Greater or equal | Bool filter |
| `Days()` | Extract days | Numeric filter |
| `Hours()` | Extract hours | Numeric filter |
| `Mins()` | Extract minutes | Numeric filter |
| `Secs()` | Extract seconds | Numeric filter |
| `Millis()` | Extract milliseconds | Numeric filter |
| `Micros()` | Extract microseconds | Numeric filter |
| `Nanos()` | Extract nanoseconds | Numeric filter |
| `Weeks()` | Extract weeks | Numeric filter |
| `Years()` | Extract years | Numeric filter |
| `Zero(bool)` | Check zero | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `IsNil()` | Is null (ptr) | Bool filter |
| `IsNotNil()` | Not null (ptr) | Bool filter |
