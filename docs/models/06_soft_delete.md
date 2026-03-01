# Soft Delete

Soft delete marks records as deleted without removing them from the database, allowing for data recovery and audit trails.

## Overview

When you delete a record normally, it's permanently removed from the database. Soft delete provides an alternative approach:

1. Records are marked with a `deleted_at` timestamp instead of being removed
2. Queries automatically exclude soft-deleted records
3. Records can be restored or permanently erased later

## Enabling Soft Delete

Embed `som.SoftDelete` in your model:

```go
type User struct {
    som.Node
    som.SoftDelete  // Adds soft delete functionality

    Name  string
    Email string
}
```

This creates a `deleted_at` field in the database that's managed automatically.

## How It Works

### Automatic Query Filtering

By default, soft-deleted records are excluded from all queries:

```go
// Create some users
client.UserRepo().Create(ctx, &model.User{Name: "Alice"})
client.UserRepo().Create(ctx, &model.User{Name: "Bob"})

// Delete Alice (soft delete)
alice, _, _ := client.UserRepo().Query().
    Where(filter.User.Name.Equal("Alice")).
    First(ctx)
client.UserRepo().Delete(ctx, alice)

// Query returns only Bob
users, _ := client.UserRepo().Query().All(ctx)
fmt.Println(len(users))  // 1 (only Bob)
```

### Including Deleted Records

Use `WithDeleted()` to include soft-deleted records in your query:

```go
// Get all users including deleted ones
allUsers, _ := client.UserRepo().Query().WithDeleted().All(ctx)
fmt.Println(len(allUsers))  // 2 (Alice and Bob)
```

### Checking Deletion Status

Every model with `som.SoftDelete` has two methods:

```go
user, _, _ := client.UserRepo().Read(ctx, userID)

// Check if deleted
if user.SoftDelete.IsDeleted() {
    fmt.Println("User was deleted at:", user.SoftDelete.DeletedAt())
}

// DeletedAt() returns time.Time
// Returns zero value if not deleted
deletedAt := user.SoftDelete.DeletedAt()
if deletedAt.IsZero() {
    fmt.Println("User is active")
}
```

## Repository Methods

Models with soft delete enabled have three deletion-related methods:

### Delete (Soft Delete)

Marks the record as deleted by setting the `deleted_at` timestamp:

```go
err := client.UserRepo().Delete(ctx, user)
if err != nil {
    // Handle error
}

// The in-memory object is automatically refreshed
fmt.Println(user.SoftDelete.IsDeleted())    // true
fmt.Println(user.SoftDelete.DeletedAt())    // 2024-01-15 10:30:00 ...
```

Deleting an already-deleted record returns an error:

```go
err := client.UserRepo().Delete(ctx, alreadyDeletedUser)
// err: "record is already deleted"
```

### Restore

Un-deletes a soft-deleted record by clearing the `deleted_at` timestamp:

```go
err := client.UserRepo().Restore(ctx, deletedUser)
if err != nil {
    // Handle error
}

// The in-memory object is automatically refreshed
fmt.Println(deletedUser.SoftDelete.IsDeleted())  // false
fmt.Println(deletedUser.SoftDelete.DeletedAt())  // 0001-01-01 00:00:00 (zero)
```

Restoring a non-deleted record returns an error:

```go
err := client.UserRepo().Restore(ctx, activeUser)
// err: "record is not deleted, cannot restore"
```

### Erase (Hard Delete)

Permanently removes the record from the database. This cannot be undone:

```go
// Permanently delete the record
err := client.UserRepo().Erase(ctx, user)
if err != nil {
    // Handle error
}

// Record is completely gone
user, exists, _ := client.UserRepo().Query().WithDeleted().
    Where(filter.User.ID.Equal(userID)).
    First(ctx)
// exists == false
```

## Querying Soft-Deleted Records

### Default Behavior

All queries automatically filter out soft-deleted records:

```go
// Returns only non-deleted users
users, _ := client.UserRepo().Query().All(ctx)

// Count excludes deleted records
count, _ := client.UserRepo().Query().Count(ctx)
```

### Include Deleted Records

Use `WithDeleted()` to query all records:

```go
// All records including soft-deleted
allUsers, _ := client.UserRepo().Query().WithDeleted().All(ctx)

// Count including deleted
totalCount, _ := client.UserRepo().Query().WithDeleted().Count(ctx)
```

### Find Only Deleted Records

Filter for records where `DeletedAt` is set:

```go
// Get only soft-deleted users
deletedUsers, _ := client.UserRepo().Query().
    WithDeleted().
    Where(filter.User.DeletedAt.Nil(false)).
    All(ctx)
```

## Fetching Related Records

### Important Behavior

Soft-delete filtering does **NOT** apply to fetched relations. When you use `Fetch()` to load related records, all records are returned regardless of their soft-delete status:

```go
type Post struct {
    som.Node
    som.SoftDelete
    Title  string
    Author *User  // User also has SoftDelete
}

// Soft delete the author
client.UserRepo().Delete(ctx, author)

// Fetch post with author - the deleted author IS returned
posts, _ := client.PostRepo().Query().
    Fetch(with.Post.Author()).
    All(ctx)

// The author is included even though soft-deleted
post := posts[0]
fmt.Println(post.Author.Name)                   // "Alice"
fmt.Println(post.Author.SoftDelete.IsDeleted()) // true
```

### Manual Filtering

If you need to filter soft-deleted records from fetched relations, do it in your application code:

