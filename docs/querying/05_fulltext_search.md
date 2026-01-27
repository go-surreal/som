# Full-Text Search

Full-text search enables BM25-based relevance searching on string fields with fulltext indexes. SOM generates type-safe search operations that integrate seamlessly with the query builder.

## Basic Search

Use `.Matches()` on string fields to perform full-text search:

```go
import "yourproject/gen/som/filter"

results, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang tutorial")).
    AllMatches(ctx)
```

The search uses SurrealDB's full-text search capabilities with BM25 scoring for relevance ranking.

## Search vs SearchAll

SOM provides two methods for combining multiple search conditions:

### Search (OR Semantics)

`Search()` combines conditions with OR - documents matching **any** term are returned:

```go
// Returns articles containing "golang" OR "rust"
results, err := client.ArticleRepo().Query().
    Search(
        filter.Article.Content.Matches("golang"),
        filter.Article.Content.Matches("rust"),
    ).
    AllMatches(ctx)
```

This is the default search engine behavior where broader results are preferred.

### SearchAll (AND Semantics)

`SearchAll()` combines conditions with AND - documents must match **all** terms:

```go
// Returns only articles containing BOTH "golang" AND "tutorial"
results, err := client.ArticleRepo().Query().
    SearchAll(
        filter.Article.Content.Matches("golang"),
        filter.Article.Content.Matches("tutorial"),
    ).
    AllMatches(ctx)
```

Use this when documents must contain all specified terms.

## Execution Methods

Search queries have specialized execution methods:

| Method | Returns | Description |
|--------|---------|-------------|
| `AllMatches(ctx)` | `[]SearchResult[M], error` | All results with search metadata |
| `FirstMatch(ctx)` | `*SearchResult[M], bool, error` | First result with metadata |
| `All(ctx)` | `[]*M, error` | Plain models without metadata |

### AllMatches

Returns all matching documents with search metadata (scores, highlights, offsets):

```go
results, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang")).
    AllMatches(ctx)

for _, result := range results {
    fmt.Printf("Score: %f, Title: %s\n", result.Score(), result.Model.Title)
}
```

### FirstMatch

Returns the first matching document:

```go
result, found, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang")).
    FirstMatch(ctx)

if found {
    fmt.Printf("Best match: %s (score: %f)\n", result.Model.Title, result.Score())
}
```

### All

Returns plain models when you don't need search metadata:

```go
articles, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang")).
    All(ctx)
```

## SearchResult Type

`AllMatches()` and `FirstMatch()` return `SearchResult` containing the model and search metadata:

```go
type SearchResult[M any] struct {
    Model      M                    // The matched model
    Scores     []float64            // BM25 relevance scores
    Highlights map[int]string       // Highlighted text by ref
    Offsets    map[int][]Offset     // Match positions by ref
}
```

### Helper Methods

```go
// Get the primary score (first in slice)
score := result.Score()

// Get highlighted text for a specific ref (defaults to 0)
highlighted := result.Highlighted()      // ref 0
highlighted := result.Highlighted(1)     // ref 1

// Get match offsets for a specific ref
offsets := result.Offset()               // ref 0
offsets := result.Offset(1)              // ref 1
```

## Search Options

### Predicate References

When searching multiple fields, use `.Ref()` to set explicit predicate references:

```go
results, err := client.ArticleRepo().Query().
    Search(
        filter.Article.Title.Matches("golang").Ref(0),
        filter.Article.Content.Matches("golang").Ref(1),
    ).
    AllMatches(ctx)

// Access scores by ref
titleScore := results[0].Scores[0]    // Score for title match
contentScore := results[0].Scores[1]  // Score for content match
```

References are auto-assigned starting from 0 if not specified.

### Highlighting

Enable highlighted snippets with custom prefix/suffix tags:

```go
results, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang").WithHighlights("<mark>", "</mark>")).
    AllMatches(ctx)

// Get highlighted text
snippet := results[0].Highlighted()
// Returns: "Learn <mark>golang</mark> with this tutorial"
```

### Offsets

Enable offset tracking to get match positions:

