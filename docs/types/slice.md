# Slice Type

The slice type handles arrays of any other type with extensive membership, set, and aggregation operations.

## Overview

| Property | Value |
|----------|-------|
| Go Type | `[]T` for any supported type T |
| Database Schema | `array<T>` / `option<array<T>>` |
| CBOR Encoding | Direct |
| Sortable | No (use element access for sorting) |

## Definition

```go
type Post struct {
    som.Node

    Title    string
    Tags     []string        // String slice
    Scores   []int           // Integer slice
    Ratings  []float64       // Float slice
    Authors  []*User         // Node slice
    Metadata []CustomStruct  // Struct slice
}
```

### Byte Slices

Byte slices have special handling:

```go
type Document struct {
    som.Node

    Data []byte  // Stored as bytes type
}
```

## Schema

Generated SurrealDB schema:

```surql
DEFINE FIELD tags ON post TYPE option<array<string>>;
DEFINE FIELD scores ON post TYPE option<array<int>>;
DEFINE FIELD ratings ON post TYPE option<array<float>>;
DEFINE FIELD authors ON post TYPE option<array<record<user>>>;
DEFINE FIELD data ON document TYPE option<bytes>;
```

## Filter Operations

### Empty/Nil Checks

```go
// Check if empty
filter.Post.Tags.IsEmpty()

// Check if not empty
filter.Post.Tags.NotEmpty()

// Explicit empty check
filter.Post.Tags.Empty(true)   // Is empty
filter.Post.Tags.Empty(false)  // Not empty
```

### Single Element Membership

```go
// Contains single element
filter.Post.Tags.Contains("golang")

// Does not contain element
filter.Post.Tags.ContainsNot("deprecated")
```

### Multiple Element Membership

```go
// Contains ALL elements
filter.Post.Tags.ContainsAll("golang", "database", "orm")

// Contains ANY element
filter.Post.Tags.ContainsAny("golang", "rust", "python")

// Contains NONE of elements
filter.Post.Tags.ContainsNone("deprecated", "obsolete")
```

### Equality Checks

```go
// Any element equals value
filter.Post.Scores.AnyEqual(100)

// All elements equal value
filter.Post.Scores.AllEqual(0)

// Any element fuzzy matches (strings)
filter.Post.Tags.AnyFuzzyMatch("go*")

// All elements fuzzy match
filter.Post.Tags.AllFuzzyMatch("*lang")
```

### Set Membership (Is-In)

```go
// Any element is in set
filter.Post.Tags.AnyIn("golang", "rust", "python")

// All elements are in set
filter.Post.Tags.AllIn("golang", "rust", "python")

// No elements are in set
filter.Post.Tags.NoneIn("deprecated", "obsolete")
```

### Element Access

```go
// First element
filter.Post.Tags.First().Equal("primary")

// Last element
filter.Post.Tags.Last().Equal("latest")

// Element at index
filter.Post.Tags.At(2).Equal("third")
```

### Length Operations

```go
// Get length
filter.Post.Tags.Len().GreaterThan(0)

// Minimum length
filter.Post.Tags.Len().GreaterThanEqual(3)

// Maximum length
filter.Post.Tags.Len().LessThanEqual(10)
```

### Array Transformations

```go
// Distinct elements
filter.Post.Tags.Distinct().Len().Equal(5)

// Reverse order
filter.Post.Tags.Reverse().First().Equal("last-tag")

// Sort ascending
filter.Post.Scores.SortAsc().First().Equal(minScore)

// Sort descending
filter.Post.Scores.SortDesc().First().Equal(maxScore)

// Slice (start, length)
filter.Post.Tags.Slice(0, 3).Contains("featured")
```

### Set Operations

```go
// Concatenate arrays
filter.Post.Tags.Concat(otherTags).Len().GreaterThan(10)

// Intersection
filter.Post.Tags.Intersect(requiredTags).NotEmpty()

// Union
filter.Post.Tags.Union(additionalTags).Contains("merged")

// Complement (elements in first but not second)
filter.Post.Tags.Complement(excludeTags).NotEmpty()

// Difference (symmetric difference)
filter.Post.Tags.Diff(otherTags).IsEmpty()
```

### Search Operations

```go
// Find index of element
filter.Post.Tags.FindIndex("featured").Equal(0)

// Filter to matching indices
filter.Post.Tags.FilterIndex("go*").NotEmpty()
```

### Aggregation Operations

```go
// Check if all elements are truthy
filter.Post.Tags.All().True()

// Check if any element is truthy
filter.Post.Tags.Any().True()

// Get minimum value (numeric)
filter.Post.Scores.Min().GreaterThan(0)

// Get maximum value (numeric)
filter.Post.Scores.Max().LessThan(100)

// Join elements (strings)
filter.Post.Tags.Join(",").Contains("golang")
```

### Boolean Array Operations

For boolean arrays:

```go
// Boolean AND of all elements
filter.Post.Flags.BooleanAnd().True()

// Boolean OR of all elements
filter.Post.Flags.BooleanOr().True()

// Boolean XOR of all elements
filter.Post.Flags.BooleanXor().True()

// Boolean NOT of all elements
filter.Post.Flags.BooleanNot()
```

### Logical Array Operations

```go
// Logical AND
filter.Post.Values.LogicalAnd()

// Logical OR
filter.Post.Values.LogicalOr()

// Logical XOR
filter.Post.Values.LogicalXor()
```

### Pattern Matching

```go
// Elements matching pattern
filter.Post.Tags.Matches("go*").NotEmpty()
```

## Sorting

Slices themselves are not directly sortable. Use element access:

