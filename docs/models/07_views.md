# Views

A view is a read-only, pre-computed table backed by a SurrealDB
`DEFINE TABLE ... AS SELECT` statement. Its rows are maintained by the database from a source
table — you never write to a view.

## Defining a View Model

Embed `som.View` instead of `som.Node[T]`:

```go
package model

import "yourproject/gen/som"

// EventSummary aggregates EventLog rows, grouped by Category.
type EventSummary struct {
    som.View

    Category string
    Total    int
    AvgValue float64
}
```

The fields of the view model are the projected columns. A view model has an `ID() string`
accessor but no timestamps, soft delete, optimistic locking or expiry — those features are not
supported on views.

## Declaring the Projection

The `SELECT` behind the view is declared in `som.define.go` (see
[Schema Definitions](../code_generation/03_definitions.md)) using the generated `define`
package:

```go
//go:build som

package myapp

import (
    "yourproject/gen/som/define"
    "yourproject/gen/som/define/aggregate"
    "yourproject/gen/som/filter"
    "yourproject/model"
)

func Definitions() define.Definitions {
    return define.Definitions{
        Views: []define.ViewDefinition{eventSummary},
    }
}

var eventSummary = define.View[model.EventSummary, model.EventLog]().
    Project(
        define.As(filter.EventSummary.Category, filter.EventLog.Category),
        define.As(filter.EventSummary.Total, aggregate.CountAll(filter.EventLog)),
        define.As(filter.EventSummary.AvgValue, aggregate.Mean(filter.EventLog.Value)),
    ).
    GroupBy(filter.EventSummary.Category)
```

`define.View[V, T]()` takes the view model `V` and the source model `T`. `define.As(column, expr)`
maps a source expression onto a view column; both sides are checked at compile time — the column
must belong to `V`, the expression to `T`, and their value types must match.

### Builder Methods

| Method | Description |
|--------|-------------|
| `Project(projections...)` | The columns the view computes |
| `Where(filters...)` | Restrict which source rows enter the view |
| `GroupBy(columns...)` | Group source rows before aggregating, by view columns |

`Where` renders its values as literals, since a table definition cannot carry query parameters.

`GroupBy` takes columns of the view, not fields of the source. SurrealDB groups by idioms only, so
a computed projection can only be grouped through its alias. Grouping by a view column always
renders that alias:

```go
var eventsByMonth = define.View[model.EventsByMonth, model.EventLog]().
    Project(
        define.As(filter.EventsByMonth.Category, filter.EventLog.Category),
        define.As(filter.EventsByMonth.Month, filter.EventLog.OccurredAt.Month()),
        define.As(filter.EventsByMonth.Total, aggregate.CountAll(filter.EventLog)),
    ).
    GroupBy(filter.EventsByMonth.Category, filter.EventsByMonth.Month)

// DEFINE TABLE events_by_month TYPE NORMAL AS
//     SELECT category AS category, time::month(occurred_at) AS month, count() AS total
//     FROM event_log
//     GROUP BY category, month;
```

### Aggregates

Aggregate functions live in the generated `define/aggregate` package. Only the incrementally
maintained set supported by SurrealDB views is exposed:

`Count`, `CountAll`, `Sum`, `Mean`, `Min`, `Max`, `Variance`, `StdDev`.

`Count(field)` counts the non-null values of a field, `CountAll(root)` counts rows. `CountAll`
takes the model's filter root (e.g. `filter.EventLog`) only to pin the source model.

Aggregates are functions rather than methods on filter fields, which keeps them out of `WHERE`,
where they would be invalid.

The example above generates:

```surql
DEFINE TABLE event_summary TYPE NORMAL AS
    SELECT category AS category, count() AS total, math::mean(value) AS avg_value
    FROM event_log
    GROUP BY category;
```

## Querying a View

Views get a repository with a single method — `Query()` — returning the same query builder used
for nodes:

```go
rows, err := client.EventSummaryRepo().Query().
    Where(filter.EventSummary.Total.GreaterThan(10)).
    Order(by.EventSummary.Total.Desc()).
    All(ctx)
```

There is no `Create`, `Update`, `Delete` or `Insert` — the database maintains the rows.

## Caveats

- **Stale linked data.** Values projected from linked tables (via record links or graph
  traversal, e.g. `->product.name`) are frozen at the time the source row is written. They do
  **not** refresh when the linked record changes.
- **Source writes drive the view.** A view only updates when rows of the source table are
  written. Changing the view definition requires re-applying the schema.
- Views cannot be the source of another view.

## Sink → View

A view is often paired with a write-only [sink](08_sinks.md): the sink accepts high-volume
records that are discarded immediately, while the view keeps only the aggregate.
