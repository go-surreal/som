# Field Constraints

Constraints declared on a field become an `ASSERT` clause on its `DEFINE FIELD` statement, so they
are enforced by the database on every write. They are never checked on read: a row written before a
rule existed stays loadable.

## Bounds

Simple bounds are declared on the field itself:

```go
type User struct {
    som.Node[som.ULID]

    Name     string   `som:"len=3..64"` // between 3 and 64 characters
    Code     string   `som:"len=4"`     // exactly 4 characters
    Nickname *string  `som:"len=2.."`   // at least 2, no upper bound
    Initials string   `som:"len=..3"`   // at most 3
    Age      int      `som:"min=0,max=130"`
    Tags     []string `som:"len=..5"`   // at most 5 items
}
```

`len` applies to strings and slices. On a slice it counts the items, not the size of any individual
element. `min` and `max` apply to numeric fields and may be combined or used on their own.

An optional field — a pointer, or a slice, which is stored as an option so that a nil slice stays
distinguishable from an empty one — is only constrained once it holds a value. Leaving it absent is
always allowed.

Tag constraints compose with the ones som derives from the Go type. An `int8 `som:"min=0"`` still
rejects anything above 127.

## Patterns and Expressions

Anything that does not fit into a tag without quoting or escaping is declared in `som.define.go` and
referenced by name, the same way a full-text search configuration is:

```go
// som.define.go
var (
    phoneFormat = define.Assert("phone_format").
        Regex(`^[0-9+\-]+$`).
        Message("must only contain digits, plus and minus")

    noPlaceholder = define.Assert("no_placeholder").
        Raw(`string::lowercase($value) != "tbd"`)
)

func Definitions() define.Definitions {
    return define.Definitions{
        Asserts: []*define.AssertBuilder{phoneFormat, noPlaceholder},
    }
}
```

```go
type Contact struct {
    som.Node[som.ULID]

    Phone string `som:"assert=phone_format"`
    Title string `som:"len=1..80,assert=no_placeholder"`
}
```

Keeping the pattern in a Go string literal means no tag escaping, real syntax highlighting, and one
place to change it when it is used on several fields.

`Regex` is evaluated by SurrealDB, not by Go's `regexp` package. Both are RE2-based, but only the
subset the two engines agree on is portable.

`Raw` is the escape hatch for constraints the builder does not model. The expression is embedded
into the schema verbatim and is not validated by som; the value under validation is `$value`.

`Message` sets what a violation reports. Without it, som derives a message from the constraint.

## Handling Violations

A rejected write returns an error matching `som.ErrAssert`, carrying the field and the rule it
broke:

```go
err := client.ContactRepo().Create(ctx, &contact)

var assertErr *som.AssertError
if errors.As(err, &assertErr) {
    fmt.Println(assertErr.Field, assertErr.Message)
    // phone   must only contain digits, plus and minus
}
```

`AssertError.Field` is the **database** path of the field, including the prefix of any struct it is
nested in — the constraint is enforced by the database, which only knows the field under that name.
The underlying database error is available via `errors.Unwrap`.

A caller that only needs to turn the rejection into a response can match the broader `som.ErrInvalid`
instead, which covers every write rejected because of a value it carried:

```go
if errors.Is(err, som.ErrInvalid) {
    http.Error(w, err.Error(), http.StatusUnprocessableEntity)
}
```

## Unsupported Fields

Constraints are rejected at generation time on fields where they could not work:

- `som.Password` — SurrealDB evaluates `ASSERT` after `VALUE`, so the constraint would be checked
  against the stored hash rather than the password it was written for.
- Record links and edges — an edge is not a real field in the database schema, so there is no
  `DEFINE FIELD` statement for the constraint to attach to.

Declaring `len` on a numeric field, or `min`/`max` on a string, is rejected for the same reason: the
constraint could never fire.
