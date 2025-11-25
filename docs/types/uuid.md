# UUID Type

The UUID type handles universally unique identifiers using `github.com/google/uuid` with binary CBOR encoding.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `uuid.UUID` / `*uuid.UUID` |
| Database Schema | `uuid` / `option<uuid>` |
| CBOR Encoding | Tag 37 with 16-byte binary |
| Sortable | Yes |

## CBOR Encoding

UUID values are encoded with CBOR Tag 37 as binary data:

```
Tag 37: <16 bytes binary UUID>
```

This provides:
- Efficient binary storage
- Proper round-tripping with SurrealDB
- Standard UUID representation

## Definition

```go
import "github.com/google/uuid"

type Document struct {
    som.Node

    ExternalID  uuid.UUID   // Required
    TrackingID  uuid.UUID   // Required
    ParentID    *uuid.UUID  // Optional
}
```

## Schema

Generated SurrealDB schema:

```surql
DEFINE FIELD external_id ON document TYPE uuid;
DEFINE FIELD tracking_id ON document TYPE uuid;
DEFINE FIELD parent_id ON document TYPE option<uuid>;
```

## Creating UUIDs

```go
import "github.com/google/uuid"

// New random UUID
doc := &model.Document{
    ExternalID: uuid.New(),
    TrackingID: uuid.New(),
}

// Parse from string
id, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
doc.ExternalID = id

// Nil UUID
doc.ParentID = nil  // Optional field
```

## Filter Operations

### Equality Operations

```go
// Exact match
targetID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
where.Document.ExternalID.Equal(targetID)

// Not equal
where.Document.ExternalID.NotEqual(excludeID)
```

### Set Membership

```go
// Value in set
where.Document.ExternalID.In(id1, id2, id3)

// Value not in set
where.Document.TrackingID.NotIn(blacklistedIDs...)
```

### Comparison Operations

UUIDs can be compared lexicographically:

```go
// Less than
where.Document.ExternalID.LessThan(referenceID)

// Less than or equal
where.Document.ExternalID.LessThanEqual(referenceID)

// Greater than
where.Document.ExternalID.GreaterThan(referenceID)

// Greater than or equal
where.Document.ExternalID.GreaterThanEqual(referenceID)
```

### Nil Operations (Pointer Types Only)

```go
// Check if nil
where.Document.ParentID.IsNil()

// Check if not nil
where.Document.ParentID.IsNotNil()
```

### Zero Value Check

```go
// Is nil UUID (00000000-0000-0000-0000-000000000000)
where.Document.ExternalID.Zero(true)

// Is not nil UUID
where.Document.ExternalID.Zero(false)
```

## Sorting

```go
// Ascending (lexicographic)
query.Order(by.Document.ExternalID.Asc())

// Descending
query.Order(by.Document.ExternalID.Desc())
```

## Common Patterns

### Find by External ID

```go
externalID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

doc, exists, _ := client.DocumentRepo().Query().
    Filter(where.Document.ExternalID.Equal(externalID)).
    First(ctx)
```

### Find Children by Parent

```go
parentID := uuid.MustParse("...")

children, _ := client.DocumentRepo().Query().
    Filter(where.Document.ParentID.Equal(parentID)).
    All(ctx)
```

### Find Root Documents (No Parent)

```go
roots, _ := client.DocumentRepo().Query().
    Filter(where.Document.ParentID.IsNil()).
    All(ctx)
```

### Bulk Lookup

```go
targetIDs := []uuid.UUID{id1, id2, id3}

docs, _ := client.DocumentRepo().Query().
    Filter(where.Document.ExternalID.In(targetIDs...)).
    All(ctx)
```

## Complete Example

```go
package main

import (
    "context"
    "github.com/google/uuid"
    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/where"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create document with UUIDs
    doc := &model.Document{
        ExternalID: uuid.New(),
        TrackingID: uuid.New(),
    }
    client.DocumentRepo().Create(ctx, doc)

    // Find by external ID
    targetID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
    found, exists, _ := client.DocumentRepo().Query().
        Filter(where.Document.ExternalID.Equal(targetID)).
        First(ctx)

    // Find documents with parent
    withParent, _ := client.DocumentRepo().Query().
        Filter(where.Document.ParentID.IsNotNil()).
        All(ctx)

    // Find root documents
    roots, _ := client.DocumentRepo().Query().
        Filter(where.Document.ParentID.IsNil()).
        All(ctx)

    // Bulk lookup
    ids := []uuid.UUID{id1, id2, id3}
    batch, _ := client.DocumentRepo().Query().
        Filter(where.Document.ExternalID.In(ids...)).
        All(ctx)

    // Exclude specific documents
    excluded := []uuid.UUID{badID1, badID2}
    filtered, _ := client.DocumentRepo().Query().
        Filter(where.Document.ExternalID.NotIn(excluded...)).
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
| `LessThan(val)` | Lexicographic < | Bool filter |
| `LessThanEqual(val)` | Lexicographic <= | Bool filter |
| `GreaterThan(val)` | Lexicographic > | Bool filter |
| `GreaterThanEqual(val)` | Lexicographic >= | Bool filter |
| `Zero(bool)` | Check nil UUID | Bool filter |
| `Truth()` | To boolean | Bool filter |
| `IsNil()` | Is null (ptr) | Bool filter |
| `IsNotNil()` | Not null (ptr) | Bool filter |
