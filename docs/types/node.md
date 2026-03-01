# Node Type

The node type represents references (links) to other database records, enabling relationships between entities.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `*OtherNode` (pointer to another Node type) |
| Database Schema | `option<record<table_name>>` |
| CBOR Encoding | Direct (as record ID) |
| Sortable | No (use nested field sorting) |

## Definition

Node references are always pointers to other Node types:

```go
type Post struct {
    som.Node

    Title   string
    Content string
    Author  *User     // Reference to User
    Editor  *User     // Another reference
}

type Comment struct {
    som.Node

    Text   string
    Author *User  // Reference to User
    Post   *Post  // Reference to Post
}
```

## Schema

Generated SurrealDB schema:

```surql
DEFINE FIELD author ON post TYPE option<record<user>>;
DEFINE FIELD editor ON post TYPE option<record<user>>;
DEFINE FIELD author ON comment TYPE option<record<user>>;
DEFINE FIELD post ON comment TYPE option<record<post>>;
```

## Creating Node References

```go
// Create user first
user := &model.User{Name: "Alice"}
client.UserRepo().Create(ctx, user)

// Create post with reference
post := &model.Post{
    Title:  "Hello World",
    Author: user,  // Direct reference
}
client.PostRepo().Create(ctx, post)
```

## Filter Operations

### ID-Based Filtering

```go
// Filter by referenced record's ID
filter.Post.Author.ID().Equal(userID)

// Not equal
filter.Post.Author.ID().NotEqual(excludeUserID)

// In set
filter.Post.Author.ID().In(authorID1, authorID2)

// Not in set
filter.Post.Author.ID().NotIn(bannedIDs...)
```

### Nested Field Filtering

Access fields of the referenced record:

```go
// Filter by author's name
filter.Post.Author.Name.Equal("Alice")

// Filter by author's email
filter.Post.Author.Email.EndsWith("@company.com")

// Filter by author's status
filter.Post.Author.IsActive.True()

// Chain nested fields
filter.Post.Author.Email.Lowercase().Contains("@admin")
```

### Deep Nesting

For deeply nested references:

```go
// Comment -> Post -> Author
filter.Comment.Post.Author.Name.Equal("Alice")

// Filter by parent post's author's status
filter.Comment.Post.Author.IsActive.True()
```

### Nil Checks

```go
// Posts with author
filter.Post.Author.IsNotNil()

// Posts without editor
filter.Post.Editor.IsNil()
```

## Sorting

Node fields themselves are not sortable, but nested fields are:

```go
// Sort by author's name
query.Order(by.Post.Author.Name.Asc())

// Sort by author's creation date
query.Order(by.Post.Author.CreatedAt.Desc())

// Multiple sort fields
query.Order(
    by.Post.Author.Name.Asc(),
    by.Post.Title.Asc(),
)
```

## Fetching References

Use `Fetch` to eager-load referenced records:

```go
import "yourproject/gen/som/with"

// Fetch author with posts
posts, _ := client.PostRepo().Query().
    Fetch(with.Post.Author...).
    All(ctx)

for _, post := range posts {
    // Author is already loaded
    fmt.Println(post.Author.Name)
}
```

### Multiple Fetches

```go
posts, _ := client.PostRepo().Query().
    Fetch(with.Post.Author...).
    Fetch(with.Post.Editor...).
    All(ctx)
```

### Nested Fetches

```go
comments, _ := client.CommentRepo().Query().
    Fetch(with.Comment.Post...).
    Fetch(with.Comment.Author...).
    All(ctx)

// Access nested: comment.Post.Title
```

## Common Patterns

### Filter by Author

```go
// Posts by specific author
authorPosts, _ := client.PostRepo().Query().
    Where(filter.Post.Author.ID().Equal(author.ID())).
    All(ctx)
```

### Filter by Author's Property

```go
// Posts by active authors
activePosts, _ := client.PostRepo().Query().
    Where(filter.Post.Author.IsActive.True()).
    All(ctx)
```

### Posts Without Editor

```go
// Posts needing an editor
needsEditor, _ := client.PostRepo().Query().
    Where(filter.Post.Editor.IsNil()).
    All(ctx)
```

### Self-Reference (e.g., Parent)

```go
type Category struct {
    som.Node

    Name   string
    Parent *Category  // Self-reference
}

// Find root categories
roots, _ := client.CategoryRepo().Query().
    Where(filter.Category.Parent.IsNil()).
    All(ctx)

// Find children of specific category
children, _ := client.CategoryRepo().Query().
    Where(filter.Category.Parent.ID().Equal(parentID)).
    All(ctx)
```

## Complete Example

```go
package main

import (
    "context"
    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/filter"
    "yourproject/gen/som/with"
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create author
    author := &model.User{Name: "Alice", IsActive: true}
    client.UserRepo().Create(ctx, author)

    // Create post with reference
    post := &model.Post{
        Title:  "Hello World",
        Author: author,
    }
    client.PostRepo().Create(ctx, post)

    // Find posts by author ID
    authorPosts, _ := client.PostRepo().Query().
        Where(filter.Post.Author.ID().Equal(author.ID())).
        All(ctx)

    // Find posts by author name
    alicePosts, _ := client.PostRepo().Query().
        Where(filter.Post.Author.Name.Equal("Alice")).
        All(ctx)

    // Find posts by active authors only
    activePosts, _ := client.PostRepo().Query().
        Where(filter.Post.Author.IsActive.True()).
        All(ctx)

    // Posts without editor
    unedited, _ := client.PostRepo().Query().
        Where(filter.Post.Editor.IsNil()).
        All(ctx)

    // Sort by author name
    sorted, _ := client.PostRepo().Query().
        Order(by.Post.Author.Name.Asc()).
        All(ctx)

    // Fetch with eager loading
    postsWithAuthor, _ := client.PostRepo().Query().
        Fetch(with.Post.Author...).
        All(ctx)

    for _, p := range postsWithAuthor {
        fmt.Printf("%s by %s\n", p.Title, p.Author.Name)
    }

    // Complex nested filter
    comments, _ := client.CommentRepo().Query().
        Where(
            filter.Comment.Post.Author.IsActive.True(),
            filter.Comment.Author.Name.StartsWith("Admin"),
        ).
        All(ctx)
}
```

## Filter Reference Table

| Operation | Description | Returns |
|-----------|-------------|---------|
| `ID()` | Access record ID | ID filter |
| `<FieldName>` | Access nested field | Field's filter type |
| `IsNil()` | Reference is null | Bool filter |
| `IsNotNil()` | Reference is not null | Bool filter |

### ID Filter Operations

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Equal(id)` | ID equals | Bool filter |
| `NotEqual(id)` | ID not equals | Bool filter |
| `In(ids...)` | ID in set | Bool filter |
| `NotIn(ids...)` | ID not in set | Bool filter |

### Nested Field Access

All fields of the referenced node type are accessible:

```go
filter.Post.Author.Name        // String filter
filter.Post.Author.Email       // Email filter
filter.Post.Author.IsActive    // Bool filter
filter.Post.Author.CreatedAt   // Time filter
```
