# Optimistic Locking

Optimistic locking prevents concurrent updates from overwriting each other by tracking a version number on each record.

## Overview

When multiple processes try to update the same record simultaneously, one update might overwrite another's changes without knowing they occurred. Optimistic locking solves this by:

1. Tracking a version number on each record
2. Incrementing the version on each successful update
3. Rejecting updates that don't match the current version

## Enabling Optimistic Locking

Embed `som.OptimisticLock` in your model:

```go
type Document struct {
    som.Node
    som.OptimisticLock  // Adds version tracking

    Title   string
    Content string
}
```

This creates a hidden `__som_lock_version` field in the database that's managed automatically.

## How It Works

### Version Tracking

```go
// Create a document
doc := &model.Document{Title: "Draft"}
client.DocumentRepo().Create(ctx, doc)
fmt.Println(doc.Version())  // 1

// Update the document
doc.Title = "Final"
client.DocumentRepo().Update(ctx, doc)
fmt.Println(doc.Version())  // 2

// Each update increments the version
doc.Content = "Hello World"
client.DocumentRepo().Update(ctx, doc)
fmt.Println(doc.Version())  // 3
```

### Conflict Detection

When an update is attempted with an outdated version, SurrealDB throws an error:

```go
// Process A reads the document (version 2)
docA, _, _ := client.DocumentRepo().Read(ctx, docID)

// Process B also reads the document (version 2)
docB, _, _ := client.DocumentRepo().Read(ctx, docID)

// Process A updates successfully
docA.Title = "Updated by A"
client.DocumentRepo().Update(ctx, docA)  // Success, version becomes 3

// Process B tries to update with stale version
docB.Title = "Updated by B"
err := client.DocumentRepo().Update(ctx, docB)  // Error!
// err is som.ErrOptimisticLock
```

## Error Handling

Import the generated `som` package to access the error:

```go
import "yourproject/gen/som"

err := client.DocumentRepo().Update(ctx, staleDoc)
if errors.Is(err, som.ErrOptimisticLock) {
    // Handle the conflict
    handleConflict(ctx, staleDoc)
}
```

### Conflict Resolution Strategies

#### Strategy 1: Retry with Fresh Data

```go
func updateWithRetry(ctx context.Context, id *som.ID, updateFn func(*model.Document)) error {
    maxRetries := 3

    for i := 0; i < maxRetries; i++ {
        // Get fresh copy
        doc, exists, err := client.DocumentRepo().Read(ctx, id)
        if err != nil || !exists {
            return err
        }

        // Apply changes
        updateFn(doc)

        // Try to update
        err = client.DocumentRepo().Update(ctx, doc)
        if err == nil {
            return nil  // Success
        }

        if !errors.Is(err, som.ErrOptimisticLock) {
            return err  // Different error
        }

        // Optimistic lock failed, retry
    }

    return fmt.Errorf("update failed after %d retries", maxRetries)
}

// Usage
updateWithRetry(ctx, docID, func(doc *model.Document) {
    doc.Title = "New Title"
})
```

#### Strategy 2: Merge Changes

```go
func mergeUpdate(ctx context.Context, staleDoc *model.Document) error {
    // Get current version
    current, _, err := client.DocumentRepo().Read(ctx, staleDoc.ID())
    if err != nil {
        return err
    }

    // Merge changes (application-specific logic)
    current.Title = staleDoc.Title  // Take our title
    // Keep current.Content from database

    // Update with merged data
    return client.DocumentRepo().Update(ctx, current)
}
```

#### Strategy 3: User Resolution

```go
func handleConflict(ctx context.Context, staleDoc *model.Document) error {
    current, _, _ := client.DocumentRepo().Read(ctx, staleDoc.ID())

    // Present both versions to user
    fmt.Println("Your version:", staleDoc.Title)
    fmt.Println("Current version:", current.Title)
    fmt.Println("Choose which to keep: (y)ours or (c)urrent?")

    // ... get user input and apply chosen version
}
```

## Version Field

Access the current version number:

```go
doc, _, _ := client.DocumentRepo().Read(ctx, docID)
version := doc.Version()
fmt.Printf("Document version: %d\n", version)
```

## Combining with Timestamps

You can use optimistic locking together with automatic timestamps:

```go
type Document struct {
    som.Node
    som.Timestamps      // CreatedAt, UpdatedAt
    som.OptimisticLock  // Version tracking

    Title   string
    Content string
}
```

## Generated Schema

The optimistic lock field is defined in SurrealDB as:

```surql
DEFINE FIELD __som_lock_version ON TABLE document TYPE int
    VALUE {
        IF $value != NONE AND $before != NONE AND $value != $before {
            THROW "optimistic_lock_failed"
        };
        RETURN IF $before THEN $before + 1 ELSE 1 END;
    };
```

This schema:
- Starts at version 1 for new records
- Increments on each successful update
- Throws an error if the submitted version doesn't match the current version

## Complete Example

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "yourproject/gen/som"
    "yourproject/model"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create document
    doc := &model.Document{
        Title:   "My Document",
        Content: "Initial content",
    }
    client.DocumentRepo().Create(ctx, doc)
    fmt.Printf("Created with version %d\n", doc.Version())

    // Simulate concurrent access
    // Process A reads
    docA, _, _ := client.DocumentRepo().Read(ctx, doc.ID())

    // Process B reads same document
    docB, _, _ := client.DocumentRepo().Read(ctx, doc.ID())

    // Process A updates successfully
    docA.Content = "Updated by A"
    err := client.DocumentRepo().Update(ctx, docA)
    if err != nil {
        panic(err)
    }
    fmt.Printf("A updated successfully, version now %d\n", docA.Version())

    // Process B tries to update with stale version
    docB.Content = "Updated by B"
    err = client.DocumentRepo().Update(ctx, docB)
    if errors.Is(err, som.ErrOptimisticLock) {
        fmt.Println("Conflict detected! B's update was rejected.")

        // Resolve by getting fresh data and retrying
        fresh, _, _ := client.DocumentRepo().Read(ctx, doc.ID())
        fresh.Content = "Updated by B (retry)"
        err = client.DocumentRepo().Update(ctx, fresh)
        if err == nil {
            fmt.Printf("B retry succeeded, version now %d\n", fresh.Version())
        }
    }
}
```

Output:
```
Created with version 1
A updated successfully, version now 2
Conflict detected! B's update was rejected.
B retry succeeded, version now 3
```

## When to Use Optimistic Locking

**Good use cases:**
- Collaborative editing (documents, wikis)
- Inventory management (stock levels)
- Financial transactions (account balances)
- Any data where concurrent updates are possible

**Consider alternatives when:**
- Data is rarely updated concurrently
- Updates are append-only (use immutable records)
- You need pessimistic locking (SurrealDB transactions)

## Performance Considerations

- Optimistic locking adds minimal overhead (one integer field)
- Conflicts are detected at update time, not read time
- No database locks are held during processing
- Suitable for high-concurrency, low-conflict scenarios