```go
// Sort by first tag
query.Order(by.Post.Tags.First().Asc())

// Sort by array length
query.Order(by.Post.Tags.Len().Desc())

// Sort by max score
query.Order(by.Post.Scores.Max().Desc())
```

## Method Chaining

Slice filters support powerful chaining:

```go
// Distinct tags, count > 3
filter.Post.Tags.Distinct().Len().GreaterThan(3)

// First tag starts with prefix
filter.Post.Tags.First().StartsWith("featured")

// Sorted scores, max value
filter.Post.Scores.SortDesc().First().GreaterThan(90)

// Top 3 tags contain value
filter.Post.Tags.Slice(0, 3).Contains("important")
```

## Common Patterns

### Has Any Tags

```go
// Posts with at least one tag
tagged, _ := client.PostRepo().Query().
    Where(filter.Post.Tags.NotEmpty()).
    All(ctx)
```

### Filter by Tag

```go
// Posts with specific tag
golangPosts, _ := client.PostRepo().Query().
    Where(filter.Post.Tags.Contains("golang")).
    All(ctx)
```

### Multiple Required Tags

```go
// Posts with all required tags
tutorials, _ := client.PostRepo().Query().
    Where(
        filter.Post.Tags.ContainsAll("golang", "tutorial", "beginner"),
    ).
    All(ctx)
```

### Score Range

```go
// Posts with high scores
highScoring, _ := client.PostRepo().Query().
    Where(filter.Post.Scores.Max().GreaterThan(90)).
    All(ctx)
```

### Unique Tags Only

```go
// Posts where all tags are unique
uniqueTags, _ := client.PostRepo().Query().
    Where(
        filter.Post.Tags.Len().Equal(filter.Post.Tags.Distinct().Len()),
    ).
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
)

func main() {
    ctx := context.Background()
    client, _ := som.NewClient(ctx, som.Config{...})

    // Create post with tags
    post := &model.Post{
        Title:  "Getting Started with Go",
        Tags:   []string{"golang", "tutorial", "beginner"},
        Scores: []int{85, 90, 88},
    }
    client.PostRepo().Create(ctx, post)

    // Find posts with any tags
    tagged, _ := client.PostRepo().Query().
        Where(filter.Post.Tags.NotEmpty()).
        All(ctx)

    // Find golang posts
    goPosts, _ := client.PostRepo().Query().
        Where(filter.Post.Tags.Contains("golang")).
        All(ctx)

    // Find tutorials for beginners
    beginnerTutorials, _ := client.PostRepo().Query().
        Where(
            filter.Post.Tags.ContainsAll("tutorial", "beginner"),
        ).
        All(ctx)

    // Find posts with good ratings
    goodPosts, _ := client.PostRepo().Query().
        Where(filter.Post.Scores.Min().GreaterThanEqual(80)).
        All(ctx)

    // Find posts with many tags
    manyTags, _ := client.PostRepo().Query().
        Where(filter.Post.Tags.Len().GreaterThan(5)).
        All(ctx)

    // Sort by number of tags
    byTagCount, _ := client.PostRepo().Query().
        Order(by.Post.Tags.Len().Desc()).
        All(ctx)

    // Posts without deprecated tags
    notDeprecated, _ := client.PostRepo().Query().
        Where(
            filter.Post.Tags.ContainsNone("deprecated", "obsolete"),
        ).
        All(ctx)
}
```

## Filter Reference Table

| Operation | Description | Returns |
|-----------|-------------|---------|
| `IsEmpty()` | Check if empty | Bool filter |
| `NotEmpty()` | Check if not empty | Bool filter |
| `Empty(bool)` | Explicit empty check | Bool filter |
| `Contains(val)` | Contains element | Bool filter |
| `ContainsNot(val)` | Doesn't contain | Bool filter |
| `ContainsAll(vals...)` | Contains all | Bool filter |
| `ContainsAny(vals...)` | Contains any | Bool filter |
| `ContainsNone(vals...)` | Contains none | Bool filter |
| `AnyEqual(val)` | Any equals | Bool filter |
| `AllEqual(val)` | All equal | Bool filter |
| `AnyFuzzyMatch(pattern)` | Any matches pattern | Bool filter |
| `AllFuzzyMatch(pattern)` | All match pattern | Bool filter |
| `AnyIn(vals...)` | Any in set | Bool filter |
| `AllIn(vals...)` | All in set | Bool filter |
| `NoneIn(vals...)` | None in set | Bool filter |
| `First()` | First element | Element filter |
| `Last()` | Last element | Element filter |
| `At(index)` | Element at index | Element filter |
| `Len()` | Array length | Numeric filter |
| `Distinct()` | Unique elements | Slice filter |
| `Reverse()` | Reversed array | Slice filter |
| `SortAsc()` | Sorted ascending | Slice filter |
| `SortDesc()` | Sorted descending | Slice filter |
| `Slice(start, len)` | Subsequence | Slice filter |
| `Concat(arr)` | Concatenate | Slice filter |
| `Intersect(arr)` | Intersection | Slice filter |
| `Union(arr)` | Union | Slice filter |
| `Complement(arr)` | Complement | Slice filter |
| `Diff(arr)` | Difference | Slice filter |
| `FindIndex(val)` | Find index | Numeric filter |
| `FilterIndex(val)` | Filter indices | Slice filter |
| `All()` | All truthy | Bool filter |
| `Any()` | Any truthy | Bool filter |
| `Min()` | Minimum value | Numeric filter |
| `Max()` | Maximum value | Numeric filter |
| `Join(sep)` | Join to string | String filter |
| `Matches(pattern)` | Pattern match | Slice filter |