```go
results, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang").WithOffsets()).
    AllMatches(ctx)

// Get match positions
for _, offset := range results[0].Offset() {
    fmt.Printf("Match at position %d-%d\n", offset.Start, offset.End)
}
```

### Combining Options

Options can be chained:

```go
filter.Article.Content.Matches("golang").
    Ref(0).
    WithHighlights("<b>", "</b>").
    WithOffsets()
```

## Score Sorting

Sort results by relevance score using the `query` package:

```go
import "yourproject/gen/som/query"

results, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang")).
    Order(query.Score(0).Desc()).
    AllMatches(ctx)
```

### Basic Score Sorting

```go
// Sort by score descending (most relevant first)
query.Score(0).Desc()

// Sort by score ascending (least relevant first)
query.Score(0).Asc()
```

### Multiple Score References

When searching multiple fields, reference specific scores:

```go
results, err := client.ArticleRepo().Query().
    Search(
        filter.Article.Title.Matches("golang").Ref(0),
        filter.Article.Content.Matches("golang").Ref(1),
    ).
    Order(query.Score(0).Desc()).  // Sort by title score
    AllMatches(ctx)
```

### Score Combination Modes

Combine multiple scores with different strategies:

```go
// Sum scores (default)
query.Score(0, 1).Sum().Desc()

// Maximum score
query.Score(0, 1).Max().Desc()

// Average score
query.Score(0, 1).Average().Desc()

// Weighted combination (title 2x, content 0.5x)
query.Score(0, 1).Weighted(2.0, 0.5).Desc()
```

### Mixed Sorting

Combine score sorting with field sorting:

```go
results, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang")).
    Order(
        query.Score(0).Desc(),           // Primary: relevance
        by.Article.CreatedAt.Desc(),     // Secondary: newest first
    ).
    AllMatches(ctx)
```

## Combining Search and Filters

Search conditions can be combined with regular filters using `Where()`:

```go
results, err := client.ArticleRepo().Query().
    Search(filter.Article.Content.Matches("golang")).
    Where(
        filter.Article.Published.IsTrue(),
        filter.Article.Category.Equal("tutorials"),
    ).
    AllMatches(ctx)
```

The search conditions and filters are combined with AND in the WHERE clause.

## Complete Example

```go
package main

import (
    "context"
    "fmt"

    "yourproject/gen/som"
    "yourproject/gen/som/query"
    "yourproject/gen/som/filter"
)

func SearchArticles(ctx context.Context, client *som.Client, terms string) {
    // Multi-field search with relevance ranking
    results, err := client.ArticleRepo().Query().
        Search(
            filter.Article.Title.Matches(terms).Ref(0).WithHighlights("<b>", "</b>"),
            filter.Article.Content.Matches(terms).Ref(1).WithHighlights("<mark>", "</mark>"),
        ).
        Where(
            filter.Article.Published.IsTrue(),
        ).
        Order(
            query.Score(0, 1).Weighted(2.0, 1.0).Desc(),  // Title matches weighted 2x
        ).
        Limit(20).
        AllMatches(ctx)

    if err != nil {
        panic(err)
    }

    for _, result := range results {
        fmt.Printf("Score: %.2f\n", result.Score())
        fmt.Printf("Title: %s\n", result.Model.Title)
        fmt.Printf("Title highlight: %s\n", result.Highlighted(0))
        fmt.Printf("Content snippet: %s\n", result.Highlighted(1))
        fmt.Println("---")
    }
}
```

## Best Practices

1. **Use appropriate semantics**: Use `Search()` for broad searches, `SearchAll()` for precise matches
2. **Set explicit refs**: When searching multiple fields, use `.Ref()` for predictable score access
3. **Weight important fields**: Use `.Weighted()` to boost matches in important fields (e.g., title)
4. **Combine with filters**: Use `Where()` for non-search criteria to narrow results efficiently
5. **Consider highlights**: Enable highlighting for user-facing search results
6. **Use `All()` when metadata isn't needed**: It's more efficient than `AllMatches()` if you don't need scores
