# Time Travel

SurrealKV can keep every past version of a record. With versioning enabled on the server, a query
can read the database as it was at any earlier point in time — `SELECT ... VERSION <datetime>` in
SurrealQL, `Version(t)` in SOM.

## Opting In

Time travel is opt-in on both sides, because a server without versioning rejects these queries.

**1. Run SurrealKV with versioning enabled.** The in-memory store (`memory`, the default of
`surreal start`) does not support it:

```bash
surreal start "surrealkv:///data/db?versioned=true"
```

Add `&retention=30d` to limit how long old versions are kept; without it they are kept forever.

**2. Generate with `--versioned`:**

```bash
go run github.com/go-surreal/som@latest -i ./model --versioned
```

Only then does the query builder have a `Version` method. Without the flag the method does not
exist, so the code cannot accidentally depend on a server feature that is not there.

## Querying the Past

```go
snapshot := time.Now().Add(-1 * time.Hour)

orders, err := client.OrderRepo().Query().
    Where(filter.Order.Total.GreaterThan(100)).
    Fetch(with.Order.Customer()).
    Version(snapshot).
    All(ctx)
```

`Version` is the last call before the terminal. Everything chained before it — filters, ordering,
fetches — applies to the past state, as do the terminals `All`, `First`, `Count`, `Exists`,
`Iterate`, `Paginate` and the fragment variants (`AllAs`, ...):

- **Fetched relations** are resolved at the same point in time.
- **Soft-deleted records** are excluded based on their state at that time.
- **Expiry** is evaluated against the given time, so records that have expired since are still
  part of the result.

## Limitations

- Live queries always follow the current state, and the full-text index only holds the current
  state as well. `Version` returns a `query.Versioned[M]`, which has neither `Live` nor `Search`,
  so combining them is a compile error.
- `query.Distinct` only accepts a plain `query.Builder[M]` and thus does not work on a versioned
  query.
- A time before a table was defined fails with a "table does not exist" error.
- Versions are stamped by the server clock. Times taken on a client with a drifting clock may
  resolve to a slightly different state.
- History is kept for the whole datastore, not per table, which grows storage with every write.
  Use `retention` to bound it.

## Time Travel vs Changefeed

| Need | Use |
|------|-----|
| Read records (with filters and relations) as they were at time T | Time travel |
| Replay the individual changes to a table since T | [Changefeed](08_changefeed.md) |
