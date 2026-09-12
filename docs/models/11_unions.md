# Unions

A union is a set of node models reachable through a common interface. It has no
table of its own. It exists so that a single link can point to records of more
than one table, which SurrealDB supports natively through the multi-table record
type `record<a|b>`.

## Declaration

A union is an interface in the model package that embeds `som.Union` with
itself as type parameter:

```go
type Shape interface {
    som.Union[Shape]
}
```

A node becomes a member of that union by embedding `som.Part` of it:

```go
type Square struct {
    som.Node[som.ULID]
    som.Part[Shape]

    Label string
    Size  float64
}

type Circle struct {
    som.Node[som.UUID]
    som.Part[Shape]

    Label  string
    Radius float64
}
```

Membership is checked by the compiler: the method `som.Union[Shape]` requires is
unexported and carries `Shape` as its parameter type, so only `som.Part[Shape]`
can satisfy it. Embedding `som.Part[SomethingElse]` does not make a model a
`Shape`.

Members do not have to agree on anything else. They may have different ID types
and different features — above, `Square` uses a ULID and `Circle` a UUID.

A model can be part of exactly one *sealed* union. Two `som.Part` embeds would
both be named `Part`, which Go rejects, and even under different names only one
of the two `unionOf` methods could be promoted. Further memberships are declared
through an open union, see below.

## Open unions

A union declared with `som.OpenUnion` requires no marker method. What a member
is checked against is what the interface declares itself, so an open union must
declare at least one method:

```go
type Colored interface {
    som.OpenUnion

    Color() string
}
```

A node joins it with a **blank** `som.Part` field, which declares the membership
without adding a method. Blank fields may repeat, so a node can be a member of
one sealed union and any number of open ones:

```go
type Square struct {
    som.Node[som.ULID]
    som.Part[Shape]

    _ som.Part[Colored]

    Label string
    Size  float64
}

func (s Square) Color() string {
    return "red"
}
```

Everything downstream — links, relation endpoints, filters, the union
repository — works the same for both kinds. The only difference is the check at
the assignment site: a sealed union accepts exactly its declared members, an
open one accepts anything implementing its methods. Assigning a non-member to an
open union field compiles, and the write then silently stores no link for that
field, so keep the interface's methods specific enough to make that unlikely.

The two forms must not be mixed: a blank field on a sealed union, or a
`som.Part` embed on an open one, is rejected by the generator.

## Methods

Next to the `som.Union` embed a union may declare methods of its own. They are
ordinary Go methods, invisible to the database, and every member has to
implement them:

```go
type Shape interface {
    som.Union[Shape]

    Area() float64
}

func (s Square) Area() float64 {
    return s.Size * s.Size
}

func (c *Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}
```

A pointer receiver is fine — a decoded member is always handed back as a
pointer. Note that a method computes on whatever the model currently holds, so
calling it on an unfetched link, which carries only the record id, operates on
zero-valued fields. Fetch the link first if the result has to be meaningful.

## Union links

A field typed as the union interface is a record link that may point to any
member:

```go
type Drawing struct {
    som.Node[som.ULID]

    Name    string
    Subject Shape
}
```

```surql
DEFINE FIELD subject ON TABLE drawing TYPE option<record<square|circle>>;
```

A union field is never a pointer — an unset link is the nil interface.

Writing a link stores the record id of whichever member the field holds:

```go
square := &model.Square{Label: "big", Size: 4}
client.SquareRepo().Create(ctx, square)

client.DrawingRepo().Create(ctx, &model.Drawing{
    Name:    "sketch",
    Subject: square,
})
```

Reading one decodes it back into that member, so a type switch recovers the
concrete model:

```go
drawing, _, _ := client.DrawingRepo().Read(ctx, id)

switch subject := drawing.Subject.(type) {
case *model.Square:
    fmt.Println(subject.Size)
case *model.Circle:
    fmt.Println(subject.Radius)
}
```

Like every link, an unfetched one carries only its record id and is marked as
partial. Use `Fetch` to load the whole record:

```go
drawings, _ := client.DrawingRepo().Query().
    Fetch(with.Drawing.Subject()).
    All(ctx)
```

## Filtering

The link itself can be matched by the member it points to, and it can be
narrowed to a single member to filter on that member's fields:

```go
// Drawings whose subject is a square
client.DrawingRepo().Query().
    Where(filter.Drawing.Subject().IsSquare()).
    All(ctx)

// Drawings of a square larger than 3
client.DrawingRepo().Query().
    Where(filter.Drawing.Subject().Square().Size.GreaterThan(3)).
    All(ctx)

// Drawings without a subject
client.DrawingRepo().Query().
    Where(filter.Drawing.Subject().Nil(true)).
    All(ctx)
```

Narrowing does not restrict the rows by itself: a link to a different member
simply holds no value for the field, so it never matches the comparison. Combine
it with `Is<Member>` when you want that stated explicitly.

## Relation endpoints

An edge may declare a union as its `in` or `out` end, which maps to a
multi-table relation:

```go
type Annotates struct {
    som.Edge

    Author  User  `som:"in"`
    Subject Shape `som:"out"`

    Note string
}
```

```surql
DEFINE TABLE annotates TYPE RELATION IN user OUT square|circle ENFORCED;
```

## Querying all members

A union has its own read-only repository. Its query selects from every member
table at once and decodes each row into the member it belongs to:

```go
shapes, err := client.ShapeRepo().Query().
    OrderBy(...).
    Limit(20).
    All(ctx)

for _, shape := range shapes {
    switch s := (*shape).(type) {
    case *model.Square:
        // ...
    }
}
```

The union repository is read-only. Creating, updating and deleting depends on
the id type, hooks and features of a concrete member, so those stay with that
member's repository.

If any member uses soft delete or expiry, the union query filters on
`deleted_at` and `expires_at` respectively. A member without the field yields no
value for it, so the filter never excludes its records.

## Limitations

- A slice of a union (`[]Shape`) is not supported yet.
- A union field is not tracked in the load state, so `Resolve` re-reads it on
  every call instead of skipping it when it was already fetched.
- A node can be a member of one *sealed* union only. Sealing rests on the single
  method `unionOf`, and a Go type cannot have that method twice. Use an open
  union for every further membership.