```go
posts, _ := client.PostRepo().Query().
    Fetch(with.Post.Author()).
    All(ctx)

for _, post := range posts {
    if post.Author != nil && !post.Author.SoftDelete.IsDeleted() {
        // Process only posts with active authors
        fmt.Println(post.Title, "by", post.Author.Name)
    }
}
```

For slice relations:

```go
type BlogPost struct {
    som.Node
    som.SoftDelete
    Title   string
    Authors []*User
}

posts, _ := client.BlogPostRepo().Query().
    Fetch(with.BlogPost.Authors()).
    All(ctx)

for _, post := range posts {
    // Filter out soft-deleted authors
    var activeAuthors []*model.User
    for _, author := range post.Authors {
        if !author.SoftDelete.IsDeleted() {
            activeAuthors = append(activeAuthors, author)
        }
    }
    fmt.Printf("%s has %d active authors\n", post.Title, len(activeAuthors))
}
```

## Combining with Other Features

Soft delete works seamlessly with Timestamps and OptimisticLock:

```go
type Document struct {
    som.Node
    som.Timestamps      // CreatedAt, UpdatedAt
    som.OptimisticLock  // Version tracking
    som.SoftDelete      // Soft delete

    Title   string
    Content string
}
```

All features work together:

```go
// Create
doc := &model.Document{Title: "Draft"}
client.DocumentRepo().Create(ctx, doc)
fmt.Println(doc.Version())                // 1
fmt.Println(doc.Timestamps.CreatedAt())   // 2024-01-15 10:30:00
fmt.Println(doc.SoftDelete.IsDeleted())   // false

// Soft delete (increments version)
client.DocumentRepo().Delete(ctx, doc)
fmt.Println(doc.Version())                // 2
fmt.Println(doc.SoftDelete.IsDeleted())   // true

// Restore (increments version)
client.DocumentRepo().Restore(ctx, doc)
fmt.Println(doc.Version())                // 3
fmt.Println(doc.SoftDelete.IsDeleted())   // false

// Update (increments version)
doc.Title = "Final"
client.DocumentRepo().Update(ctx, doc)
fmt.Println(doc.Version())                // 4
```

## Generated Schema

The soft delete field is defined in SurrealDB as:

```surql
DEFINE FIELD deleted_at ON TABLE user TYPE option<datetime> DEFAULT NONE;
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "yourproject/gen/som"
    "yourproject/gen/som/filter"
    "yourproject/model"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create a user
    user := &model.User{
        Name:  "Alice",
        Email: "alice@example.com",
    }
    client.UserRepo().Create(ctx, user)
    fmt.Println("Created user:", user.Name)

    // Verify not deleted
    fmt.Println("Is deleted:", user.SoftDelete.IsDeleted())  // false

    // Soft delete the user
    err := client.UserRepo().Delete(ctx, user)
    if err != nil {
        panic(err)
    }
    fmt.Println("Soft deleted at:", user.SoftDelete.DeletedAt())

    // Normal query doesn't find the user
    users, _ := client.UserRepo().Query().All(ctx)
    fmt.Println("Active users:", len(users))  // 0

    // WithDeleted() finds the user
    allUsers, _ := client.UserRepo().Query().WithDeleted().All(ctx)
    fmt.Println("All users:", len(allUsers))  // 1

    // Restore the user
    err = client.UserRepo().Restore(ctx, user)
    if err != nil {
        panic(err)
    }
    fmt.Println("Restored, is deleted:", user.SoftDelete.IsDeleted())  // false

    // User appears in normal queries again
    users, _ = client.UserRepo().Query().All(ctx)
    fmt.Println("Active users:", len(users))  // 1

    // Soft delete and then permanently remove
    client.UserRepo().Delete(ctx, user)
    err = client.UserRepo().Erase(ctx, user)
    if err != nil {
        panic(err)
    }
    fmt.Println("User permanently erased")

    // User is completely gone
    allUsers, _ = client.UserRepo().Query().WithDeleted().All(ctx)
    fmt.Println("All users:", len(allUsers))  // 0
}
```

Output:
```
Created user: Alice
Is deleted: false
Soft deleted at: 2024-01-15 10:30:00.123456789 +0000 UTC
Active users: 0
All users: 1
Restored, is deleted: false
Active users: 1
User permanently erased
All users: 0
```

## When to Use Soft Delete

**Good use cases:**
- Audit requirements (keeping history of all records)
- User data recovery (accidental deletion protection)
- Maintaining referential integrity (related records still reference the deleted record)
- Legal/compliance requirements for data retention
- Implementing "trash" or "archive" functionality

**Consider alternatives when:**
- Storage constraints are a concern (soft-deleted records consume space)
- GDPR or other regulations require actual data deletion
- High-volume temporary data that should be truly purged
- Performance is critical and you have many soft-deleted records

## Best Practices

1. **Regularly purge old soft-deleted records** if storage is a concern:
   ```go
   // Find records deleted more than 30 days ago
   cutoff := time.Now().AddDate(0, 0, -30)
   oldDeleted, _ := client.UserRepo().Query().
       WithDeleted().
       Where(filter.User.DeletedAt.LessThan(cutoff)).
       All(ctx)

   // Permanently remove them
   for _, user := range oldDeleted {
       client.UserRepo().Erase(ctx, user)
   }
   ```

2. **Handle fetched relations explicitly** - remember that soft-delete filtering doesn't apply to `Fetch()` results

3. **Combine with Timestamps** for complete audit trails showing when records were created, updated, and deleted
