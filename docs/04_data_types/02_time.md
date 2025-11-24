# Time Types

SOM provides support for Go's time-related types.

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

### Querying Time Fields

```go
// Find events starting after a specific time
events, err := client.EventRepo().Query().
    Filter(where.Event.StartTime.GreaterThan(time.Now())).
    All(ctx)

// Find events within a range
events, err := client.EventRepo().Query().
    Filter(
        where.Event.StartTime.GreaterThanOrEqual(startDate),
        where.Event.EndTime.LessThanOrEqual(endDate),
    ).
    All(ctx)
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

### Working with Durations

```go
task := &model.Task{
    Name:     "Build",
    Duration: 30 * time.Minute,
    Timeout:  2 * time.Hour,
}
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

## Timezone Handling

SOM stores times in UTC. Convert to local time when displaying:

```go
event, _ := client.EventRepo().Read(ctx, id)
localTime := event.StartTime.Local()
```
