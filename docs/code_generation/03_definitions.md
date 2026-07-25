# Schema Definitions

Some parts of the schema cannot be expressed with struct tags: full-text analyzers, search index
configurations and [view](../models/07_views.md) projections. Those live in a definition file
that the generator compiles and executes during generation.

## The Definition File

Create a file in the **module root** (next to `go.mod`), guarded by the `som` build tag, that
exposes a `Definitions()` function:

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
        Searches: []*define.SearchBuilder{searchEnglish, searchAutocomplete},
        Views:    []define.ViewDefinition{eventSummary},
    }
}
```

Two details matter:

- **`//go:build som`** — the file is excluded from your normal build and only compiled by the
  generator. Definitions reference generated packages, so without the tag your project would not
  build before the first generation.
- **Module root, not the model package** — a view definition references generated filter refs,
  which import the model package. Putting the definition into the model package would create an
  import cycle.

If no `//go:build som` file exists in the module root, the generator simply skips this step.

## Analyzers

An analyzer describes how text is tokenized and normalized for full-text search:

```go
var (
    english = define.FulltextAnalyzer("english").
        Tokenizers(define.Blank, define.Punct).
        Filters(define.Lowercase, define.Snowball(define.English))

    autocomplete = define.FulltextAnalyzer("autocomplete").
        Tokenizers(define.Class).
        Filters(define.Lowercase, define.Edgengram(1, 10))
)
```

### Tokenizers

| Tokenizer | Splits on |
|-----------|-----------|
| `define.Blank` | Whitespace (space, tab, newline) |
| `define.Camel` | Uppercase transitions (camelCase / PascalCase) |
| `define.Class` | Unicode class changes (digit, letter, punctuation, blank) |
| `define.Punct` | Punctuation characters |

### Filters

| Filter | Effect |
|--------|--------|
| `define.Ascii` | Strip diacritics to ASCII equivalents |
| `define.Lowercase` / `define.Uppercase` | Normalize case |
| `define.Snowball(lang)` | Snowball stemming + lowercase |
| `define.Edgengram(min, max)` | Prefix tokens (autocomplete) |
| `define.Ngram(min, max)` | N-gram tokens |
| `define.Mapper(path)` | Lemmatization via a dictionary file |

`define.Snowball` takes a language constant such as `define.English`, `define.German`,
`define.French`, and so on.

Additional builder options: `Function(fn)` for a custom analyzer function and `Comment(text)`.

## Search Configurations

A search configuration binds an analyzer to a full-text index and sets its scoring options:

```go
var (
    searchEnglish = define.Search("english_search").
        FulltextAnalyzer(english).
        BM25(1.2, 0.75).
        Highlights()

    searchAutocomplete = define.Search("autocomplete_search").
        FulltextAnalyzer(autocomplete)
)
```

| Method | Description |
|--------|-------------|
| `FulltextAnalyzer(a)` | Analyzer used by the index |
| `BM25(k1, b)` | BM25 ranking parameters |
| `Highlights()` | Enable highlight support |
| `Concurrently()` | Build the index concurrently |

Reference the configuration by name from your model:

```go
type Article struct {
    som.Node[som.ULID]

    Title   string `som:"fulltext=english_search"`
    Content string `som:"fulltext=english_search"`
}
```

Highlighting is only available when the configuration enables `Highlights()`. See
[Full-Text Search](../querying/05_fulltext_search.md) for querying.

## Views

View projections are declared with `define.View` and listed under `Views`:

```go
var eventSummary = define.View[model.EventSummary, model.EventLog]().
    Project(
        define.As(filter.EventSummary.Category, filter.EventLog.Category),
        define.As(filter.EventSummary.Total, aggregate.Count(filter.EventLog.Category)),
        define.As(filter.EventSummary.AvgValue, aggregate.Mean(filter.EventLog.Value)),
    ).
    GroupBy(filter.EventLog.Category)
```

See [Views](../models/07_views.md) for the full description.

## Applying the Schema

Analyzers, search indexes and views end up in the generated schema and are created by:

```go
err := client.ApplySchema(ctx)
```
