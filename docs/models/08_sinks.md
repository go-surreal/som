# Sinks

A sink is a write-only ingestion table backed by a SurrealDB `DEFINE TABLE ... DROP` statement.
Records written to a sink are accepted — firing any [views](07_views.md) and events that select
from the table — but are **not persisted**: the row is discarded immediately after the write.

Typical use: feeding an aggregating view from a high-volume event or log stream where only the
aggregate matters, not the raw records.

## Defining a Sink Model

Embed `som.Sink`:

```go
package model

import "yourproject/gen/som"

// EventLog is a write-only ingestion sink.
type EventLog struct {
    som.Sink

    Category string
    Value    float64
}
```

A sink has **no ID** and no `ID()` accessor. Timestamps, soft delete, optimistic locking, expiry
and changefeeds are not supported on sinks.

Generated DDL:

```surql
DEFINE TABLE event_log DROP SCHEMAFULL TYPE NORMAL PERMISSIONS FULL;
```

## Writing to a Sink

The generated repository exposes only `Create` and `Insert`:

```go
err := client.EventLogRepo().Create(ctx, &model.EventLog{
    Category: "checkout",
    Value:    19.99,
})

err = client.EventLogRepo().Insert(ctx, []*model.EventLog{
    {Category: "checkout", Value: 19.99},
    {Category: "signup", Value: 0},
})
```

Nothing is returned beyond the error — there is no record to read back. Reading, querying,
updating and deleting are not available.

## Sink → View Pattern

Pair a sink with a view to keep aggregates without storing raw events:

```go
// Write-only ingestion
type EventLog struct {
    som.Sink

    Category string
    Value    float64
}

// Read-only aggregate over the discarded rows
type EventSummary struct {
    som.View

    Category string
    Total    int
    AvgValue float64
}
```

```go
var eventSummary = define.View[model.EventSummary, model.EventLog]().
    Project(
        define.As(filter.EventSummary.Category, filter.EventLog.Category),
        define.As(filter.EventSummary.Total, aggregate.Count(filter.EventLog.Category)),
        define.As(filter.EventSummary.AvgValue, aggregate.Mean(filter.EventLog.Value)),
    ).
    GroupBy(filter.EventLog.Category)
```

Writes go to the sink, reads come from the view:

```go
_ = client.EventLogRepo().Create(ctx, &model.EventLog{Category: "checkout", Value: 19.99})

rows, err := client.EventSummaryRepo().Query().All(ctx)
```

## When to Use

- High-volume events where only aggregates are kept (metrics, counters, telemetry)
- Feeding several views from one ingestion stream
- Avoiding a storage-heavy raw table plus a separate cleanup job

If you need the raw records later, use a normal node — optionally with
[expiry](09_expiry.md) — instead.
