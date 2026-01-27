# Time Types

SOM provides support for Go's time-related types with automatic CBOR serialization for SurrealDB.

## time.Time

Use `time.Time` for timestamps and dates:

```go
import "time"

type Event struct {
    som.Node

    Name      string
    StartTime time.Time
    EndTime   time.Time
}
```

### CBOR Encoding

Internally, SOM uses a custom `DateTime` wrapper that handles CBOR serialization. Time values are encoded using:

- **CBOR Tag 12** for datetime values
- **Two-element array**: `[unix_seconds, nanoseconds]`

This ensures nanosecond precision and proper round-tripping with SurrealDB.

### Querying Time Fields

```go
// Find events starting after a specific time
events, err := client.EventRepo().Query().
    Where(filter.Event.StartTime.GreaterThan(time.Now())).
    All(ctx)

// Find events within a range
events, err := client.EventRepo().Query().
    Where(
        filter.Event.StartTime.GreaterThanOrEqual(startDate),
        filter.Event.EndTime.LessThanOrEqual(endDate),
    ).
    All(ctx)
```

### Time Filter Operations

| Operation | Description |
|-----------|-------------|
| `Before(time)` | Before specified time |
| `BeforeOrEqual(time)` | Before or equal to |
| `After(time)` | After specified time |
| `AfterOrEqual(time)` | After or equal to |
| `Add(duration)` | Add duration to time |
| `Sub(duration)` | Subtract duration |
| `Floor(duration)` | Floor to duration unit |
| `Round(duration)` | Round to duration unit |
| `Format(format)` | Format as string |

```go
// Created in last 7 days
filter.User.CreatedAt.After(time.Now().AddDate(0, 0, -7))

// Compare with date math
filter.Event.StartTime.Add(2 * time.Hour).Before(deadline)
```

## time.Duration

Use `time.Duration` for time intervals:

```go
import "time"

type Task struct {
    som.Node

    Name     string
    Duration time.Duration
    Timeout  time.Duration
}
```

### CBOR Encoding

Duration values are encoded using:

- **CBOR Tag 14** for duration values
- **Two-element array**: `[seconds, nanoseconds]`

### Working with Durations

```go
task := &model.Task{
    Name:     "Build",
    Duration: 30 * time.Minute,
    Timeout:  2 * time.Hour,
}
```

### Duration Filter Operations

| Operation | Description |
|-----------|-------------|
| `Before(duration)` | Shorter than |
| `After(duration)` | Longer than |
| `Add(duration)` | Add durations |
| `Sub(duration)` | Subtract durations |

```go
// Sessions longer than 1 hour
filter.Session.Duration.After(time.Hour)

// Timeouts under 5 minutes
filter.Task.Timeout.Before(5 * time.Minute)
```

## Optional Time Fields

Use pointers for optional time values:

```go
type User struct {
    som.Node

    CreatedAt  time.Time   // Always set
    DeletedAt  *time.Time  // Optional, nil if not deleted
    LastLogin  *time.Time  // Optional
}
```

Query optional time fields:

```go
// Find soft-deleted users
filter.User.DeletedAt.IsNotNil()

// Find users who have never logged in
filter.User.LastLogin.IsNil()
```

## Automatic Timestamps

Use `som.Timestamps` for automatic tracking:

```go
type User struct {
    som.Node
    som.Timestamps  // Adds CreatedAt and UpdatedAt

    Name string
}
```

- `CreatedAt` - Set automatically on create (readonly)
- `UpdatedAt` - Updated automatically on every save

Both fields are readonly in your application code and managed by SOM.

## Timezone Handling

SOM stores times in UTC. Convert to local time when displaying:

```go
event, _, _ := client.EventRepo().Read(ctx, id)
localTime := event.StartTime.Local()
```

Best practice: always work with UTC internally and convert for display only.

## Examples

### Event Scheduling

```go
// Find upcoming events in the next week
nextWeek := time.Now().AddDate(0, 0, 7)
events, err := client.EventRepo().Query().
    Where(
        filter.Event.StartTime.After(time.Now()),
        filter.Event.StartTime.Before(nextWeek),
    ).
    Order(by.Event.StartTime.Asc()).
    All(ctx)
```

### Session Expiry

```go
// Find expired sessions (created more than 24 hours ago)
cutoff := time.Now().Add(-24 * time.Hour)
expired, err := client.SessionRepo().Query().
    Where(filter.Session.CreatedAt.Before(cutoff)).
    All(ctx)
```
